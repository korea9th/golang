//20250527 ocean9th@naver.com

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		
		ext := filepath.Ext(path)
		fmt.Printf("%v\t%v\n", path, ext)

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
