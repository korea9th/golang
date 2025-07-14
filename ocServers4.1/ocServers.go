//20250618 Ocean
package main

//import ("fmt"; "io/ioutil"; "log")
import (
	"fmt"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
//	"syscall"
//	"os/exec"
//	"unicode/utf8"
//	"golang.org/x/text/encoding/charmap"
)

var ocServersVersion = "ocGoServers V4.1(Golang)"


var MYPATH = ""

	
func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Configs struct {
    EvaluationFacility    string `json:"evaluationfacility"`
    ReceiptNumber         string `json:"receiptnumber"`
    DeveloperName         string `json:"developername"`
    ProductName           string `json:"productname"`
    Description           string `json:"description"`
    TargetDir             string `json:"targetdir"`
    HashedAlgorithm       string `json:"hashedalgorithm"`
}


type ServerConfigs struct {
    HttpIp               string `json:"httpip"`
    HttpPort             string `json:"httpport"`
    HttpDir              string `json:"httpdir"`
    Pop3Ip              string `json:"pop3ip"`
    Pop3Port             string `json:"pop3port"`
    Pop3Dir              string `json:"pop3dir"`
    FtpIp                string `json:"ftpip"`
    FtpPort              string `json:"ftpport"`
    FtpId                string `json:"ftpid"`
    FtpPassword          string `json:"ftppassword"`
    FtpDir               string `json:"ftpdir"`
    SmtpIp              string `json:"smtpip"`
    SmtpPort             string `json:"smtpport"`
    SmtpDir              string `json:"smtpdir"`
}


func printHelp(myName string) {
	fmt.Printf("%s\n", ocServersVersion);
	fmt.Printf("Usage: %s [OPTION]\n", myName);
//	fmt.Printf("\t-b configfileslist\tmake hashs by all jsonfiles in configfileslist\n");
	fmt.Printf("\t-h     display this help and exit\n");
	fmt.Printf("\t-http     http file server(for download)\n");
	fmt.Printf("\t-ftp      ftp server(for download)\n");
	fmt.Printf("\t-pop3     pop3 server(for download)\n");
	fmt.Printf("\t-smtp     smtp server(for upload)\n");
	fmt.Printf("\t-hash     http server(for hash)\n");
	fmt.Printf("\t-cmd      make hash\n");
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

func readConfigJson(jsonFileName string) (Configs){
    // Open our jsonFile
    jsonFile, err := os.Open(jsonFileName)
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
    var cfg Configs

    // we unmarshal our byteArray which contains our
    // jsonFile's content into 'users' which we defined above
    json.Unmarshal(byteValue, &cfg)//&users)
	
	return cfg
}

func writeConfigJson(jsonFileName string, config Configs) {
/*
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	b, _ := json.MarshalIndent(config, "", "	")//&users)

	// Open our jsonFile
	jsonFile, err := os.Open(jsonFileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	
	jsonFile.Write(b)
	jsonFile.Sync()
	fmt.Println(string(b))
*/	
	b, _ := json.MarshalIndent(config, "", "	")//&users)
//	b, _ := json.Marshal(config)

/*
	err := ioutil.WriteFile(jsonFileName, b, os.FileMode(0644)) // articles.json 파일에 JSON 문서 저장
	if err != nil {
		fmt.Println(err)
		return
	}
*/
	f1, err := os.Create(jsonFileName)
	check(err)
	defer f1.Close()
	fmt.Fprintf(f1, string(b))
}


func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	
	MYPATH = dir

	CreateDirIfNotExist(MYPATH + "/" + RESULT)

	args := os.Args
	

	if len(args) == 3 && string(args[1]) == string("-b") {
		readBatch(args[2])
		return
	} else if len(args) == 2 && string(args[1]) == string("-h") {
		printHelp(args[0])
		return
	} else if len(args) == 2 && string(args[1]) == string("-hash") {
		ocHttpd()
		return
	} else if len(args) == 2 && string(args[1]) == string("-http") {
		ocHttpFileServer()
		return
	} else if len(args) == 2 && string(args[1]) == string("-pop3") {
		ocPop3Server()
		return
	} else if len(args) == 2 && string(args[1]) == string("-ftp") {
		ocFtpServer()
		return
	} else if len(args) == 2 && string(args[1]) == string("-smtp") {
		ocSmtpServer()
		return
	} else if len(args) == 2 && string(args[1]) == string("-cmd") {
		ocGoHash()
		return
	} else {
		printHelp(args[0])
		return
	}
}

