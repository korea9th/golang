//20250620 ocean9th

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

/*
CONF	conf, yml, config, properties, gradle, dockerfile
JAVA	java, conf, yml
JAVASCRIPT	js
JSP	jsp
HTML	html
XML	xml, config, properties, gradle
C/C++	c, cc, cpp, c++, cp, cxx, pc, h, hpp
IOS	m, mm, plist
PHP	php, inc
PYTHON	py
SWIFT	swift
ASP	asp
C#	cs
VBSCRIPT	vbs
Flutter(DART)	dart
KOTLIN	kt
TYPESCRIPT	ts, tsx, vue
GO	go
TOBESOFT	xadl, xfdl, xjs, js, xml, config, properties, gradle
ASP.NET	aspx, xml, config, properties, gradle, cs, html
Android	java, conf, yml, xml, config, properties, gradle

var CONF = "CONF, conf, yml, config, properties, gradle, dockerfile"
var JAVA = "JAVA	java, conf, yml"
var JAVASCRIPT = "JAVASCRIPT	js"
var JSP = "JSP	jsp"
var HTML = "HTML	html"
var XML = "XML	xml, config, properties, gradle"
var CCPP = "C/C++	c, cc, cpp, c++, cp, cxx, pc, h, hpp"
var IOS = "IOS	m, mm, plist"
var PHP = "PHP	php, inc"
var PYTHON = "PYTHON	py"
var SWIFT = "SWIFT	swift"
var ASP = "ASP	asp"
var C3 = "C#	cs"
var VBSCRIPT = "VBSCRIPT	vbs"
var Flutter_DART = "Flutter(DART)	dart"
var KOTLIN = "KOTLIN	kt"
var TYPESCRIPT = "TYPESCRIPT	ts, tsx, vue"
var GO = "GO	go"
var TOBESOFT = "TOBESOFT	xadl, xfdl, xjs, js, xml, config, properties, gradle"
var ASPNET = "ASP.NET	aspx, xml, config, properties, gradle, cs, html"
var Android = "Android	java, conf, yml, xml, config, properties, gradle"
*/

var CONF = []string{"CONF", "conf", "yml", "config", "properties", "gradle", "dockerfile"}
var JAVA = []string{"JAVA", "java", "conf", "yml"}
var JAVASCRIPT = []string{"JAVASCRIPT", "js"}
var JSP = []string{"JSP", "jsp"}
var HTML = []string{"HTML", "html"}
var XML = []string{"XML", "xml", "config", "properties", "gradle"}
var CCPP = []string{"C/C++", "c", "cc", "cpp", "c++", "cp", "cxx", "pc", "h", "hpp"}
var IOS = []string{"IOS", "m", "mm", "plist"}
var PHP = []string{"PHP", "php", "inc"}
var PYTHON = []string{"PYTHON", "py"}
var SWIFT = []string{"SWIFT", "swift"}
var ASP = []string{"ASP", "asp"}
var C3 = []string{"C#", "cs"}
var VBSCRIPT = []string{"VBSCRIPT", "vbs"}
var Flutter_DART = []string{"Flutter(DART)", "dart"}
var KOTLIN = []string{"KOTLIN", "kt"}
var TYPESCRIPT = []string{"TYPESCRIPT", "ts", "tsx", "vue"}
var GO = []string{"GO", "go"}
var TOBESOFT = []string{"TOBESOFT", "xadl", "xfdl", "xjs", "js", "xml", "config", "properties", "gradle"}
var ASPNET = []string{"ASP.NET", "aspx", "xml", "config", "properties", "gradle", "cs", "html"}
var Android = []string{"Android", "java", "conf", "yml", "xml", "config", "properties", "gradle"}


