// Package pop3 provides simple POP3 client.
package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"log"
	"regexp"
	"os"
)

var (
	EOF = errors.New("skip the all mail remaining")
)

// MessageInfo has Number, Size, and Uid fields,
// and used as a return value of ListAll and UidlAll.
// When used as the return value of the method ListAll,
// MessageInfo contain only the Number and Size values.
// When used as the return value of the method UidlAll,
// MessageInfo contain only the Number and Uid values.
type MessageInfo struct {
	Number int
	Size   uint64
	Uid    string
}

// A Client represents a client connection to an POP server.
type Client struct {
	// Text is the pop3.Conn used by the Client.
	Text *Conn
	// keep a reference to the connection so it can be used to create a TLS
	// connection later
	conn net.Conn
}

// Dial returns a new Client connected to an POP server at addr.
// The addr must include a port number.
func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	return NewClient(conn)
}

// NewClient returns a new Client using an existing connection.
func NewClient(conn net.Conn) (*Client, error) {
	text := NewConn(conn)

	_, err := text.ReadResponse()

	if err != nil {
		return nil, err
	}

	return &Client{Text: text, conn: conn}, nil
}

// User issues a USER command to the server using the provided user name.
func (c *Client) User(user string) error {
	return c.cmdSimple("USER %s", user)
}

// Pass issues a PASS command to the server using the provided password.
func (c *Client) Pass(pass string) error {
	return c.cmdSimple("PASS %s", pass)
}

// Stat issues a STAT command to the server
// and returns mail count and total size.
func (c *Client) Stat() (int, uint64, error) {
	return c.cmdStatOrList("STAT", "STAT")
}

// Retr issues a RETR command to the server using the provided mail number
// and returns mail data.
func (c *Client) Retr(number int) (string, error) {
	var err error

	err = c.Text.WriteLine("RETR %d", number)

	if err != nil {
		return "", err
	}

	_, err = c.Text.ReadResponse()

	if err != nil {
		return "", err
	}

	return c.Text.ReadToPeriod()
}

// List issues a LIST command to the server using the provided mail number
// and returns mail number and size.
func (c *Client) List(number int) (string, int, uint64, error) {
//	return c.cmdStatOrList("LIST", "LIST %d", number)
//	return 
	count, size, err := c.cmdStatOrList("LIST", "LIST") 
	filelist, err := c.Text.ReadToPeriod()
	
	return filelist, count, size, err
}

