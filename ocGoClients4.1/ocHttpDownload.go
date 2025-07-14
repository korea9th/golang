package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"bufio"
)

func ocDownloadFiles() {
	var cfg ServerConfigs
	
	cfg = readServerConfigJson("serverconfig.json")

	ip := cfg.HttpIp
	port := cfg.HttpPort
	filelist := cfg.HttpFileList
	DOWNLOADED = cfg.HttpDownloadDir
	
	CreateDirIfNotExist(MYPATH + "/" + DOWNLOADED)


	fo, err := os.Open(filelist)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	targetFile := ""

	reader := bufio.NewReader(fo)
	i := 0
	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix || err != nil {
			break
		}
		targetFile = string(line)
		err = DownloadFile(targetFile, "http://" + ip + ":" + port + "/" + targetFile)
		if err != nil {
			//panic(err)
			fmt.Printf("%s : ", targetFile)
			fmt.Println(err)
		} else {
			fmt.Printf("Get(%04d) : %s\n", i+1, targetFile)
			i++
		}
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
//		fmt.Println(resp.Status)
		err2 := fmt.Errorf("Server error : %s", resp.Status)
		return err2
	}

//	fmt.Println( resp.StatusCode)

	// Create the file
	out, err := os.Create(DOWNLOADED + "/"  + filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

