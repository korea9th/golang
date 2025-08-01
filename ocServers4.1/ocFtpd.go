// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//package server
package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
)

// Version returns the library version
func Version() string {
	return "0.3.0"
}

// ServerOpts contains parameters for server.NewServer()
type ServerOpts struct {
	// The factory that will be used to create a new FTPDriver instance for
	// each client connection. This is a mandatory option.
	Factory DriverFactory

	Auth Auth

	// Server Name, Default is Go Ftp Server
	Name string

	// The hostname that the FTP server should listen on. Optional, defaults to
	// "::", which means all hostnames on ipv4 and ipv6.
	Hostname string

	// Public IP of the server
	PublicIp string

	// Passive ports
	PassivePorts string

	// The port that the FTP should listen on. Optional, defaults to 3000. In
	// a production environment you will probably want to change this to 21.
	Port int

	// use tls, default is false
	TLS bool

	// if tls used, cert file is required
	CertFile string

	// if tls used, key file is required
	KeyFile string

	// If ture TLS is used in RFC4217 mode
	ExplicitFTPS bool

	WelcomeMessage string

	// A logger implementation, if nil the StdLogger is used
	Logger Logger
}

// Server is the root of your FTP application. You should instantiate one
// of these and call ListenAndServe() to start accepting client connections.
//
// Always use the NewServer() method to create a new Server.
type Server struct {
	*ServerOpts
	listenTo  string
	logger    Logger
	listener  net.Listener
	tlsConfig *tls.Config
	ctx       context.Context
	cancel    context.CancelFunc
	feats     string
}

// ErrServerClosed is returned by ListenAndServe() or Serve() when a shutdown
// was requested.
var ErrServerClosed = errors.New("ftp: Server closed")

// serverOptsWithDefaults copies an ServerOpts struct into a new struct,
// then adds any default values that are missing and returns the new data.
func serverOptsWithDefaults(opts *ServerOpts) *ServerOpts {
	var newOpts ServerOpts
	if opts == nil {
		opts = &ServerOpts{}
	}
	if opts.Hostname == "" {
		newOpts.Hostname = "::"
	} else {
		newOpts.Hostname = opts.Hostname
	}
	if opts.Port == 0 {
		newOpts.Port = 3000
	} else {
		newOpts.Port = opts.Port
	}
	newOpts.Factory = opts.Factory
	if opts.Name == "" {
		newOpts.Name = "Go FTP Server"
	} else {
		newOpts.Name = opts.Name
	}

	if opts.WelcomeMessage == "" {
		newOpts.WelcomeMessage = defaultWelcomeMessage
	} else {
		newOpts.WelcomeMessage = opts.WelcomeMessage
	}

	if opts.Auth != nil {
		newOpts.Auth = opts.Auth
	}

	newOpts.Logger = &StdLogger{}
	if opts.Logger != nil {
		newOpts.Logger = opts.Logger
	}

	newOpts.TLS = opts.TLS
	newOpts.KeyFile = opts.KeyFile
	newOpts.CertFile = opts.CertFile
	newOpts.ExplicitFTPS = opts.ExplicitFTPS

	newOpts.PublicIp = opts.PublicIp
	newOpts.PassivePorts = opts.PassivePorts

	return &newOpts
}

// NewServer initialises a new FTP server. Configuration options are provided
// via an instance of ServerOpts. Calling this function in your code will
// probably look something like this:
//
//     factory := &MyDriverFactory{}
//     server  := server.NewServer(&server.ServerOpts{ Factory: factory })
//
// or:
//
//     factory := &MyDriverFactory{}
//     opts    := &server.ServerOpts{
//       Factory: factory,
//       Port: 2000,
//       Hostname: "127.0.0.1",
//     }
//     server  := server.NewServer(opts)
//
func NewServer(opts *ServerOpts) *Server {
	opts = serverOptsWithDefaults(opts)
	s := new(Server)
	s.ServerOpts = opts
	s.listenTo = net.JoinHostPort(opts.Hostname, strconv.Itoa(opts.Port))
	s.logger = opts.Logger
	return s
}

