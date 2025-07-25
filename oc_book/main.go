//20250717 

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
//	"bufio"
	"io/ioutil"
)





func findfiles(rootpath string) (int, int) {

	dirCount := 0
	fileCount := 0

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			//			fmt.Printf("[%d]\n", dirCount)
			//fmt.Printf("\nDir : %s\n", path)
			return nil
		}
/*	
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
		}
		
		b, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
*/
		// 출력파일 생성
		fo, err := os.Create(path + ".txt")
		if err != nil {
			panic(err)
		}
		defer fo.Close()

		
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		
		
		result := strings.Split(string(b), "novel_content")
		result2 := strings.Split(result[1], "responsive")
		
		
//		fmt.Printf("%s", result2[0])

		if err := ioutil.WriteFile(path + ".txt", []byte(result2[0]), 0666); err != nil {
			log.Fatal(err)
		}
		
/*		
		// 출력파일 생성
		fo, err := os.Create(path + ".txt")
		if err != nil {
			panic(err)
		}
		defer fo.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(ScanLines)

		for scanner.Scan() {
//			fmt.Println(scanner.Text())
//			result := strings.SplitAfterN(scanner.Text(), " ", 3)
			fmt.Println(scanner.Text())
			fmt.Println(strings.Contains(scanner.Text(), "novel_content"))
//			result := scanner.Text() + ".txt"
		

				//		result = strings.SplitAfter(scanner.Text(), " ")
//		fmt.Println(result) // 출력: [apple, banana, cherry,date]
		}
		
//		ext := filepath.Ext(path)
//		fmt.Printf("%v\t%v\n", path, ext)
*/		

		fileCount++
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	return dirCount, fileCount
	//	return list
}


func main() {
	args := os.Args
	targetDir := ""

	if len(args) != 2 {
		return
	}
	targetDir = string(args[1])

//	if len(args) == 3 && string(args[1]) == string("-b") {
//		readBatch(args[2])
//		return
//	}

	f, err := os.Open(targetDir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()


	dirs, files := findfiles(targetDir)


	fmt.Printf("\n\n")
	fmt.Printf("Total Dir      : %d\n", dirs);
	fmt.Printf("Total File     : %d\n", files);
	

}
