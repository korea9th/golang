package main

import (
//	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"bytes"
)

var UPLOADEDFILES = "uploadedFiles"
var WILLDOWNLOAD = "files"
var FILELIST = "list"
var UPLOADFILE = "upload"

func defaultHandler(w http.ResponseWriter, r *http.Request) {
//	var cfg Configs

//	currentTime := time.Now()
	
	fmt.Printf("\n\n\n--------------------------------------------------------\n")
	r.ParseForm()
	//Get 파라미터 및 정보 출력
	fmt.Println("default : ", r.Form)
	fmt.Println("path : ", r.URL.Path)
	fmt.Println("Method : ", r.Method)
	fmt.Println("param : ", r.Form["description"])
	
	//Parameter 전체 출력
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	if len(r.Form["command"]) != 0 && r.Form["command"][0] == "hash" {
		hashFile(w, r)
		return
	}

	localPath := r.URL.Path

	if localPath == "/"+UPLOADFILE {
		_, _, err := r.FormFile("uploadFile")
		if err != nil {
			//		fmt.Println("....................Error Retrieving the File")
			fmt.Printf("NOT file [%s]\n", err)
			//		fmt.Println(err)
			//		return
		} else {
			uploadFile(w, r)
			return
		}

		printUploadHtml(w)
		return
	} else if localPath == "/list" {
		printFileListHtml(w)
		return
	} else if localPath == "/" {
		printIndexHtml(w)
		return
	} else {
		responseFile(localPath, w)
	}
	

	
	
/*
	if r.URL.Path == "/ineedfilelist.a" {
		listUpload()
	}
*/

/*
//	localPath := "www" + r.URL.Path
	localPath := r.URL.Path
	if r.URL.Path == "/" {
		localPath = "www/index.html"
		printIndexHtml(w)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	
	if strings.Contains(localPath, ".txt") {
		localPath = "result/" + localPath
	}

	content, err := ioutil.ReadFile(dir + "/" + localPath)
	if err != nil {
		fmt.Println(err) //jerryprint
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	contentType := getContentType(localPath)
	w.Header().Add("Content-Type", contentType)
*/

//	w.Write(replaceContent(content))//content)
	
	
//	w.Write([]byte("<br>JerryBlack!!!!!!!!!!!!!!!!!!!!!!!!!!!!<br>"))

	//기본 출력
	//	fmt.Fprintf(w, "Golang WebServer Working!")

//	printResults(dir + "/" + RESULT, w)
	
//	fmt.Fprintf(w, "<a href=\""+ hashedfile01 + "\">" + hashedfile01 + "</a><br>")
//	fmt.Fprintf(w, "<a href=\""+ hashedfile02 + "\">" + hashedfile02 + "</a><br>")
	
//	fmt.Fprintf(w, hashedfile01 + "<br>" + hashedfile02)
	
	/*
	   w.Header().Set("Content-Disposition", "attachment; filename=WHATEVER_YOU_WANT")
	   w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	*/

}



