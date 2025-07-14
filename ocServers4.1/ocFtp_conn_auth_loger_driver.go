// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"crypto/subtle"
	"errors"
	"os"
)

const (
	defaultWelcomeMessage = "Welcome to the Go FTP Server"
)

type Conn struct {
	conn          net.Conn
	controlReader *bufio.Reader
	controlWriter *bufio.Writer
	dataConn      DataSocket
	driver        Driver
	auth          Auth
	logger        Logger
	server        *Server
	tlsConfig     *tls.Config
	sessionID     string
	namePrefix    string
	reqUser       string
	user          string
	renameFrom    string
	lastFilePos   int64
	appendData    bool
	closed        bool
	tls           bool
}

func (conn *Conn) LoginUser() string {
	return conn.user
}

func (conn *Conn) IsLogin() bool {
	return len(conn.user) > 0
}

func (conn *Conn) PublicIp() string {
	return conn.server.PublicIp
}

func (conn *Conn) passiveListenIP() string {
	var listenIP string
	if len(conn.PublicIp()) > 0 {
		listenIP = conn.PublicIp()
	} else {
		listenIP = conn.conn.LocalAddr().(*net.TCPAddr).IP.String()
	}

	lastIdx := strings.LastIndex(listenIP, ":")
	if lastIdx <= 0 {
		return listenIP
	}
	return listenIP[:lastIdx]
}

func (conn *Conn) PassivePort() int {
	if len(conn.server.PassivePorts) > 0 {
		portRange := strings.Split(conn.server.PassivePorts, "-")

		if len(portRange) != 2 {
			log.Println("empty port")
			return 0
		}

		minPort, _ := strconv.Atoi(strings.TrimSpace(portRange[0]))
		maxPort, _ := strconv.Atoi(strings.TrimSpace(portRange[1]))

		return minPort + mrand.Intn(maxPort-minPort)
	}
	// let system automatically chose one port
	return 0
}

// returns a random 20 char string that can be used as a unique session ID
func newSessionID() string {
	hash := sha256.New()
	_, err := io.CopyN(hash, rand.Reader, 50)
	if err != nil {
		return "????????????????????"
	}
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr[0:20]
}

// Serve starts an endless loop that reads FTP commands from the client and
// responds appropriately. terminated is a channel that will receive a true
// message when the connection closes. This loop will be running inside a
// goroutine, so use this channel to be notified when the connection can be
// cleaned up.
func (conn *Conn) Serve() {
	conn.logger.Print(conn.sessionID, "Connection Established")
	// send welcome
	conn.writeMessage(220, conn.server.WelcomeMessage)
	// read commands
	for {
		line, err := conn.controlReader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				conn.logger.Print(conn.sessionID, fmt.Sprint("read error:", err))
			}

			break
		}
		conn.receiveLine(line)
		// QUIT command closes connection, break to avoid error on reading from
		// closed socket
		if conn.closed == true {
			break
		}
	}
	conn.Close()
	conn.logger.Print(conn.sessionID, "Connection Terminated")
}

// Close will manually close this connection, even if the client isn't ready.
func (conn *Conn) Close() {
	conn.conn.Close()
	conn.closed = true
	if conn.dataConn != nil {
		conn.dataConn.Close()
		conn.dataConn = nil
	}
}

func (conn *Conn) upgradeToTLS() error {
	conn.logger.Print(conn.sessionID, "Upgrading connectiion to TLS")
	tlsConn := tls.Server(conn.conn, conn.tlsConfig)
	err := tlsConn.Handshake()
	if err == nil {
		conn.conn = tlsConn
		conn.controlReader = bufio.NewReader(tlsConn)
		conn.controlWriter = bufio.NewWriter(tlsConn)
		conn.tls = true
	}
	return err
}