// NewConn constructs a new object that will handle the FTP protocol over
// an active net.TCPConn. The TCP connection should already be open before
// it is handed to this functions. driver is an instance of FTPDriver that
// will handle all auth and persistence details.
func (server *Server) newConn(tcpConn net.Conn, driver Driver) *Conn {
	c := new(Conn)
	c.namePrefix = "/"
	c.conn = tcpConn
	c.controlReader = bufio.NewReader(tcpConn)
	c.controlWriter = bufio.NewWriter(tcpConn)
	c.driver = driver
	c.auth = server.Auth
	c.server = server
	c.sessionID = newSessionID()
	c.logger = server.logger
	c.tlsConfig = server.tlsConfig

	driver.Init(c)
	return c
}

func simpleTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	config := &tls.Config{}
	if config.NextProtos == nil {
		config.NextProtos = []string{"ftp"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// ListenAndServe asks a new Server to begin accepting client connections. It
// accepts no arguments - all configuration is provided via the NewServer
// function.
//
// If the server fails to start for any reason, an error will be returned. Common
// errors are trying to bind to a privileged port or something else is already
// listening on the same port.
//
func (server *Server) ListenAndServe() error {
	var listener net.Listener
	var err error
	var curFeats = featCmds

	if server.ServerOpts.TLS {
		server.tlsConfig, err = simpleTLSConfig(server.CertFile, server.KeyFile)
		if err != nil {
			return err
		}

		curFeats += " AUTH TLS\n PBSZ\n PROT\n"

		if server.ServerOpts.ExplicitFTPS {
			listener, err = net.Listen("tcp", server.listenTo)
		} else {
			listener, err = tls.Listen("tcp", server.listenTo, server.tlsConfig)
		}
	} else {
		listener, err = net.Listen("tcp", server.listenTo)
	}
	if err != nil {
		return err
	}
	server.feats = fmt.Sprintf(feats, curFeats)

	sessionID := ""
	server.logger.Printf(sessionID, "%s listening on %d", server.Name, server.Port)

	return server.Serve(listener)
}

// Serve accepts connections on a given net.Listener and handles each
// request in a new goroutine.
//
func (server *Server) Serve(l net.Listener) error {
	server.listener = l
	server.ctx, server.cancel = context.WithCancel(context.Background())
	sessionID := ""
	for {
		tcpConn, err := server.listener.Accept()
		if err != nil {
			select {
			case <-server.ctx.Done():
				return ErrServerClosed
			default:
			}
			server.logger.Printf(sessionID, "listening error: %v", err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return err
		}
		driver, err := server.Factory.NewDriver()
		if err != nil {
			server.logger.Printf(sessionID, "Error creating driver, aborting client connection: %v", err)
			tcpConn.Close()
		} else {
			ftpConn := server.newConn(tcpConn, driver)
			go ftpConn.Serve()
		}
	}
}

// Shutdown will gracefully stop a server. Already connected clients will retain their connections
func (server *Server) Shutdown() error {
	if server.cancel != nil {
		server.cancel()
	}
	if server.listener != nil {
		return server.listener.Close()
	}
	// server wasnt even started
	return nil
}

func ocFtpServer() {
//	fmt.Printf("usage : GoFptd -root=c:\\ftpdir -host=192.168.100.123 -port=2121\n")
	
	var cfg ServerConfigs
	
	cfg = readServerConfigJson("serverconfig.json")

	root := cfg.FtpDir
	user := cfg.FtpId
	pass := cfg.FtpPassword
	port := cfg.FtpPort
	host := cfg.FtpIp//"127.0.0.1"

	CreateDirIfNotExist(MYPATH + "/" + root)


	//	factory := &filedriver.FileDriverFactory{
	factory := DriverFactory{
		RootPath: root,
		//		Perm:     server.NewSimplePerm("user", "group"),
		Perm: NewSimplePerm("user", "group"),
	}
	
	intPort, err := strconv.Atoi(port)//ParseInt(port, 10, 32)
	check(err)


	//	opts := &server.ServerOpts{
	opts := &ServerOpts{
		Factory:  factory,
		Port:     intPort,
		Hostname: host,
		//		Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
		Auth: &SimpleAuth{Name: user, Password: pass},
	}

//	fmt.Println("[FTP Server]ListenAndServe Started! -> Port(" + port + ")")
//	fmt.Println("[FTP Server]ListenAndServe Started! -> Username(" + user + ") + Password(" + pass + ")")

	log.Printf("[FTP Server]ListenAndServe Started! -> Port(%s), User/Passwd(%s/%s)", port, user, pass)
//	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
//	log.Printf("Username %v, Password %v", user, pass)
	//	server := server.NewServer(opts)
	server := NewServer(opts)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

