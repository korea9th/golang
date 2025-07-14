//20250620 ocean9th

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
//	"slices"
	"strings"
)





func findfiles(rootpath string, debugFlag int) (int, int) {

	dirCount := 0
	fileCount := 0
//	list := []string{}

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			//			fmt.Printf("[%d]\n", dirCount)
			//fmt.Printf("\nDir : %s\n", path)
			return nil
		}
		
		filedir := filepath.Dir(path)
		file_name := filepath.Base(path)

//		result := strings.SplitAfterN(path, " - ", 3)
		result := strings.Split(info.Name(), " - ")
//		temp_date := result[len(result)-1];
		
		file_date := strings.Split(result[len(result)-1], ".")
		
		file_fullname := filedir + "/" + file_name  // 전체 파일명

		new_name := strings.Replace(file_name, " - " + file_date[0], "", 1) //원본 문자열, 바꿀문자열, 신규문자열, 바꿀횟수
//		tmp_name := strings.TrimRight(new_name, " ")
//		tmp_name = strings.TrimRight(tmp_name, "-")
		new_name = file_date[0] + " - " + new_name
		new_fullname := filedir + "/" + new_name  // 신규 파일명(경로포함)

		os.Rename(file_fullname, new_fullname) // 변경
		

//		fmt.Printf("%v, %v\n", filepath.Dir(path), filepath.Base(path))
//		fmt.Printf("%v\n", result)
//		fmt.Printf("%v\n", new_name)
//		fmt.Printf("%v\n", result[len(result)-1])
		
		
//		ext := filepath.Ext(path)
//		ext = strings.Trim(ext, ".")
		if (debugFlag != 0) {
			fmt.Printf("%v\n", path)
//			fmt.Printf("%v\t%v\n", path, strings.Trim(ext, "."))
			
		}
		
//		list = append(list, ext)
//		list = append(list, strings.ToLower(ext))

		fileCount++
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	
	return dirCount, fileCount//, list
	//	return list
}


func main() {
	args := os.Args
	targetDir := ""
	debugFlag := 0
//	list := []string{}
//	rd_list := []string{}

	if len(args) != 2 && len(args) !=3 {
		return
	}
	targetDir = string(args[1])
	if targetDir == "-h" {
//		printLang()
		return
	}
	
	if len(args) == 3 {
		debugFlag = 1
	}

//	if len(args) == 3 && string(args[1]) == string("-b") {
//		readBatch(args[2])
//		return
//	}

	f, err := os.Open(targetDir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()


//	dirs, files, list := findfiles(targetDir, debugFlag)
	dirs, files := findfiles(targetDir, debugFlag)


	fmt.Printf("\n\n")
	fmt.Printf("Target Dir        : %s\n", targetDir);
	fmt.Printf("Total Dir         : %d\n", dirs);
	fmt.Printf("Total File        : %d\n", files);
	fmt.Printf("\n")
	

	if debugFlag != 0 {
//		fmt.Printf("Total exts        : %s\n", list);
//		fmt.Printf("\n")
//		fmt.Println(list)
	}
//	slices.Sort(list) // sort
//	fmt.Println(list) 
//	rd_list = slices.Compact(list)
//	fmt.Println(slices.Compact(list)) // Efficiently Removing Duplicates from the Slice
//	fmt.Println(rd_list) 
//	fmt.Println(len(rd_list))
//	fmt.Printf("Total exts(rm dup): %s\n", rd_list);
}