func printLang() {
/*
var CONF = "CONF, conf, yml, config, properties, gradle, dockerfile"
var JAVA = "JAVA	java, conf, yml"
var JAVASCRIPT = "JAVASCRIPT	js"
var JSP = "JSP	jsp"
var HTML = "HTML	html"
var XML = "XML	xml, config, properties, gradle"
var CCPP = "C/C++	c, cc, cpp, c++, cp, cxx, pc, h, hpp"
var IOS = "IOS	m, mm, plist"
var PHP = "PHP	php, inc"
var PYTHON = "PYTHON	py"
var SWIFT = "SWIFT	swift"
var ASP = "ASP	asp"
var C3 = "C#	cs"
var VBSCRIPT = "VBSCRIPT	vbs"
var Flutter_DART = "Flutter(DART)	dart"
var KOTLIN = "KOTLIN	kt"
var TYPESCRIPT = "TYPESCRIPT	ts, tsx, vue"
var GO = "GO	go"
var TOBESOFT = "TOBESOFT	xadl, xfdl, xjs, js, xml, config, properties, gradle"
var ASPNET = "ASP.NET	aspx, xml, config, properties, gradle, cs, html"
var Android = "Android	java, conf, yml, xml, config, properties, gradle"
	for i := 1; i < len(rd_list); i++ {
		if slices.Contains(CONF, rd_list[i]) { fmt.Println(CONF[0], rd_list[i]); continue} 
		if slices.Contains(JAVA, rd_list[i]) { fmt.Println(JAVA[0], rd_list[i]); continue} 
		if slices.Contains(JAVASCRIPT, rd_list[i]) { fmt.Println(JAVASCRIPT[0], rd_list[i]); continue} 
		if slices.Contains(JSP, rd_list[i]) { fmt.Println(JSP[0], rd_list[i]); continue} 
		if slices.Contains(HTML, rd_list[i]) { fmt.Println(HTML[0], rd_list[i]); continue} 
		if slices.Contains(XML, rd_list[i]) { fmt.Println(XML[0], rd_list[i]); continue} 
		if slices.Contains(CCPP, rd_list[i]) { fmt.Println(CCPP[0], rd_list[i]); continue} 
		if slices.Contains(IOS, rd_list[i]) { fmt.Println(IOS[0], rd_list[i]); continue} 
		if slices.Contains(PHP, rd_list[i]) { fmt.Println(PHP[0], rd_list[i]); continue} 
		if slices.Contains(PYTHON, rd_list[i]) { fmt.Println(PYTHON[0], rd_list[i]); continue} 
		if slices.Contains(SWIFT, rd_list[i]) { fmt.Println(SWIFT[0], rd_list[i]); continue} 
		if slices.Contains(ASP, rd_list[i]) { fmt.Println(ASP[0], rd_list[i]); continue} 
		if slices.Contains(C3, rd_list[i]) { fmt.Println(C3[0], rd_list[i]); continue} 
		if slices.Contains(VBSCRIPT, rd_list[i]) { fmt.Println(VBSCRIPT[0], rd_list[i]); continue} 
		if slices.Contains(Flutter_DART, rd_list[i]) { fmt.Println(Flutter_DART[0], rd_list[i]); continue} 
		if slices.Contains(KOTLIN, rd_list[i]) { fmt.Println(KOTLIN[0], rd_list[i]); continue} 
		if slices.Contains(TYPESCRIPT, rd_list[i]) { fmt.Println(TYPESCRIPT[0], rd_list[i]); continue} 
		if slices.Contains(GO, rd_list[i]) { fmt.Println(GO[0], rd_list[i]); continue} 
		if slices.Contains(TOBESOFT, rd_list[i]) { fmt.Println(TOBESOFT[0], rd_list[i]); continue} 
		if slices.Contains(ASPNET, rd_list[i]) { fmt.Println(ASPNET[0], rd_list[i]); continue} 
		if slices.Contains(Android, rd_list[i]) { fmt.Println(Android[0], rd_list[i]); continue} 
//		if strings.Contains(CONF, rd_list[i]) {	fmt.Println(CONF[0], rd_list[i]) }
	}
*/	
	fmt.Println(CONF) 
	fmt.Println(JAVA) 
	fmt.Println(JAVASCRIPT) 
	fmt.Println(JSP) 
	fmt.Println(HTML) 
	fmt.Println(XML) 
	fmt.Println(CCPP) 
	fmt.Println(IOS) 
	fmt.Println(PHP) 
	fmt.Println(PYTHON) 
	fmt.Println(SWIFT) 
	fmt.Println(ASP) 
	fmt.Println(C3) 
	fmt.Println(VBSCRIPT) 
	fmt.Println(Flutter_DART) 
	fmt.Println(KOTLIN) 
	fmt.Println(TYPESCRIPT) 
	fmt.Println(GO) 
	fmt.Println(TOBESOFT) 
	fmt.Println(ASPNET) 
	fmt.Println(Android) 
	fmt.Println(CONF) 
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


func main() {
	args := os.Args
	targetDir := ""
	debugFlag := 0
	list := []string{}
	rd_list := []string{}

	if len(args) != 2 && len(args) !=3 {
		return
	}
	targetDir = string(args[1])
	if targetDir == "-h" {
		printLang()
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
}