func replaceContent(content []byte) ([]byte) {
	var cfg Configs
	currentTime := time.Now()

	cfg = readConfigJson("config.json")
	
	content = bytes.ReplaceAll(content, []byte("ttttt11111"), []byte(productVersion))
	
	content = bytes.ReplaceAll(content, []byte("aaaaa11111"), []byte(currentTime.Format("2006-01-02")))
	content = bytes.ReplaceAll(content, []byte("bbbbb11111"), []byte(cfg.EvaluationFacility))
	content = bytes.ReplaceAll(content, []byte("ccccc11111"), []byte(cfg.ReceiptNumber))
	content = bytes.ReplaceAll(content, []byte("ddddd11111"), []byte(cfg.DeveloperName))
	content = bytes.ReplaceAll(content, []byte("eeeee11111"), []byte(cfg.ProductName))
	content = bytes.ReplaceAll(content, []byte("fffff11111"), []byte(cfg.Description))
	content = bytes.ReplaceAll(content, []byte("ggggg11111"), []byte(cfg.TargetDir))

	if cfg.HashedAlgorithm == "SHA224" {
		content = bytes.ReplaceAll(content, []byte("sssss11111"), []byte("selected"))
		content = bytes.ReplaceAll(content, []byte("sssss22222"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss33333"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss44444"), []byte(""))
	} else if cfg.HashedAlgorithm == "SHA384" {
		content = bytes.ReplaceAll(content, []byte("sssss11111"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss22222"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss33333"), []byte("selected"))
		content = bytes.ReplaceAll(content, []byte("sssss44444"), []byte(""))
	} else if cfg.HashedAlgorithm == "SHA512" {
		content = bytes.ReplaceAll(content, []byte("sssss11111"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss22222"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss33333"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss44444"), []byte("selected"))
	} else {
		content = bytes.ReplaceAll(content, []byte("sssss11111"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss22222"), []byte("selected"))
		content = bytes.ReplaceAll(content, []byte("sssss33333"), []byte(""))
		content = bytes.ReplaceAll(content, []byte("sssss44444"), []byte(""))
	}

	return content
}

func responseFile(filename string, w http.ResponseWriter) {
	filename = strings.Replace(filename, "files/", "empty/", 1)
	filename = strings.Replace(filename, RESULT + "/", "empty/", 1)

	filename = strings.Replace(filename, "txt/", RESULT + "/", 1)
	filename = strings.Replace(filename, "list/", WILLDOWNLOAD + "/", 1)
	
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	
	content, err := ioutil.ReadFile(dir + "/" + filename)
	if err != nil {
		fmt.Println(err) //jerryprint
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	contentType := getContentType(filename)
	w.Header().Add("Content-Type", contentType)
	w.Write(content)

	return
}



func printResults(rootpath string, w http.ResponseWriter) {
	fmt.Fprintf(w, "<br>")

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		} 

		fmt.Fprintf(w, "<a href=\""+ "txt/" + info.Name() + "\">" + info.Name() + "</a><br>")

		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return
}

func printFileListHtml(w http.ResponseWriter) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "<meta charset=\"utf-8\">\n")
	fmt.Fprintf(w, "</head>\n")
	fmt.Fprintf(w, "<body style=\"font-family:consolas;\">\n")
	fmt.Fprintf(w, "<input type=\"button\" value=\"HOME\" style=\"width:100px;height:22px;\" onClick=\"location.href='/'\">\n")


	fmt.Fprintf(w, "<br>")

	err = filepath.Walk(dir + "/" + WILLDOWNLOAD , func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		} 

		fmt.Fprintf(w, "<a href=\""+ "list/" + info.Name() + "\">" + info.Name() + "</a><br>")

		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
	return
}


func printUploadHtml(w http.ResponseWriter) {
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "<meta charset=\"utf-8\">\n")
	fmt.Fprintf(w, "</head>\n")
	fmt.Fprintf(w, "<body style=\"font-family:consolas;\">\n")
	fmt.Fprintf(w, "<fieldset style =\"width:640\">\n")

//	fmt.Fprintf(w, "<pre style=\"line-height:220%%;\">\n")
	fmt.Fprintf(w, "<form action=\"/%s\" method=\"POST\" enctype=\"multipart/form-data\"> \n", UPLOADFILE)
	fmt.Fprintf(w, "Select a file: <input name=\"uploadFile\" type=\"file\"> \n")
	fmt.Fprintf(w, "<input name=\"submit\" type=\"submit\"> \n")

	fmt.Fprintf(w, "</form>\n")
	fmt.Fprintf(w, "<input type=\"button\" value=\"HOME\" style=\"width:100px;height:22px;\" onClick=\"location.href='/'\">\n")

//	fmt.Fprintf(w, "</pre>\n")
	
	fmt.Fprintf(w, "</fieldset>\n")
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
	return
}


func printIndexHtml(w http.ResponseWriter) {
	var cfg Configs
	selected224 := ""
	selected256 := ""
	selected384 := ""
	selected512 := ""

	currentTime := time.Now()

	cfg = readConfigJson("config.json")

	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "<meta charset=\"utf-8\">\n")
	fmt.Fprintf(w, "</head>\n")
	fmt.Fprintf(w, "<body style=\"font-family:consolas;\">\n")
	fmt.Fprintf(w, "<form action=\"/\" method=\"POST\">\n")
	fmt.Fprintf(w, "<fieldset style =\"width:640\">\n")
	fmt.Fprintf(w, "<hr>\n")
	fmt.Fprintf(w, "Code Integrity Tool<br>\n")
	fmt.Fprintf(w, "%s by Jerry<br>\n", productVersion)
	fmt.Fprintf(w, "<hr>\n")
	
	fmt.Fprintf(w, "<legend>Korea Security Evaluation Laboratory</legend>\n")

	fmt.Fprintf(w, "<pre style=\"line-height:220%%;\">\n")
	
	fmt.Fprintf(w, "<input type=\"hidden\" id=\"command\" name=\"command\" value=\"hash\">\n")
	fmt.Fprintf(w, "Date                : <input tyep=\"text\" name=\"date\" disabled size=\"50\" value=\"%s\"/>\n", currentTime.Format("2006-01-02"))
	fmt.Fprintf(w, "Evaluation Facility : <input tyep=\"text\" name=\"evaluationfacility\" size=\"50\" value=\"%s\"/>\n", cfg.EvaluationFacility)
	fmt.Fprintf(w, "Receipt Number      : <input tyep=\"text\" name=\"receiptnumber\" size=\"50 \"value=\"%s\"/>\n", cfg.ReceiptNumber)
	fmt.Fprintf(w, "Developer Name      : <input tyep=\"text\" name=\"developername\" size=\"50 \"value=\"%s\"/>\n", cfg.DeveloperName)
	fmt.Fprintf(w, "Product Name        : <input tyep=\"text\" name=\"productname\" size=\"50\" value=\"%s\"/>\n", cfg.ProductName)
	fmt.Fprintf(w, "Description         : <input tyep=\"text\" name=\"description\" size=\"50\" value=\"%s\"/>\n", cfg.Description)
	fmt.Fprintf(w, "Target dir          : <input tyep=\"text\" name=\"targetdir\" size=\"50\" value=\"%s\"/>\n", cfg.TargetDir)
	fmt.Fprintf(w, "Hashed Algorithm    : <select tyep=\"text\" name=\"hashedalgorithm\"  style=\"width:200px;height:22px;\"/>\n")
	
	if cfg.HashedAlgorithm == "SHA224" {
		selected224 = "selected"
	} else if cfg.HashedAlgorithm == "SHA256" {
		selected256 = "selected"
	} else if cfg.HashedAlgorithm == "SHA384" {
		selected384 = "selected"
	} else if cfg.HashedAlgorithm == "SHA512" {
		selected512 = "selected"
	}

	fmt.Fprintf(w, "	<option value=\"SHA224\" %s>SHA224</option>\n", selected224)
	fmt.Fprintf(w, "	<option value=\"SHA256\" %s>SHA256</option>\n", selected256)
	fmt.Fprintf(w, "	<option value=\"SHA384\" %s>SHA384</option>\n", selected384)
	fmt.Fprintf(w, "	<option value=\"SHA512\" %s>SHA512</option>\n", selected512)
	fmt.Fprintf(w, "</select>\n")
	fmt.Fprintf(w, "</pre>\n")
	
	fmt.Fprintf(w, "<input type=\"submit\" value=\"submit\" style=\"width:100px;height:22px;\">&nbsp\n")
//	fmt.Fprintf(w, "<input type = \"reset\" value = \"reset\" style=\"width:50px;height:22px;\"/>\n")
	fmt.Fprintf(w, "</form>\n")
	fmt.Fprintf(w, "<input type=\"button\" value=\"upload file\" style=\"width:100px;height:22px;\" onClick=\"location.href='/%s'\">\n", UPLOADFILE)
	fmt.Fprintf(w, "<input type=\"button\" value=\"file list\" style=\"width:100px;height:22px;\" onClick=\"location.href='/%s'\">\n", FILELIST)

	
	fmt.Fprintf(w, "</fieldset>\n")
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")

	return
}


