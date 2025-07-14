//20250527 ocean9th@naver.com
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

	

// 파일 확장자를 변경하는 함수
func changeFileExtension(fileName, newExtension string) string {
	ext := filepath.Ext(fileName)
	return fileName[:len(fileName)-len(ext)] + newExtension
}

func extchange(rootpath string, oldext string, newext string) (int, int) {

	dirCount := 0
	fileCount := 0

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			//			fmt.Printf("[%d]\n", dirCount)
			//fmt.Printf("\nDir : %s\n", path)
			return nil
		}
		
		ext := filepath.Ext(path)

		if string(ext) == oldext {
			// 새 파일 이름 생성
			newpath := changeFileExtension(path, newext)
			oldpath := path

			// 파일 이름 변경
			err := os.Rename(path, newpath)
			if err != nil {
				fmt.Printf("파일 이름 변경 실패: %v\n", err)
				return nil
			}			

			fmt.Printf("%v\t%v\n", oldpath, newpath)
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


func main() {
	args := os.Args
	targetDir := ""
	oldext := ""
	newext := ""

	if len(args) != 4 {
		fmt.Printf("Usage: %s \t [targetDir]	[oldExt]	[newExt]\n", args[0]);
		fmt.Printf("");
		
		return
	}
	targetDir = string(args[1])
	oldext = string(args[2])
	newext = string(args[3])

//	if len(args) == 3 && string(args[1]) == string("-b") {
//		readBatch(args[2])
//		return
//	}

	f, err := os.Open(targetDir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()


	dirs, files := extchange(targetDir, oldext, newext)


	fmt.Printf("\n\n")
	fmt.Printf("Total Dir      : %d\n", dirs);
	fmt.Printf("Total File     : %d\n", files);
}
