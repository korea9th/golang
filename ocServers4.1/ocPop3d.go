package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"log"
	//"launchpad.net/mgo"
	//	"launchpad.net/mgo/bson"
	//	"launchpad.net/mgo/bson"
)

type Message struct {
	Headers string
	Message string
	Subject string
}
type User struct {
	Username string
	Password string
	Messages []Message
}

type Message_head struct {
	Id   int
	Size int
}

var (
//	session    *mgo.Session
//	collection *mgo.Collection
)
var (
	eol = "\r\n"
)

const (
	host = "experimental.zapto.org"
)

const (
	STATE_UNAUTHORIZED = 1
	STATE_TRANSACTION  = 2
	STATE_UPDATE       = 3
)

var (
	serviceDir = ""
)

func ocPop3Server() {
//	fmt.Printf("usage : GoPop3d -dir=c:\\emls -ip=192.168.100.123 -port=3110\n")

	//	var emls emlFiles
	var cfg ServerConfigs
	
	cfg = readServerConfigJson("serverconfig.json")

	port := cfg.Pop3Port

	ip := cfg.Pop3Ip//"127.0.0.1"
	serviceDir = cfg.Pop3Dir//"receivedFiles"

	CreateDirIfNotExist(MYPATH + "/" + serviceDir)

	//	emls = searchFiles(*dir)

	//	connectDatabase()
	// Server IP and PORT
	//	service := "192.168.1.4:3110"
	service := ip + ":" + port                       //"127.0.0.1:3110"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service) //"ip4", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error.. %s", err.Error())
		fmt.Printf("\nERRORRRRRRRRRR \n")
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error.. %s", err.Error())
	}

	log.Printf("[Pop3 Server]ListenAndServe Started! -> Port(%s)", port)