// receiveLine accepts a single line FTP command and co-ordinates an
// appropriate response.
func (conn *Conn) receiveLine(line string) {
	command, param := conn.parseLine(line)
	conn.logger.PrintCommand(conn.sessionID, command, param)
	cmdObj := commands[strings.ToUpper(command)]
	if cmdObj == nil {
		conn.writeMessage(500, "Command not found")
		return
	}
	if cmdObj.RequireParam() && param == "" {
		conn.writeMessage(553, "action aborted, required param missing")
	} else if cmdObj.RequireAuth() && conn.user == "" {
		conn.writeMessage(530, "not logged in")
	} else {
		cmdObj.Execute(conn, param)
	}
}

func (conn *Conn) parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}

// writeMessage will send a standard FTP response back to the client.
func (conn *Conn) writeMessage(code int, message string) (wrote int, err error) {
	conn.logger.PrintResponse(conn.sessionID, code, message)
	line := fmt.Sprintf("%d %s\r\n", code, message)
	wrote, err = conn.controlWriter.WriteString(line)
	conn.controlWriter.Flush()
	return
}

// writeMessage will send a standard FTP response back to the client.
func (conn *Conn) writeMessageMultiline(code int, message string) (wrote int, err error) {
	conn.logger.PrintResponse(conn.sessionID, code, message)
	line := fmt.Sprintf("%d-%s\r\n%d END\r\n", code, message, code)
	wrote, err = conn.controlWriter.WriteString(line)
	conn.controlWriter.Flush()
	return
}

// buildPath takes a client supplied path or filename and generates a safe
// absolute path within their account sandbox.
//
//    buildpath("/")
//    => "/"
//    buildpath("one.txt")
//    => "/one.txt"
//    buildpath("/files/two.txt")
//    => "/files/two.txt"
//    buildpath("files/two.txt")
//    => "/files/two.txt"
//    buildpath("/../../../../etc/passwd")
//    => "/etc/passwd"
//
// The driver implementation is responsible for deciding how to treat this path.
// Obviously they MUST NOT just read the path off disk. The probably want to
// prefix the path with something to scope the users access to a sandbox.
func (conn *Conn) buildPath(filename string) (fullPath string) {
	if len(filename) > 0 && filename[0:1] == "/" {
		fullPath = filepath.Clean(filename)
	} else if len(filename) > 0 && filename != "-a" {
		fullPath = filepath.Clean(conn.namePrefix + "/" + filename)
	} else {
		fullPath = filepath.Clean(conn.namePrefix)
	}
	fullPath = strings.Replace(fullPath, "//", "/", -1)
	fullPath = strings.Replace(fullPath, string(filepath.Separator), "/", -1)
	return
}

// sendOutofbandData will send a string to the client via the currently open
// data socket. Assumes the socket is open and ready to be used.
func (conn *Conn) sendOutofbandData(data []byte) {
	bytes := len(data)
	if conn.dataConn != nil {
		conn.dataConn.Write(data)
		conn.dataConn.Close()
		conn.dataConn = nil
	}
	message := "Closing data connection, sent " + strconv.Itoa(bytes) + " bytes"
	conn.writeMessage(226, message)
}

func (conn *Conn) sendOutofBandDataWriter(data io.ReadCloser) error {
	conn.lastFilePos = 0
	bytes, err := io.Copy(conn.dataConn, data)
	if err != nil {
		conn.dataConn.Close()
		conn.dataConn = nil
		return err
	}
	message := "Closing data connection, sent " + strconv.Itoa(int(bytes)) + " bytes"
	conn.writeMessage(226, message)
	conn.dataConn.Close()
	conn.dataConn = nil

	return nil
}


// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(string, string) (bool, error)
}

var (
	_ Auth = &SimpleAuth{}
)

// SimpleAuth implements Auth interface to provide a memory user login auth
type SimpleAuth struct {
	Name     string
	Password string
}

// CheckPasswd will check user's password
func (a *SimpleAuth) CheckPasswd(name, pass string) (bool, error) {
	return constantTimeEquals(name, a.Name) && constantTimeEquals(pass, a.Password), nil
}

