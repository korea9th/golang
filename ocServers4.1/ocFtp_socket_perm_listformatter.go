// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"bytes"
	"fmt"
)

// DataSocket describes a data socket is used to send non-control data between the client and
// server.
type DataSocket interface {
	Host() string

	Port() int

	// the standard io.Reader interface
	Read(p []byte) (n int, err error)

	// the standard io.ReaderFrom interface
	ReadFrom(r io.Reader) (int64, error)

	// the standard io.Writer interface
	Write(p []byte) (n int, err error)

	// the standard io.Closer interface
	Close() error
}

type ftpActiveSocket struct {
	conn   *net.TCPConn
	host   string
	port   int
	logger Logger
}

func newActiveSocket(remote string, port int, logger Logger, sessionID string) (DataSocket, error) {
	connectTo := net.JoinHostPort(remote, strconv.Itoa(port))

	logger.Print(sessionID, "Opening active data connection to "+connectTo)

	raddr, err := net.ResolveTCPAddr("tcp", connectTo)

	if err != nil {
		logger.Print(sessionID, err)
		return nil, err
	}

	tcpConn, err := net.DialTCP("tcp", nil, raddr)

	if err != nil {
		logger.Print(sessionID, err)
		return nil, err
	}

	socket := new(ftpActiveSocket)
	socket.conn = tcpConn
	socket.host = remote
	socket.port = port
	socket.logger = logger

	return socket, nil
}

func (socket *ftpActiveSocket) Host() string {
	return socket.host
}

func (socket *ftpActiveSocket) Port() int {
	return socket.port
}

func (socket *ftpActiveSocket) Read(p []byte) (n int, err error) {
	return socket.conn.Read(p)
}

func (socket *ftpActiveSocket) ReadFrom(r io.Reader) (int64, error) {
	return socket.conn.ReadFrom(r)
}

func (socket *ftpActiveSocket) Write(p []byte) (n int, err error) {
	return socket.conn.Write(p)
}

func (socket *ftpActiveSocket) Close() error {
	return socket.conn.Close()
}

type ftpPassiveSocket struct {
	conn      net.Conn
	port      int
	host      string
	ingress   chan []byte
	egress    chan []byte
	logger    Logger
	lock      sync.Mutex // protects conn and err
	err       error
	tlsConfig *tls.Config
}

// Detect if an error is "bind: address already in use"
//
// Originally from https://stackoverflow.com/a/52152912/164234
func isErrorAddressAlreadyInUse(err error) bool {
	errOpError, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	errSyscallError, ok := errOpError.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	errErrno, ok := errSyscallError.Err.(syscall.Errno)
	if !ok {
		return false
	}
	if errErrno == syscall.EADDRINUSE {
		return true
	}
	const WSAEADDRINUSE = 10048
	if runtime.GOOS == "windows" && errErrno == WSAEADDRINUSE {
		return true
	}
	return false
}

func newPassiveSocket(host string, port func() int, logger Logger, sessionID string, tlsConfig *tls.Config) (DataSocket, error) {
	socket := new(ftpPassiveSocket)
	socket.ingress = make(chan []byte)
	socket.egress = make(chan []byte)
	socket.logger = logger
	socket.host = host
	socket.tlsConfig = tlsConfig
	const retries = 10
	var err error
	for i := 1; i <= retries; i++ {
		socket.port = port()
		err = socket.GoListenAndServe(sessionID)
		if err != nil && socket.port != 0 && isErrorAddressAlreadyInUse(err) {
			// choose a different port on error already in use
			continue
		}
		break
	}
	return socket, err
}

func (socket *ftpPassiveSocket) Host() string {
	return socket.host
}

func (socket *ftpPassiveSocket) Port() int {
	return socket.port
}

func (socket *ftpPassiveSocket) Read(p []byte) (n int, err error) {
	socket.lock.Lock()
	defer socket.lock.Unlock()
	if socket.err != nil {
		return 0, socket.err
	}
	return socket.conn.Read(p)
}