//	fmt.Fprintf(os.Stdout, "Server listening, host: %s\n", tcpAddr.String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error.. %s", err.Error())
			continue
		}
		// run as goroutine
		go handleClient(conn)
	}

	//	fmt.Println(tcpAddr)
	//	fmt.Println("hejsan")
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("\n handleClient \n")

	var emls emlFiles

	var (
		tmp_user = ""
		eol      = "\r\n"
		state    = 1
	)
	reader := bufio.NewReader(conn)
	//	writer := bufio.NewWriter(conn)

	// State
	// 1 = Unauthorized
	// 2 = Transaction mode
	// 3 = update mode

	// First welcome the new connection
	fmt.Fprintf(conn, "+OK simple POP3 server %s powered by Go"+eol, host)
	//nr, _ := writer.WriteString("+OK simple POP3 server (" + host + ") powered by Go")

	//writer.Flush()	// Måste använda flush eftersom inte buffern  fylls med ett meddelande...

	for {
		// Reads a line from the client
		raw_line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error!!" + err.Error())
			return
		}

		// Parses the command
		cmd, args := getCommand(raw_line)

		fmt.Println(">" + cmd + "<")
		arg, _ := getSafeArg(args, 0)
		fmt.Println("line>" + arg + "<")
		fmt.Println(raw_line)

		if cmd == "USER" && state == STATE_UNAUTHORIZED {
			tmp_user, _ = getSafeArg(args, 0)
			if userExists(tmp_user) {
				fmt.Fprintf(conn, "+OK"+eol)
				fmt.Println(">+OK" + eol)
			} else {
				fmt.Fprintf(conn, "-ERR The user %s doesn't belong here!"+eol, tmp_user)
				fmt.Println(">-ERR The user " + tmp_user + " doesn't belong here!" + eol)
			}
		} else if cmd == "PASS" && state == STATE_UNAUTHORIZED {
			pass, _ := getSafeArg(args, 0)
			if authUser(tmp_user, pass) {
				fmt.Fprintf(conn, "+OK User signed in"+eol)
				fmt.Println(">+OK User signed in" + eol)

				state = 2

			} else {
				fmt.Fprintf(conn, "-ERR Password incorrect!"+eol)
				fmt.Println(">-ERR Password incorrect!" + eol)
			}
		} else if cmd == "STAT" && state == STATE_TRANSACTION {
			fmt.Println("STAT accepted")

			if emls.Filecounts == 0 {
				emls = searchFiles(serviceDir)
			}

			fmt.Fprintf(conn, "+OK "+strconv.Itoa(emls.Filecounts)+" "+strconv.Itoa(int(emls.Filesizes))+eol)
			fmt.Println(">+OK " + strconv.Itoa(emls.Filecounts) + " " + strconv.Itoa(int(emls.Filesizes)) + eol)

			/*
				nr_messages, size_messages := getStat(tmp_user)
				fmt.Fprintf(conn, "+OK "+strconv.Itoa(nr_messages)+" "+strconv.Itoa(size_messages)+eol)
				fmt.Println(">+OK " + strconv.Itoa(nr_messages) + " " + strconv.Itoa(size_messages) + eol)
			*/
		} else if cmd == "LIST" && state == STATE_TRANSACTION {
			fmt.Println("List accepted")

			if emls.Filecounts == 0 {
				emls = searchFiles(serviceDir)
			}

			fmt.Fprintf(conn, "+OK %d messages (%d octets)\r\n%s\r\n.\r\n", emls.Filecounts, len(strings.Join(emls.Filenames, "")), strings.Join(emls.Filenames, "\r\n"))
			fmt.Printf(">+OK %d messages (%d octets)\r\n%s\r\n.", emls.Filecounts, len(strings.Join(emls.Filenames, "")), strings.Join(emls.Filenames, "\r\n"))
			//			fmt.Fprintf(conn, "+OK %d messages (%d octets)\r\n%s", emls.Filecounts, emls.Filesizes, strings.Join(emls.Filenames, "\r\n"))
			//			fmt.Printf("+OK %d messages (%d octets)\r\n%s", emls.Filecounts, emls.Filesizes, strings.Join(emls.Filenames, "\r\n"))
			// Ending
			//			fmt.Fprintf(conn, ".\r\n")

			//			return "+OK %d messages (%i octets)\r\n%s\r\n." % (i, size, '\r\n'.join(names))

			/*
				filecount := 0
				var filesize int64
				filesize = 0
				filenames := ""

				path, err := os.Getwd()
				if err != nil {
					fmt.Printf("listUpload error [%v]\n", err)
					return
				} else {
					//		fmt.Printf(path)
				}
				err = filepath.Walk(path+"/"+serviceDir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}

					//		fmt.Fprintf(fo, "%s\n", info.Name())
					fmt.Printf("%s %d\n", info.Name(), info.Size())
					filenames = filenames + info.Name() + "\r\n"
					filecount++
					filesize = filesize + info.Size()
					return nil
				})
				if err != nil {
					fmt.Printf("walk error [%v]\n", err)
				}

				fmt.Fprintf(conn, "+OK %d messages (%d octets)\r\n%s", filecount, filesize, filenames)
				// Ending
				fmt.Fprintf(conn, ".\r\n")
			*/
			/*
				nr, tot_size, Message_head := getList(tmp_user)
				fmt.Fprintf(conn, "+OK %d messages (%d octets)\r\n", nr, tot_size)
				// Print all messages
				for _, val := range Message_head {
					fmt.Fprintf(conn, "%d %d\r\n", val.Id, val.Size)
				}
			*/

		} else if cmd == "UIDL" && state == STATE_TRANSACTION {

			// Retreive one message but don't delete it from the server..
			//message, size, _ := getMessage(tmp_user, 1)
			//fmt.Fprintf(conn, "+OK " + strconv.Itoa(size) + " octets" + eol)
			//fmt.Fprintf(conn, message.message + eol + "." + eol)
			fmt.Fprintf(conn, "-ERR Command not implemented"+eol)

		} else if cmd == "TOP" && state == STATE_TRANSACTION {
			arg, _ := getSafeArg(args, 0)
			nr, _ := strconv.Atoi(arg)
			headers := getTop(tmp_user, nr)

			fmt.Fprintf(conn, "+OK Top message followes"+eol)
			fmt.Fprintf(conn, headers+eol+eol+"."+eol)

		} else if cmd == "RETR" && state == STATE_TRANSACTION {
			fmt.Println("RETR accepted")
			if emls.Filecounts == 0 {
				emls = searchFiles(serviceDir)
			}

			fileinfo := emls.Fileinfo

			arg, _ := getSafeArg(args, 0)
			nr, _ := strconv.Atoi(arg)

			nr--

			path, err := os.Getwd()
			if err != nil {
				fmt.Printf("listUpload error [%v]\n", err)
				//		return nil
			} else {
				//		fmt.Printf(path)
			}

			fmt.Printf("filename : %s\n", path+"/"+serviceDir+"/"+fileinfo[nr].Name)

			body, err := ioutil.ReadFile(path + "/" + serviceDir + "/" + fileinfo[nr].Name)
			check(err)

			str := string(body)
			bodylen := len(str)

			fmt.Fprintf(conn, "+OK "+strconv.Itoa(bodylen)+" octets\r\n"+str+"\r\n.\r\n")
			fmt.Printf("+OK " + strconv.Itoa(bodylen) + " octets\r\n") // + str + "\r\n.")

		} else if cmd == "DELE" && state == STATE_TRANSACTION {
			arg, _ := getSafeArg(args, 0)
			nr, _ := strconv.Atoi(arg)
			deleteMessage(tmp_user, nr)
			fmt.Fprintf(conn, "+OK"+eol)
		} else if cmd == "QUIT" {
			fmt.Fprintf(conn, "+OK pypopper POP3 server signing off"+eol)
			fmt.Printf("+OK pypopper POP3 server signing off" + eol)
			return
		}
	}
}

