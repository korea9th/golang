//20250717 

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"bufio"
	"io/ioutil"
	"strings"

)





func findfiles(rootpath string) (int, int) {

	dirCount := 0
	fileCount := 0
	i := 0

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			//			fmt.Printf("[%d]\n", dirCount)
			//fmt.Printf("\nDir : %s\n", path)
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
		}
		


		
	


		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if string(scanner.Text()[0]) == "\"" {
				i = 0
			}
			i++
//			fmt.Println(scanner.Text())
//			result := strings.SplitAfterN(scanner.Text(), " ", 3)
			if i == 6 {
				fmt.Println(scanner.Text())
			} else {
//				fmt.Println("_")

			}
//			if i == 11 {
//				i=0
//			}
//			result := scanner.Text() + ".txt"
		

				//		result = strings.SplitAfter(scanner.Text(), " ")
//		fmt.Println(result) // 출력: [apple, banana, cherry,date]
		}
		
//		ext := filepath.Ext(path)
//		fmt.Printf("%v\t%v\n", path, ext)
	

		fileCount++
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	return dirCount, fileCount
	//	return list
}


func printlinetext(linefile string, codefile string) (int) {


	line := 0



	code_file, err := os.Open(codefile)
	if err != nil {
		fmt.Println(err)
	}
	
	b, err := ioutil.ReadFile(linefile)
	if err != nil {
		log.Fatal(err)
	}
	
	
	result := strings.Split(string(b), "\n")
	fmt.Println(result[line])


	scanner := bufio.NewScanner(code_file)

	for scanner.Scan() {
		code := strings.Split(scanner.Text(), ":")
//		fmt.Println(code[0])
		if strings.Contains(code[0], "\"") {
			line++
		}
//			fmt.Println(strings.Compare(strings.TrimSpace(code[0]), strings.TrimSpace(result[line-1])))

//		fmt.Println(result[line-1])
		if strings.Compare(strings.TrimSpace(code[0]), strings.TrimSpace(result[line-1])) == 0 {
			fmt.Println(strings.TrimSpace(result[line-1]) + "|"+ scanner.Text())
		} else if strings.Compare(strings.Trim(code[0], "\""), strings.TrimSpace(result[line-1])) == 0 {
			fmt.Println(strings.TrimSpace(result[line-1]) + "|")
		}





			//		result = strings.SplitAfter(scanner.Text(), " ")
//		fmt.Println(result) // 출력: [apple, banana, cherry,date]
	}
	
//		ext := filepath.Ext(path)
//		fmt.Printf("%v\t%v\n", path, ext)




	
	return line
	//	return list
}



func main() {




//	if len(args) == 3 && string(args[1]) == string("-b") {
//		readBatch(args[2])
//		return
//	}



	files := printlinetext("line.txt", "code.txt")


	fmt.Printf("\n\n")
	fmt.Printf("Total File     : %d\n", files);
	

}