// List issues a LIST command to the server
// and returns array of MessageInfo.
func (c *Client) ListAll() ([]MessageInfo, error) {
	list := make([]MessageInfo, 0)

	err := c.cmdReadLines("LIST", func(line string) error {
		number, size, err := c.convertNumberAndSize(line)

		if err != nil {
			return err
		}

		list = append(list, MessageInfo{Number: number, Size: size})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

// Uidl issues a UIDL command to the server using the provided mail number
// and returns mail number and unique id.
func (c *Client) Uidl(number int) (int, string, error) {
	var err error

	err = c.Text.WriteLine("UIDL %d", number)

	if err != nil {
		return 0, "", err
	}

	var msg string

	msg, err = c.Text.ReadResponse()

	if err != nil {
		return 0, "", err
	}

	var val int
	var uid string

	val, uid, err = c.convertNumberAndUid(msg)

	if err != nil {
		return 0, "", err
	}

	return val, uid, nil
}

// Uidl issues a UIDL command to the server
// and returns array of MessageInfo.
func (c *Client) UidlAll() ([]MessageInfo, error) {
	list := make([]MessageInfo, 0)

	err := c.cmdReadLines("UIDL", func(line string) error {
		number, uid, err := c.convertNumberAndUid(line)

		if err != nil {
			return err
		}

		list = append(list, MessageInfo{Number: number, Uid: uid})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

// Dele issues a DELE command to the server using the provided mail number.
func (c *Client) Dele(number int) error {
	return c.cmdSimple("DELE %d", number)
}

// Noop issues a NOOP command to the server.
func (c *Client) Noop() error {
	return c.cmdSimple("NOOP")
}

// Rset issues a RSET command to the server.
func (c *Client) Rset() error {
	return c.cmdSimple("RSET")
}

// Quit issues a QUIT command to the server.
func (c *Client) Quit() error {
	return c.cmdSimple("QUIT")
}

// ReceiveMail connects to the server at addr,
// and authenticates with user and pass,
// and calling receiveFn for each mail.
func ReceiveMail(addr, user, pass string, receiveFn ReceiveMailFunc) error {
	c, err := Dial(addr)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil && err != EOF {
			c.Rset()
		}

		c.Quit()
		c.Close()
	}()

	if err = c.User(user); err != nil {
		return err
	}

	if err = c.Pass(pass); err != nil {
		return err
	}

	var mis []MessageInfo

	if mis, err = c.UidlAll(); err != nil {
		return err
	}

	for _, mi := range mis {
		var data string

		data, err = c.Retr(mi.Number)

		del, err := receiveFn(mi.Number, mi.Uid, data, err)

		if err != nil && err != EOF {
			return err
		}

		if del {
			if err = c.Dele(mi.Number); err != nil {
				return err
			}
		}

		if err == EOF {
			break
		}
	}

	return nil
}

// ReceiveMailFunc is the type of the function called for each mail.
// Its arguments are mail's number, uid, data, and mail receiving error.
// if this function returns false value, the mail will be deleted,
// if its returns EOF, skip the all mail of remaining.
// (after deleting mail, if necessary)
type ReceiveMailFunc func(number int, uid, data string, err error) (bool, error)

func (c *Client) cmdSimple(format string, args ...interface{}) error {
	var err error

	err = c.Text.WriteLine(format, args...)

	if err != nil {
		return err
	}

	_, err = c.Text.ReadResponse()

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) cmdStatOrList(name, format string, args ...interface{}) (int, uint64, error) {
	var err error

	err = c.Text.WriteLine(format, args...)

	if err != nil {
		return 0, 0, err
	}

	var msg string

	msg, err = c.Text.ReadResponse()

	if err != nil {
		return 0, 0, err
	}

	s := strings.Split(msg, " ")

	if len(s) < 2 {
		return 0, 0, ResponseError(fmt.Sprintf("invalid response format: %s", msg))
	}

	var val int
	var size uint64

	val, size, err = c.convertNumberAndSize(msg)

	if err != nil {
		return 0, 0, err
	}

	return val, size, nil
}

func (c *Client) cmdReadLines(cmnd string, lineFn lineFunc) error {
	var err error

	err = c.Text.WriteLine(cmnd)

	if err != nil {
		return err
	}

	_, err = c.Text.ReadResponse()

	if err != nil {
		return err
	}

	var lines []string

	lines, err = c.Text.ReadLines()

	if err != nil {
		return err
	}

	for _, line := range lines {
		err = lineFn(line)

		if err != nil {
			return err
		}
	}

	return nil
}

type lineFunc func(line string) error

func (c *Client) Close() error {
	return c.Text.Close()
}

func (c *Client) convertNumberAndSize(line string) (int, uint64, error) {
	var err error

	re := regexp.MustCompile("[0-9]+")
	s := re.FindAllString(line, -1)

//	s := strings.Split(numberLine, " ")
	if len(s) < 2 {
		return 0, 0, errors.New(fmt.Sprintf("the length of the array is less than 2: %s", line))
	}

	var val int
	var size uint64

	if val, err = strconv.Atoi(s[0]); err != nil {
		return 0, 0, errors.New(fmt.Sprintf("can not convert element[0] to int type: %s", line))
	}

	if size, err = strconv.ParseUint(s[1], 10, 64); err != nil {
		return 0, 0, errors.New(fmt.Sprintf("can not convert element[1] to uint64 type: %s", line))
	}

	return val, size, nil
}

func (c *Client) convertNumberAndUid(line string) (int, string, error) {
	var err error

	s := strings.Split(line, " ")

	if len(s) < 2 {
		return 0, "", errors.New(fmt.Sprintf("the length of the array is less than 2: %s", line))
	}

	var val int

	if val, err = strconv.Atoi(s[0]); err != nil {
		return 0, "", errors.New(fmt.Sprintf("can not convert element[0] to int type: %s", line))
	}

	return val, s[1], nil
}


func ocPop3Client() {
/*
//	ReceiveMail(addr, user, pass string, receiveFn ReceiveMailFunc) error {
//	if err := pop3.ReceiveMail(address, user, pass,	func(number int, uid, data string, err error) (bool, error) {
	if err := ReceiveMail("127.0.0.1:3110", "a", "b",	func(number int, uid, data string, err error) (bool, error) {
		log.Printf("%d, %s\n", number, uid)

    // implement your own logic here

		return false, nil
	}); err != nil {
	log.Fatalf("%v\n", err)
	}
*/
	var cfg ServerConfigs
	
	cfg = readServerConfigJson("serverconfig.json")

	ip := cfg.Pop3Ip
	port := cfg.Pop3Port
	user := cfg.Pop3User
	pass := cfg.Pop3Password
	
	DOWNLOADED := cfg.Pop3DownloadDir
	
	server := ip + ":" + port

	CreateDirIfNotExist(MYPATH + "/" + DOWNLOADED)


	client, err := Dial(server,)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	defer func() {
		client.Quit()
		client.Close()
	}()

	if err = client.User(user); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	if err = client.Pass(pass); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	var count int
	var size uint64
	var filenames string

	if count, size, err = client.Stat(); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	log.Printf("Count: %d, Size: %d\n", count, size)

//	if count, size, err = client.List(6); err != nil {
	if filenames, count, size, err = client.List(count); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	
	filelist := strings.Split(filenames, "\n")
	log.Printf("Number: %d, Size: %d\n", count, size)
/*
	var mis []MessageInfo
	fmt.Println( "11111 ............................");//jerryprint
	if mis, err = client.ListAll(); err != nil {
		log.Printf("Error: %v\n", err)
	fmt.Println( "22222 ............................");//jerryprint
		return
	}
	fmt.Println( "33333 ............................");//jerryprint

	for _, mi := range mis {
		log.Printf("Number: %d, Size: %d\n", mi.Number, mi.Size)
	}
	fmt.Println( "44444 ............................");//jerryprint
	var number int
	var uid string

	if number, uid, err = client.Uidl(6); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	log.Printf("Number: %d, Uid: %s\n", number, uid)

	if mis, err = client.UidlAll(); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	for _, mi := range mis {
		log.Printf("Number: %d, Uid: %s\n", mi.Number, mi.Uid)
	}
*/

	var content string

	for i := 0; i< count; i++ {
		if content, err = client.Retr(i+1); err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
//		log.Printf("Content:\n%s\n", content)
		temp := strings.Trim(strings.Trim(filelist[i], "\r"), "\n")
		temp  = strings.Trim(strings.Trim(temp, "\n"), "\r")
		fmt.Printf("Get(%04d) : %s\n", i+1, temp)
		
		out, err := os.Create(DOWNLOADED + "/" + temp)//filepath)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		defer out.Close()

		// Write the body to file
		//_, err = io.Copy(out, []byte(content))
		out.WriteString(content)
	}

	if err = client.Quit(); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

/*	
	if content, err = client.Retr(8); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	log.Printf("Content:\n%s\n", content)

	if err = client.Dele(6); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	if err = client.Noop(); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	if err = client.Rset(); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
*/

}