func (socket *ftpPassiveSocket) ReadFrom(r io.Reader) (int64, error) {
	socket.lock.Lock()
	defer socket.lock.Unlock()
	if socket.err != nil {
		return 0, socket.err
	}

	// For normal TCPConn, this will use sendfile syscall; if not,
	// it will just downgrade to normal read/write procedure
	return io.Copy(socket.conn, r)
}

func (socket *ftpPassiveSocket) Write(p []byte) (n int, err error) {
	socket.lock.Lock()
	defer socket.lock.Unlock()
	if socket.err != nil {
		return 0, socket.err
	}
	return socket.conn.Write(p)
}

func (socket *ftpPassiveSocket) Close() error {
	socket.lock.Lock()
	defer socket.lock.Unlock()
	if socket.conn != nil {
		return socket.conn.Close()
	}
	return nil
}

func (socket *ftpPassiveSocket) GoListenAndServe(sessionID string) (err error) {
	laddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort("", strconv.Itoa(socket.port)))
	if err != nil {
		socket.logger.Print(sessionID, err)
		return
	}

	var tcplistener *net.TCPListener
	tcplistener, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		socket.logger.Print(sessionID, err)
		return
	}

	// The timeout, for a remote client to establish connection
	// with a PASV style data connection.
	const acceptTimeout = 60 * time.Second
	err = tcplistener.SetDeadline(time.Now().Add(acceptTimeout))
	if err != nil {
		socket.logger.Print(sessionID, err)
		return
	}

	var listener net.Listener = tcplistener
	add := listener.Addr()
	parts := strings.Split(add.String(), ":")
	port, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		socket.logger.Print(sessionID, err)
		return
	}

	socket.port = port
	if socket.tlsConfig != nil {
		listener = tls.NewListener(listener, socket.tlsConfig)
	}

	socket.lock.Lock()
	go func() {
		defer socket.lock.Unlock()

		conn, err := listener.Accept()
		if err != nil {
			socket.err = err
			return
		}
		socket.err = nil
		socket.conn = conn
		_ = listener.Close()
	}()
	return nil
}

type Perm interface {
	GetOwner(string) (string, error)
	GetGroup(string) (string, error)
	GetMode(string) (os.FileMode, error)

	ChOwner(string, string) error
	ChGroup(string, string) error
	ChMode(string, os.FileMode) error
}

type SimplePerm struct {
	owner, group string
}

func NewSimplePerm(owner, group string) *SimplePerm {
	return &SimplePerm{
		owner: owner,
		group: group,
	}
}

func (s *SimplePerm) GetOwner(string) (string, error) {
	return s.owner, nil
}

func (s *SimplePerm) GetGroup(string) (string, error) {
	return s.group, nil
}

func (s *SimplePerm) GetMode(string) (os.FileMode, error) {
	return os.ModePerm, nil
}

func (s *SimplePerm) ChOwner(string, string) error {
	return nil
}

func (s *SimplePerm) ChGroup(string, string) error {
	return nil
}

func (s *SimplePerm) ChMode(string, os.FileMode) error {
	return nil
}

type listFormatter []FileInfo

// Short returns a string that lists the collection of files by name only,
// one per line
func (formatter listFormatter) Short() []byte {
	var buf bytes.Buffer
	for _, file := range formatter {
		fmt.Fprintf(&buf, "%s\r\n", file.Name())
	}
	return buf.Bytes()
}

// Detailed returns a string that lists the collection of files with extra
// detail, one per line
func (formatter listFormatter) Detailed() []byte {
	var buf bytes.Buffer
	for _, file := range formatter {
		fmt.Fprintf(&buf, file.Mode().String())
		fmt.Fprintf(&buf, " 1 %s %s ", file.Owner(), file.Group())
		fmt.Fprintf(&buf, lpad(strconv.FormatInt(file.Size(), 10), 12))
		fmt.Fprintf(&buf, file.ModTime().Format(" Jan _2 15:04 "))
		fmt.Fprintf(&buf, "%s\r\n", file.Name())
	}
	return buf.Bytes()
}

func lpad(input string, length int) (result string) {
	if len(input) < length {
		result = strings.Repeat(" ", length-len(input)) + input
	} else if len(input) == length {
		result = input
	} else {
		result = input[0:length]
	}
	return
}
