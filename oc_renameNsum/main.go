//20250717 

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
//	"strings"
//	"bufio"
	"io/ioutil"
	"time"
	"regexp"
	"strconv"
)


func rename_files(rootpath string) (int, int) {

	dirCount := 0
	fileCount := 0

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			return nil
		}

		re := regexp.MustCompile("[0-9]+")
		file_number := re.FindAllString(info.Name(), -1)  // output: [123 987]
		
		num, err := strconv.Atoi(file_number[len(file_number)-1])
		if err != nil {
			fmt.Printf("Atoi error [%v]\n", err)
		}
		
		
		prefix := fmt.Sprintf("%04d", num)
		
		fmt.Println(prefix, file_number, len(file_number))

		file_fullname := rootpath + "/" + info.Name()  // 전체 파일명
		new_fullname := rootpath + "/" + prefix + " " + info.Name()  // 신규 파일명(경로포함)

		os.Rename(file_fullname, new_fullname) // 변경
		
		fileCount++
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	return dirCount, fileCount
	//	return list
}

func sum_files(rootpath string) (int, int) {

	dirCount := 0
	fileCount := 0
	
	sumfile_name := ""
	a := ""

	now := time.Now()
	formattedDate := now.Format("2006-01-02 150405")
	fmt.Println("포맷된 날짜 및 시간:", formattedDate)

	writeFilePath := formattedDate

	// 추가 쓰기 작업
	file, err := os.OpenFile(writeFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
		
		sumfile_name = info.Name()
		fileCount++
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	a, err = os.Getwd()
	if err != nil {
		fmt.Printf("Getwd error [%v]\n", err)
	}

	sumfile_name = fmt.Sprintf("%s\\%s 0001-%04d.txt", a, sumfile_name, fileCount)
	fmt.Println(sumfile_name)

	file.Close()
	err = os.Rename(writeFilePath, sumfile_name) // 변경
	if err != nil {
		fmt.Printf("Rename error [%v]\n", err)
	}

	
	return dirCount, fileCount
	//	return list
}



func main() {
	args := os.Args
//	targetDir := ""
	dirs := 0
	files := 0
	
	if len(args) < 2 {
		return
	}
	
	targetDir := string(args[1])
	
	if len(args) == 3 {
		if string(args[2]) == "rename" {
			dirs, files = rename_files(targetDir)
			fmt.Printf("rename Total Dir : %d, file : %d\n", dirs, files)
		} else if string(args[2]) == "sum" {
			dirs, files = sum_files(targetDir)
			fmt.Printf("sum    Total Dir : %d, file : %d\n", dirs, files)
		}
	} else {
		dirs, files = rename_files(targetDir)
		fmt.Printf("rename Total Dir : %d, file : %d\n", dirs, files)

		dirs, files = sum_files(targetDir)
		fmt.Printf("sum    Total Dir : %d, file : %d\n", dirs, files)
	}
}