// cuts the line into command and arguments
func getCommand(line string) (string, []string) {
	line = strings.Trim(line, "\r \n")
	cmd := strings.Split(line, " ")
	return cmd[0], cmd[1:]
}
func getSafeArg(args []string, nr int) (string, error) {
	if nr < len(args) {
		return args[nr], nil
	}
	return "", errors.New("Out of range")
}

// Checks if a user exists or not
func userExists(user string) bool {
	//	cnt, _ := collection.Find(bson.M{"username": user}).Count()
	//	return cnt == 1
	return true
}

// Checks if the login info is correct
func authUser(user string, pass string) bool {
	//	cnt, _ := collection.Find(bson.M{"username": user, "password": pass}).Count()
	//	return cnt == 1
	return true
}

// Returns the number of messages in the maildrop and
// the size in bytes of the maildrop
func getStat(user string) (int, int) {

	var filesize int64
	filecount := 0

	filesize = 0

	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("listUpload error [%v]\n", err)
		return 1, 1
	} else {
		//		fmt.Printf(path)
	}
	err = filepath.Walk(path+"/receivedFiles", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		//		fmt.Fprintf(fo, "%s\n", info.Name())
		fmt.Printf("%s %d\n", info.Name(), info.Size())
		filecount++
		filesize = filesize + info.Size()
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	/*
	       for i in range(0, len(filenames)):
	           filesize = filesize + os.path.getsize(filenames[i])
	   #        print(filenames[i])
	   	return "+OK %i %i" % (i+1, filesize)
	*/
	/*test := User{ Username: "test1223", Password: "pass1223" }
	errore := collection.Insert(&test)
	if errore != nil {
		fmt.Println(errore.Error())
	}*/
	result := User{}
	/*
		err := collection.Find(bson.M{"username": user}).One(&result)
		if err != nil {
			fmt.Println("Error:" + err.Error())
			// If the document cannot be found, error occurs here
		}
	*/
	fmt.Println(result)
	fmt.Println(len(result.Messages))
	// To count the total octets..
	var (
		sum = 0
	)
	// Count how many letters there are in all the headers and messages
	for _, v := range result.Messages {
		sum = sum + len(v.Headers) + len(v.Message) // headers_cnt+message_cnt
	}

	// return the count and the size in octets (bytes)
	return filecount, int(filesize) //len(result.Messages), sum * 8
}