func getContentType(localPath string) string {
	var contentType string
	ext := filepath.Ext(localPath)

	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	default:
		contentType = "text/plain; charset=utf-8"
	}

	return contentType
}

func hashFile(w http.ResponseWriter, r *http.Request) {
	var cfg Configs

	cfg.EvaluationFacility = r.Form["evaluationfacility"][0]
	cfg.ReceiptNumber = r.Form["receiptnumber"][0]
	cfg.DeveloperName = r.Form["developername"][0]
	cfg.ProductName = r.Form["productname"][0]
	cfg.Description = r.Form["description"][0]
	cfg.TargetDir = r.Form["targetdir"][0]
	cfg.HashedAlgorithm = r.Form["hashedalgorithm"][0]
	
//	go makeHash(cfg)
	writeConfigJson("config.json", cfg)
//		hashedfile01, hashedfile02 = makeHash(cfg)
	makeHash(cfg)
	
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	
/*
	localPath := "www/index.html"
	content, err := ioutil.ReadFile(dir + "/" + localPath)
	if err != nil {
		fmt.Println(err) //jerryprint
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	contentType := getContentType(localPath)
	w.Header().Add("Content-Type", contentType)
	

	w.Write(replaceContent(content))//content)
*/
	printIndexHtml(w)
	
	printResults(dir + "/" + RESULT, w)

//		return
}


