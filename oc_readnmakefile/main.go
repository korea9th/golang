//20250620 ocean9th

package main

import (
	"fmt"
//	"log"
	"os"
	"path/filepath"
//	"slices"
	"strings"
	"bufio"
)



func printHelp() {
	fmt.Println("help~") 

}




func findfiles(rootpath string, debugFlag int) (int, int, []string) {

	dirCount := 0
	fileCount := 0
	list := []string{}

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			//			fmt.Printf("[%d]\n", dirCount)
			//fmt.Printf("\nDir : %s\n", path)
			return nil
		}
		
		ext := filepath.Ext(path)
		ext = strings.Trim(ext, ".")
		if (debugFlag != 0) {
			fmt.Printf("%v\t%v\n", path, ext)
//			fmt.Printf("%v\t%v\n", path, strings.Trim(ext, "."))
			
		}
		
//		list = append(list, ext)
		list = append(list, strings.ToLower(ext))

		fileCount++
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	
	return dirCount, fileCount, list
	//	return list
}

func countOccurrences(slice []string, target string) int {
	count := 0
	for _, num := range slice {
		if num == target {
			count++
		}
	}
	return count
}

func main() {
	args := os.Args
	targetDir := ""
//	debugFlag := 0
//	list := []string{}
//	list2 := []string{}
//	rd_list := []string{}

	if len(args) != 2 && len(args) !=3 {
		return
	}
	targetDir = string(args[1])
	if targetDir == "-h" {
		printHelp()
		return
	}
	
	if len(args) == 3 {
//		debugFlag = 1
	}

    file, err := os.Open("a.txt")
    if err != nil {
        fmt.Println(err)
    }
	
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
//		fmt.Println(scanner.Text())
//		result := strings.SplitAfterN(scanner.Text(), " ", 3)
		result := scanner.Text() + ".txt"
		
		// 출력파일 생성
		fo, err := os.Create(result)
		if err != nil {
			panic(err)
		}
		defer fo.Close()
 	

		//		result = strings.SplitAfter(scanner.Text(), " ")
//		fmt.Println(result) // 출력: [apple, banana, cherry,date]
    }
    if err := scanner.Err(); err != nil {
        fmt.Println(err)
    }

/*
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		list2 = append(list2, scanner.Text())
		if slices.Contains(list, scanner.Text()) {
//			fmt.Println("-----")
		} else {
			list = append(list, scanner.Text())
		}
 //       fmt.Println(scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        fmt.Println(err)
    }
	
	for i:=0; i<len(list); i++ {
		count := countOccurrences(list2, list[i])
		fmt.Printf("%s\t%d\n", list[i], count)
	}
*/

    file.Close()
//	fmt.Println(list)
	
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


	dirs, files, list := findfiles(targetDir, debugFlag)


	fmt.Printf("\n\n")
	fmt.Printf("Target Dir        : %s\n", targetDir);
	fmt.Printf("Total Dir         : %d\n", dirs);
	fmt.Printf("Total File        : %d\n", files);
	fmt.Printf("\n")
	

	if debugFlag != 0 {
		fmt.Printf("Total exts        : %s\n", list);
		fmt.Printf("\n")
//		fmt.Println(list)
	}
	slices.Sort(list) // sort
//	fmt.Println(list) 
	rd_list = slices.Compact(list)
//	fmt.Println(slices.Compact(list)) // Efficiently Removing Duplicates from the Slice
//	fmt.Println(rd_list) 
//	fmt.Println(len(rd_list))
	fmt.Printf("Total exts(rm dup): %s\n", rd_list);
*/

}