// Returns all message heads in maildrop if no argument
// or the message head for the mail id
// in format:
//
func getList(user string) (int, int, []Message_head) {
	result := User{}
	/*
		err := collection.Find(bson.M{"username": user}).One(&result)
		if err != nil {
			fmt.Println("Error:" + err.Error())
			// If the document cannot be found, error occurs here
		}
	*/

	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("listUpload error [%v]\n", err)
	} else {
		//		fmt.Printf(path)
	}
	err = filepath.Walk(path+"/receivedFiles", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		//		fmt.Fprintf(fo, "%s\n", info.Name())
		fmt.Printf("%s\n", info.Name())
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	fmt.Println(result)
	fmt.Println(len(result.Messages))
	// To count the total octets..
	var (
		sum      = 0
		messages []Message_head
	)
	// Add all messages into a header struct
	for i, v := range result.Messages {
		size := (len(v.Headers) + len(v.Message)) * 8
		m := Message_head{Id: i + 1,
			Size: size}
		messages = append(messages, m)
		sum = sum + len(v.Headers) + len(v.Message) // headers_cnt+message_cnt
	}
	fmt.Println(messages)
	return len(messages), sum * 8, messages
	/*
		m := Message_head{
			Id: 1,
			Size: 180}
		messages := []Message_head{ m }
		return 1, 180, messages*/
}

//func getListN(
// Returns the message of the id
func getMessage(user string, id int) (Message, int, error) {

	result := User{}
	/*
		err := collection.Find(bson.M{"username": user}).One(&result)
		if err != nil {
			fmt.Println("Error:" + err.Error())
			// If the document cannot be found, error occurs here
		}
	*/
	fmt.Println(result)
	fmt.Println(len(result.Messages))

	// Get the specified message
	i := id - 1

	size := len(result.Messages[i].Message) * 8
	return result.Messages[i], size, nil

	/*
			message := Message{
				Headers: "From: test@test.se"+ eol +
		"Subject: This is the subject"+ eol +
		"To: pr_125@hotmail.com"+ eol +
		"X-PHP-Originating-Script: 0:smtp-server-test.php",
				Subject: "This is the subject",
				Message: "This is the message"}
			return message, len(message.Message)*8, nil
	*/
}

//						TO-DO:	int arg1, int arg2
func getTop(user string, id int) string {
	result := User{}
	/*
		err := collection.Find(bson.M{"username": user}).One(&result)
		if err != nil {
			fmt.Println("Error:" + err.Error())
			// If the document cannot be found, error occurs here
		}
	*/
	// Get the specified message
	return result.Messages[id-1].Headers
	/*
			m := Message {
				Headers: "From: test@test.se"+ eol +
		"Subject: This is the subject"+ eol +
		"To: pr_125@hotmail.com"+ eol +
		"X-PHP-Originating-Script: 0:smtp-server-test.php",
				Message: "This is the message"}
			// Depending on the arg1 and arg2 it returns different
			// 1 0 means return header but no body
			return m.Headers*/
}

// Marks a message for deletion, the message is removed
// in the UPDATE-state
func deleteMessage(user string, id int) {
	/*
		i := id - 1
		err := collection.Update(bson.M{"username": user}, bson.M{"$unset": bson.M{"messages": i}})
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
	*/
}

type fileInfo struct {
	Name       string
	Size       int64
	Filenumber int
}

type emlFiles struct {
	Filecounts int
	Filesizes  int64
	Filenames  []string
	Fileinfo   []fileInfo
}

func searchFiles(dir string) emlFiles {

	var a emlFiles
	var finfo fileInfo
	filecount := 0
	var filesize int64
	filesize = 0

	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("listUpload error [%v]\n", err)
		//		return nil
	} else {
		//		fmt.Printf(path)
	}
	err = filepath.Walk(path+"/"+serviceDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		//		fmt.Fprintf(fo, "%s\n", info.Name())
		//		fmt.Printf("%s %d\n", info.Name(), info.Size())
		filecount++
		filesize = filesize + info.Size()

		finfo.Name = info.Name()
		finfo.Size = info.Size()
		a.Filenames = append(a.Filenames, finfo.Name)
		finfo.Filenumber = filecount

		a.Fileinfo = append(a.Fileinfo, finfo)
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	a.Filecounts = filecount
	a.Filesizes = filesize

	return a
}