func constantTimeEquals(a, b string) bool {
	return len(a) == len(b) && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

type Logger interface {
	Print(sessionId string, message interface{})
	Printf(sessionId string, format string, v ...interface{})
	PrintCommand(sessionId string, command string, params string)
	PrintResponse(sessionId string, code int, message string)
}

// Use an instance of this to log in a standard format
type StdLogger struct{}

func (logger *StdLogger) Print(sessionId string, message interface{}) {
	log.Printf("%s  %s", sessionId, message)
}

func (logger *StdLogger) Printf(sessionId string, format string, v ...interface{}) {
	logger.Print(sessionId, fmt.Sprintf(format, v...))
}

func (logger *StdLogger) PrintCommand(sessionId string, command string, params string) {
	if command == "PASS" {
		log.Printf("%s > PASS ****", sessionId)
	} else {
		log.Printf("%s > %s %s", sessionId, command, params)
	}
}

func (logger *StdLogger) PrintResponse(sessionId string, code int, message string) {
	log.Printf("%s < %d %s", sessionId, code, message)
}

// Silent logger, produces no output
type DiscardLogger struct{}

func (logger *DiscardLogger) Print(sessionId string, message interface{})                  {}
func (logger *DiscardLogger) Printf(sessionId string, format string, v ...interface{})     {}
func (logger *DiscardLogger) PrintCommand(sessionId string, command string, params string) {}
func (logger *DiscardLogger) PrintResponse(sessionId string, code int, message string)     {}




type FileDriver struct {
	RootPath string
	//	server.Perm
	Perm
}

type FileInfo struct {
	os.FileInfo

	mode  os.FileMode
	owner string
	group string
}

type Driver interface {
	// Init init
	Init(*Conn)

	// params  - a file path
	// returns - a time indicating when the requested path was last modified
	//         - an error if the file doesn't exist or the user lacks
	//           permissions
	Stat(string) (*FileInfo, error)

	// params  - path
	// returns - true if the current user is permitted to change to the
	//           requested path
	ChangeDir(string) error

	// params  - path, function on file or subdir found
	// returns - error
	//           path
	ListDir(string, func(*FileInfo) error) error

	// params  - path
	// returns - nil if the directory was deleted or any error encountered
	DeleteDir(string) error

	// params  - path
	// returns - nil if the file was deleted or any error encountered
	DeleteFile(string) error

	// params  - from_path, to_path
	// returns - nil if the file was renamed or any error encountered
	Rename(string, string) error

	// params  - path
	// returns - nil if the new directory was created or any error encountered
	MakeDir(string) error

	// params  - path
	// returns - a string containing the file data to send to the client
	GetFile(string, int64) (int64, io.ReadCloser, error)

	// params  - destination path, an io.Reader containing the file data
	// returns - the number of bytes writen and the first error encountered while writing, if any.
	PutFile(string, io.Reader, bool) (int64, error)
}

func (f *FileInfo) Mode() os.FileMode {
	return f.mode
}

func (f *FileInfo) Owner() string {
	return f.owner
}

func (f *FileInfo) Group() string {
	return f.group
}

func (driver *FileDriver) realPath(path string) string {
	paths := strings.Split(path, "/")
	return filepath.Join(append([]string{driver.RootPath}, paths...)...)
}

//func (driver *FileDriver) Init(conn *server.Conn) {
func (driver *FileDriver) Init(conn *Conn) {
	//driver.conn = conn
}

func (driver *FileDriver) ChangeDir(path string) error {
	rPath := driver.realPath(path)
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	return errors.New("Not a directory")
}

//func (driver *FileDriver) Stat(path string) (server.FileInfo, error) {
func (driver *FileDriver) Stat(path string) (*FileInfo, error) {
	basepath := driver.realPath(path)
	rPath, err := filepath.Abs(basepath)
	if err != nil {
		return nil, err
	}
	f, err := os.Lstat(rPath)
	if err != nil {
		return nil, err
	}
	mode, err := driver.Perm.GetMode(path)
	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		mode |= os.ModeDir
	}
	owner, err := driver.Perm.GetOwner(path)
	if err != nil {
		return nil, err
	}
	group, err := driver.Perm.GetGroup(path)
	if err != nil {
		return nil, err
	}
	return &FileInfo{f, mode, owner, group}, nil
}

