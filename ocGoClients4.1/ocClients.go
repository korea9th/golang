package main

import (
	"fmt"
	"os"
	"log"
	"encoding/json"
	"io/ioutil"
)


type ServerConfigs struct {
    HttpIp               string `json:"httpip"`
    HttpPort             string `json:"httpport"`
    HttpFileList         string `json:"httpfilelist"`
    HttpDownloadDir      string `json:"httpdownloaddir"`
    Pop3Ip               string `json:"pop3ip"`
    Pop3Port             string `json:"pop3port"`
    Pop3User             string `json:"pop3user"`
    Pop3Password         string `json:"pop3password"`
    Pop3DownloadDir      string `json:"pop3downloaddir"`
    FtpIp                string `json:"ftpip"`
    FtpPort              string `json:"ftpport"`
    FtpId                string `json:"ftpid"`
    FtpPassword          string `json:"ftppassword"`
    FtpDownloadDir       string `json:"ftpdownloaddir"`
    SmtpIp               string `json:"smtpip"`
    SmtpPort             string `json:"smtpport"`
    SmtpFileList         string `json:"smtpfilelist"`
}

var goClientsVersion = "ocClients V4.1(Golang)"

var DOWNLOADED = "downloaded"
var MYPATH = ""

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readServerConfigJson(serverConfigName string) (ServerConfigs){
    // Open our jsonFile
    jsonFile, err := os.Open(serverConfigName)
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }

    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    // read our opened xmlFile as a byte array.
    byteValue, _ := ioutil.ReadAll(jsonFile)

    // we initialize our Users array
//    var users Users
    var cfg ServerConfigs

    // we unmarshal our byteArray which contains our
    // jsonFile's content into 'users' which we defined above
    err = json.Unmarshal(byteValue, &cfg)//&users)
	check(err)
	
	return cfg
}


func CreateDirIfNotExist(dir string) {
	fmt.Println(dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func printHelp(myName string) {
	fmt.Printf("%s\n", goClientsVersion);
	fmt.Printf("Usage: %s [OPTION]\n", myName);
//	fmt.Printf("\t-b configfileslist\tmake hashs by all jsonfiles in configfileslist\n");
	fmt.Printf("\t-h     display this help and exit\n");
	fmt.Printf("\t-http     http client(for download)\n");
	fmt.Printf("\t-ftp      ftp client(for download)\n");
	fmt.Printf("\t-pop3     pop3 client(for download)\n");
	fmt.Printf("\t-smtp     smtp client(for upload)\n");
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	
	MYPATH = dir

//	CreateDirIfNotExist(MYPATH + "/" + DOWNLOADED)

	args := os.Args

	if len(args) == 2 && string(args[1]) == string("-h") {
		printHelp(args[0])
		return
	} else if len(args) == 2 && string(args[1]) == string("-http") {
		ocDownloadFiles()
		return
	} else if len(args) == 2 && string(args[1]) == string("-smtp") {
//		ocSmtpClient()
		ocEmailClient()
		return
	} else if len(args) == 2 && string(args[1]) == string("-pop3") {
		ocPop3Client()
		return
	} else if len(args) == 2 && string(args[1]) == string("-ftp") {
		ocFtpClient()
		return
	} else {
		printHelp(args[0])
		return
	}
	
	
/*

	fileUrl := "http://192.168.100.191:8080/bin/004.jpg"
	err = DownloadFile("004.jpg", fileUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded: " + fileUrl)
*/
}
