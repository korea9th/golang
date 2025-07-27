//20250717 

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"os/exec"
//	"strings"
//	"bufio"
//	"io/ioutil"
//	"time"
)

/**
*	execCommandStart
**/
func execCommandStart(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	err := cmd.Start()
	return err
}

/**
*	execCommandRun
**/
func execCommandRun(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

/**
*	execCommandOutput
**/
func execCommandOutput(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func findfilessnrun(rootpath string, alg string) (int, int) {

	dirCount := 0
	fileCount := 0
	current_path, _ := os.Getwd()

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			return nil
		}
//		fmt.Printf("%s, %s, %s, %s\n", alg, path, info.Name(), current_path)
		
		exec_err := execCommandRun(current_path + "\\" + "oceanCryptoCavp.exe", alg, path, info.Name(), current_path + "\\result" )
//		exec_err := execCommandRun("go", "version")
		if exec_err != nil {
			fmt.Printf("exec error [%v]\n", exec_err)

		}
		
		fileCount++
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	return dirCount, fileCount
}





func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Printf("exe req_path algorithm\n")
		return
	}
	targetDir := string(args[1])
	alg := string(args[2])

	f, err := os.Open(targetDir)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()


	dirs, files := findfilessnrun(targetDir, alg)


	fmt.Printf("\n\n")
	fmt.Printf("Total Dir : %d, file : %d\n", dirs, files);

}