//func (driver *FileDriver) ListDir(path string, callback func(server.FileInfo) error) error {
func (driver *FileDriver) ListDir(path string, callback func(*FileInfo) error) error {
	fmt.Printf("//////////////////////////////////////////////////////////////////////////////\n")
	basepath := driver.realPath(path)
	fmt.Printf(basepath + "\n") //////////////////////////////////////////////////////////////////////////////
	return filepath.Walk(basepath, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rPath, _ := filepath.Rel(basepath, f)
		if rPath == info.Name() {
			mode, err := driver.Perm.GetMode(rPath)
			if err != nil {
				return err
			}
			if info.IsDir() {
				mode |= os.ModeDir
			}
			owner, err := driver.Perm.GetOwner(rPath)
			if err != nil {
				return err
			}
			group, err := driver.Perm.GetGroup(rPath)
			if err != nil {
				return err
			}
			//			err = callback(&FileInfo{info, mode, owner, group})
			err = callback(&FileInfo{info, mode, owner, group})
			if err != nil {
				return err
			}
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
}

func (driver *FileDriver) DeleteDir(path string) error {
	rPath := driver.realPath(path)
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return os.Remove(rPath)
	}
	return errors.New("Not a directory")
}

func (driver *FileDriver) DeleteFile(path string) error {
	rPath := driver.realPath(path)
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return os.Remove(rPath)
	}
	return errors.New("Not a file")
}

func (driver *FileDriver) Rename(fromPath string, toPath string) error {
	oldPath := driver.realPath(fromPath)
	newPath := driver.realPath(toPath)
	return os.Rename(oldPath, newPath)
}

func (driver *FileDriver) MakeDir(path string) error {
	rPath := driver.realPath(path)
	return os.MkdirAll(rPath, os.ModePerm)
}

func (driver *FileDriver) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	rPath := driver.realPath(path)
	f, err := os.Open(rPath)
	if err != nil {
		return 0, nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return 0, nil, err
	}

	f.Seek(offset, os.SEEK_SET)

	return info.Size(), f, nil
}

func (driver *FileDriver) PutFile(destPath string, data io.Reader, appendData bool) (int64, error) {
	rPath := driver.realPath(destPath)
	var isExist bool
	f, err := os.Lstat(rPath)
	if err == nil {
		isExist = true
		if f.IsDir() {
			return 0, errors.New("A dir has the same name")
		}
	} else {
		if os.IsNotExist(err) {
			isExist = false
		} else {
			return 0, errors.New(fmt.Sprintln("Put File error:", err))
		}
	}

	if appendData && !isExist {
		appendData = false
	}

	if !appendData {
		if isExist {
			err = os.Remove(rPath)
			if err != nil {
				return 0, err
			}
		}
		f, err := os.Create(rPath)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		bytes, err := io.Copy(f, data)
		if err != nil {
			return 0, err
		}
		return bytes, nil
	}

	of, err := os.OpenFile(rPath, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return 0, err
	}
	defer of.Close()

	_, err = of.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(of, data)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}

//type FileDriverFactory struct {
type DriverFactory struct {
	RootPath string
	//	server.Perm
	Perm
}

//func (factory *FileDriverFactory) NewDriver() (server.Driver, error) {
func (factory *DriverFactory) NewDriver() (Driver, error) {
	return &FileDriver{factory.RootPath, factory.Perm}, nil
}
