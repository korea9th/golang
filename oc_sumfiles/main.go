//20250717 

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
//	"strings"
	"bufio"
	"io/ioutil"
	"time"
)


func read_files(rootpath string) (int, int) {

	dirCount := 0
	fileCount := 0

	now := time.Now()
	formattedDate := now.Format("2006-01-02 150405")
	fmt.Println("포맷된 날짜 및 시간:", formattedDate)

	writeFilePath := formattedDate

	// 추가 쓰기 작업
	file, err := os.OpenFile(writeFilePath + ".txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	err = filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			//			fmt.Printf("[%d]\n", dirCount)
			//fmt.Printf("\nDir : %s\n", path)
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		
		
		sub_string := info.Name() + "\n" + string(b) + "\n\n\n"
		_, err = file.Write([]byte(sub_string))//result_string))
		if err != nil {
			log.Fatalf("Failed to write additional content to file: %s", err)
		}

		fileCount++
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	return dirCount, fileCount
	//	return list
}





func readnprint(rootpath string) {

    // open the file
    file, err := os.Open(rootpath)

    //handle errors while opening
    if err != nil {
        log.Fatalf("Error when opening file: %s", err)
    }

    fileScanner := bufio.NewScanner(file)

    // read line by line
    for fileScanner.Scan() {
        fmt.Println(fileScanner.Text())
    }
    // handle first encountered error while reading
    if err := fileScanner.Err(); err != nil {
        log.Fatalf("Error while reading file: %s", err)
	}

    file.Close()
}



func main() {
	args := os.Args
//	targetDir := ""

	if len(args) != 2 {
		return
	}
	targetDir := string(args[1])

//	if len(args) == 3 && string(args[1]) == string("-b") {
//		readBatch(args[2])
//		return
//	}
/*
	f, err := os.Open(targetDir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
*/

	dirs, files := read_files(targetDir)


	fmt.Printf("\n\n")
	fmt.Printf("Total Dir : %d, file : %d\n", dirs, files);
//	fmt.Printf("Total Dir      : %d\n", dirs);
//	fmt.Printf("Total File     : %d\n", files);
	

}