func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	extention := filepath.Ext(handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our UPLOADEDFILES directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(UPLOADEDFILES, "upload-*"+extention) //png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	//	fmt.Fprintf(w, "Successfully Uploaded File\n")

	// 결과 출력
	printUploadHtml(w)
/*
	content, err := ioutil.ReadFile("www/index.html")
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		fmt.Println(err)
		return
	}

	contentType := getContentType("www/index.html")
	w.Header().Add("Content-Type", contentType)
	w.Write(content)
*/	
	fmt.Fprintf(w, "\n[%s] uploaded.\n", handler.Filename)
	//	listUpload()
}

func listUpload() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("listUpload error [%v]\n", err)
		return
	} else {
		//		fmt.Printf(path)
	}

	// 출력파일 생성
	fo, err := os.Create(path + "/wwwroot/ineedfilelist.a")
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	err = filepath.Walk(path+"/wwwroot/files", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fmt.Fprintf(fo, "%s\n", info.Name())
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return

}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func ocHttpd() {
	var cfg ServerConfigs
	
	cfg = readServerConfigJson("serverconfig.json")

	port := cfg.HttpPort
	WILLDOWNLOAD = cfg.HttpDir
/*
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
*/	
	fmt.Println("My Home : " + MYPATH)

	CreateDirIfNotExist(MYPATH + "/" + UPLOADEDFILES)
	CreateDirIfNotExist(MYPATH + "/" + WILLDOWNLOAD)

	//기본 Url 핸들러 메소드 지정
	http.HandleFunc("/", defaultHandler)
	//서버 시작

	log.Printf("[HTTP Server]ListenAndServe Started! -> Port(%s)", port)
//	fmt.Println("[HTTP Server]ListenAndServe Started! -> Port(" + port + ")")
	err := http.ListenAndServe(":"+port, nil) //":9090", nil)
	//예외 처리
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("ListenAndServe Started! -> Port(" + port + ")")
	}
}

func ocHttpFileServer() {
	var cfg ServerConfigs
	
	cfg = readServerConfigJson("serverconfig.json")

	port := cfg.HttpPort
	dir := cfg.HttpDir

	log.Printf("[File Server]ListenAndServe Started! -> Port(%s)", port)

//	fmt.Println("[File Server]ListenAndServe Started! -> Port(" + port + ")")
	http.Handle("/", http.FileServer(http.Dir("./"+dir)))
	http.ListenAndServe(":"+port, nil)
}




/*
func ocHttpd() {
	fmt.Printf("usage : GoHttpd -port=8080\n")

	var (
		port = flag.String("port", "8080", "port")
	)
	flag.Parse()
	if *port == "" {
		log.Fatalf("Please set a PORT to serve with -port")
	}

	port := "8080"

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	fmt.Println("My Home : " + dir)

	CreateDirIfNotExist(dir + "/" + UPLOADEDFILES)
	CreateDirIfNotExist(dir + "/" + WILLDOWNLOAD)

	//기본 Url 핸들러 메소드 지정
	http.HandleFunc("/", defaultHandler)
	//서버 시작
	fmt.Println("ListenAndServe Started! -> Port(" + *port + ")")
	err = http.ListenAndServe(":"+*port, nil) //":9090", nil)
	//예외 처리
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("ListenAndServe Started! -> Port(" + *port + ")")
	}
}
*/
