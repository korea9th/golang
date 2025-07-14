//20201106 Ocean
package main

//import ("fmt"; "io/ioutil"; "log")
import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"bufio"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
//	"syscall"
//	"os/exec"
//	"unicode/utf8"
//	"golang.org/x/text/encoding/charmap"
)
var productVersion = "ocGoHash V3.3(Golang)"
var hashs = [4]string{"SHA224", "SHA256", "SHA384", "SHA512"}

var titleLine    = "--------------------------------------------------"
var titleVersion = "    Code Integrity Tool"
var titleAuth    = "    " + productVersion + " by Jerry"

	
func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Configs struct {
    EvaluationFacility    string `json:"evaluationfacility"`
    ReceiptNumber         string `json:"receiptnumber"`
    DeveloperName         string `json:"developername"`
    ProductName           string `json:"productname"`
    Description           string `json:"description"`
    TargetDir             string `json:"targetdir"`
    HashedAlgorithm       string `json:"hashedalgorithm"`
}


func readConfigJson(jsonFileName string) (Configs){
    // Open our jsonFile
    jsonFile, err := os.Open(jsonFileName)
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }

    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    // read our opened xmlFile as a byte array.
    byteValue, _ := ioutil.ReadAll(jsonFile)

    // we initialize our Users array
//    var users Users
    var cfg Configs

    // we unmarshal our byteArray which contains our
    // jsonFile's content into 'users' which we defined above
    json.Unmarshal(byteValue, &cfg)//&users)
	
	return cfg
}

func writeConfigJson(jsonFileName string, config Configs) {
	b, _ := json.MarshalIndent(config, "", "	")//&users)

	f1, err := os.Create(jsonFileName)
	check(err)
	defer f1.Close()
	fmt.Fprintf(f1, string(b))
}

func printMainMenu() {
//	fmt.Printf("\t1. Change Vendor/Product Name\n");
	fmt.Printf("\t1. Change [Evaluation Facility]\n");
	fmt.Printf("\t2. Change [Receipt Number]\n");
	fmt.Printf("\t3. Change [Developer Name]\n");
	fmt.Printf("\t4. Change [Product Name]\n");
	fmt.Printf("\t5. Change [Target Dir]\n");
//	fmt.Printf("\t3. Change Result File Name\n");
	fmt.Printf("\t6. Change Hashed Algorithm\n");
	fmt.Printf("\t9. Make Hash\n");
	fmt.Printf("\t0. Exit\n");
	fmt.Printf("\n");
	fmt.Printf("\n\tSelect Number : ");
}

func printHelp(myName string) {
	fmt.Printf("Usage: %s [OPTION] [FILE]\n", myName);
	fmt.Printf("\t-b configfileslist\tmake hashs by all jsonfiles in configfileslist\n");
	fmt.Printf("\t-h     display this help and exit\n");
	
	
}

func main() {

	args := os.Args

	if len(args) == 3 && string(args[1]) == string("-b") {
		readBatch(args[2])
		return
	} else if len(args) == 2 && string(args[1]) == string("-h") {
		printHelp(args[0])
		return
	}
	
	currentTime := time.Now()
	
	hashID := 2
	reader := bufio.NewReader(os.Stdin)
	
	var cfg Configs
	cfg = readConfigJson("config.json")
	
	for i := range hashs {
		if hashs[i] == cfg.HashedAlgorithm {
			hashID = i + 1
			break;
		}
	}

	evaluationFacility := cfg.EvaluationFacility
	receiptNumber := cfg.ReceiptNumber
	developerName := cfg.DeveloperName
	productName := cfg.ProductName
	description := cfg.Description
	targetDir := cfg.TargetDir
	hashedAlgorithm := cfg.HashedAlgorithm
//	resultfileName := "result.txt"
	

	for {
		fmt.Printf ("\n\n")
		fmt.Println(titleLine)
		fmt.Println(titleVersion)
		fmt.Println(titleAuth)
		fmt.Println(titleLine)
		fmt.Printf ("Date                : %s\n", currentTime.Format("2006-01-02"))
		fmt.Printf ("Evaluation Facility : %s\n", evaluationFacility)
		fmt.Printf ("Receipt Number      : %s\n", receiptNumber)
		fmt.Printf ("Developer Name      : %s\n", developerName);// developerName);
		fmt.Printf ("Product Name        : %s\n", productName);// productName);
		fmt.Printf ("Description         : %s\n", description);// description);
		fmt.Printf ("Target dir          : %s\n", targetDir);//targetDir)
		fmt.Printf ("Hashed Algorithm    : %s\n", hashs[hashID-1])
		
//		fmt.Printf ("Result File Name    : %s\n", resultfileName)
		fmt.Printf ("\n\n")
		
		printMainMenu()

//		fmt.Print("Enter text: ")
		input, _ := reader.ReadString('\n')
//		fmt.Println(input)
		
		
		selectedNum, e := strconv.Atoi(strings.Trim(strings.Trim(input,"\n"), "\r"))
		if e != nil {
			fmt.Print("Only Number(0~9)\t")
			fmt.Println(e)
//			log.Fatal(e)
			continue
//			return
		}
		
		if (selectedNum == 0) {
			return
		} else if (selectedNum == 1) {
			fmt.Printf ("Evaluation Facility       : ")
			input, _ = reader.ReadString('\n')
			evaluationFacility = strings.Trim(strings.Trim(input,"\n"), "\r")
			cfg.EvaluationFacility = evaluationFacility
		} else if (selectedNum == 2) {
			fmt.Printf ("Receipt Number       : ")
			input, _ = reader.ReadString('\n')
			receiptNumber = strings.Trim(strings.Trim(input,"\n"), "\r")
			cfg.ReceiptNumber = receiptNumber
		} else if (selectedNum == 3) {
			fmt.Printf ("Developer Name       : ")
			input, _ = reader.ReadString('\n')
			developerName = strings.Trim(strings.Trim(input,"\n"), "\r")
			cfg.DeveloperName = developerName
		} else if (selectedNum == 4) {
			fmt.Printf ("Product Name       : ")
			input, _ = reader.ReadString('\n')
			productName = strings.Trim(strings.Trim(input,"\n"), "\r")
			cfg.ProductName = productName
		} else if (selectedNum == 5) {
			fmt.Printf ("Target Dir           : ")
			input, _ = reader.ReadString('\n')
			targetDir = strings.Trim(strings.Trim(input,"\n"), "\r")
			cfg.TargetDir = targetDir
		} else if (selectedNum == 6) {
			fmt.Printf ("Hashed Algoritm(1:SHA224, 2:SHA256, 3:SHA384, 4:SHA512)      : ")
			input, _ = reader.ReadString('\n')
			hashedAlgorithm = strings.Trim(strings.Trim(input,"\n"), "\r")
			cfg.HashedAlgorithm = hashedAlgorithm

			selectedHashid, e := strconv.Atoi(strings.Trim(strings.Trim(input,"\n"), "\r"))
			if e != nil {
				fmt.Print("Only Number(0~9)\t")
				fmt.Println(e)
	//			log.Fatal(e)
				continue
				//return
			}
			if (selectedHashid == 1) {
				hashID = 1
			} else if (selectedHashid == 2) {
				hashID = 2
			} else if (selectedHashid == 3) {
				hashID = 3
			} else if (selectedHashid == 4) {
				hashID = 4
			} else {
				hashID = 2
			}
		} else if (selectedNum == 99) {
/*		
			var cmdName = "tree"
			var cmd = exec.Command(cmdName, "C:\\00.JerryBlack\\01.Music\\BTS\\The Most Beautiful Moment In Life Young Forever", "/A", "/F")
			
			grepOut, _ := cmd.StdoutPipe()
			cmd.Start()
			grepBytes, _ := ioutil.ReadAll(grepOut)
			cmd.Wait()
			fmt.Printf("%x\n", grepBytes)
*/			
			break
		} else if (selectedNum == 9) {
			
			fmt.Printf("\n")
			writeConfigJson("config.json", cfg)
//			makeHashBatch("config.json")
			makeHash(cfg)
//			break
		}
	}
}

func readBatch(configFileName string) {
	f, err := os.Open(configFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	reader := bufio.NewReader(f)
	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix || err != nil {
			break
		}
		
		makeHashBatch(string(line))
	}

}


func makeHashBatch(jsonFileName string) {

	var cfg Configs
	cfg = readConfigJson(jsonFileName)
	
	makeHash(cfg)
	
}

func makeHash(cfg Configs) {

	currentTime := time.Now()
	
	hashID := 2

	for i := range hashs {
		if hashs[i] == cfg.HashedAlgorithm {
			hashID = i + 1
			break;
		}
	}

	evaluationFacility := cfg.EvaluationFacility
	receiptNumber := cfg.ReceiptNumber
	developerName := cfg.DeveloperName
	productName := cfg.ProductName
	description := cfg.Description
	targetDir := cfg.TargetDir
	resultfileName := "result.txt"

//	resultfileName = fmt.Sprintf("%s_%s_%s_%s.txt", currentTime.Format("2006-01-02"), receiptNumber, productName, description)
	resultfileName = fmt.Sprintf("%s_%s_%s_%s.txt", currentTime.Format("20060102-150405"), receiptNumber, productName, description)

	fo, er := os.Create(resultfileName)
	check(er)
	defer fo.Close()


	_, allFiles := getCounts(targetDir)

	t := time.Now()
	dirs, files := makeHashDirCSV(targetDir, hashID, fo, cfg, allFiles) //"C:\\go\\001.project\\include")
	t2 := time.Now()

	fmt.Fprintf(fo, "Total Scan Dir,%d\n", dirs)
	fmt.Fprintf(fo, "Total Scan File,%d\n", files)

	fmt.Printf("\n\n")
	fmt.Printf("Total Scan Dir      : %d\n", dirs);
	fmt.Printf("Total Scan File     : %d\n", files);
	fmt.Printf("Total Process times : %s\n", t2.Sub(t))
	
	tempString := ""
	newName := strings.Split(resultfileName, ".")
    for i := 0; i < len(newName); i++ {
        tempString = tempString + newName[i] + "."
    }

	newResultName := fmt.Sprintf("%shashed.txt", tempString)//newName[0])
	fonew, er := os.Create(newResultName)
	check(er)
	defer fonew.Close()
	
	fmt.Fprintf(fonew, "%s\n", titleLine)
	fmt.Fprintf(fonew, "%s\n", titleVersion)
	fmt.Fprintf(fonew, "%s\n", titleAuth)
	fmt.Fprintf(fonew, "%s\n", titleLine)
	fmt.Fprintf(fonew, "Date                : %s\n", currentTime.Format("2006-01-02"))
	fmt.Fprintf(fonew, "Evaluation Facility : %s\n", evaluationFacility)
	fmt.Fprintf(fonew, "Receipt Number      : %s\n", receiptNumber)
	fmt.Fprintf(fonew, "Developer Name      : %s\n", developerName);// productName);
	fmt.Fprintf(fonew, "Product Name        : %s\n", productName);// productName);
	fmt.Fprintf(fonew, "Description         : %s\n", description);// description);
	fmt.Fprintf(fonew, "Hashed Algorithm    : %s\n", hashs[hashID-1])
	fmt.Fprintf(fonew, "Target File         : %s\n", resultfileName)
	fmt.Fprintf(fonew, "Result File Name    : %s\n", newResultName)
	fmt.Fprintf(fonew, "\n")
	
	fmt.Fprintf(fonew, "Total Scan Dir      : %d\n", dirs);
	fmt.Fprintf(fonew, "Total Scan File     : %d\n", files);
	fmt.Fprintf(fonew, "Total Process times(%s) : %s\n", targetDir, t2.Sub(t))
	fmt.Fprintf(fonew, "\n\n")
	
//	dstFileName := makeHashFile(resultfileName, hashID, fonew, targetDir)
	makeHashFile(resultfileName, hashID, fonew, targetDir)
	
//	currentTime = time.Now()
//	dstFileName = currentTime.Format("20060102_150405")
//	resultName_test := fmt.Sprintf("%s.%s", resultfileName, dstFileName)
	
	
	fo.Close()
//	er = os.Rename(resultfileName, resultName_test)
//	if er != nil {
//		log.Fatal(er)
//	}

}


func makeHashFile(filename string, hash int, fo *os.File, rootpath string) (string) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	h := sha256.New224()

	if hash == 1 {
		h = sha256.New224()
	} else if hash == 2 {
		h = sha256.New()
	} else if hash == 3 {
		h = sha512.New384()
	} else {
		h = sha512.New()
	}

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	
	returnString := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Fprintf(fo, "%s\t%s\n", filename, returnString)
	
	fmt.Fprintf(fo, "\n\n")
	

//218
	return returnString
}


func makeHashDir(rootpath string, hash int, fo *os.File) (int, int) {

	//	list := make([]string, 0, 10)
	dirCount := 0
	fileCount := 0

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		h := sha256.New224()
		if info.IsDir() {
			dirCount++
			//			fmt.Printf("[%d]\n", dirCount)
			//fmt.Printf("\nDir : %s\n", path)

			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if hash == 1 {
			h = sha256.New224()
		} else if hash == 2 {
			h = sha256.New()
		} else if hash == 3 {
			h = sha512.New384()
		} else {
			h = sha512.New()
		}

		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

//		fmt.Fprintf(fo, "[%s%s%s] <%x>\n", path, PATH_SEPARATOR, info.Name(), h.Sum(nil))
		fmt.Fprintf(fo, "%s\t%x\n", path, h.Sum(nil))
		fileCount++

		if (fileCount%10 == 0) {
			fmt.Printf(".")
		}
		if (fileCount%500 == 0) {
			fmt.Printf("\n")
		}

		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	
	return dirCount, fileCount
	//	return list
}
/*
	progressBar := [20]string{
	"[=                             ]",
	"[==                            ]",
	"[===                           ]",
	"[====                          ]",
	"[=====                         ]",
	"[======                        ]",
	"[=======                       ]",
	"[========                      ]",
	"[=========                     ]",
	"[==========                    ]",
	"[===========                   ]",
	"[============                  ]",
	"[=============                 ]",
	"[==============                ]",
	"[===============               ]",
	"[================              ]",
	"[=================             ]",
	"[==================            ]",
	"[===================           ]",
	"[====================          ]",
	"[=====================         ]",
	"[======================        ]",
	"[=======================       ]",
	"[========================      ]",
	"[=========================     ]",
	"[==========================    ]",
	"[===========================   ]",
	"[============================  ]",
	"[============================= ]",
	"[==============================]"
	}

*/

func makeHashDirCSV(rootpath string, hash int, fo *os.File, cfg Configs, allFiles int) (int, int) {

	progressBar := [30]string{
	"[#                             ]",
	"[##                            ]",
	"[###                           ]",
	"[####                          ]",
	"[#####                         ]",
	"[######                        ]",
	"[#######                       ]",
	"[########                      ]",
	"[#########                     ]",
	"[##########                    ]",
	"[###########                   ]",
	"[############                  ]",
	"[#############                 ]",
	"[##############                ]",
	"[###############               ]",
	"[################              ]",
	"[#################             ]",
	"[##################            ]",
	"[###################           ]",
	"[####################          ]",
	"[#####################         ]",
	"[######################        ]",
	"[#######################       ]",
	"[########################      ]",
	"[#########################     ]",
	"[##########################    ]",
	"[###########################   ]",
	"[############################  ]",
	"[############################# ]",
	"[##############################]" }

	dirCount := 0
	fileCount := 0
	

	fmt.Printf("\n\nTarget dir,%s\n", cfg.TargetDir);//targetDir)

	fmt.Fprintf(fo, "Title,Desc\n")
	fmt.Fprintf(fo, "Evaluation Facility,%s\n", cfg.EvaluationFacility)
	fmt.Fprintf(fo, "Receipt Number,%s\n", cfg.ReceiptNumber)
	fmt.Fprintf(fo, "Developer Name,%s\n", cfg.DeveloperName);// productName);
	fmt.Fprintf(fo, "Product Name,%s\n", cfg.ProductName);// productName);
	fmt.Fprintf(fo, "Description,%s\n", cfg.Description);// productName);
	fmt.Fprintf(fo, "Target dir,%s\n", cfg.TargetDir);//targetDir)
	fmt.Fprintf(fo, "Hashed Algorithm,%s\n", cfg.HashedAlgorithm)
	fmt.Fprintf(fo, "Program Version,%s\n", productVersion)
	fmt.Fprintf(fo, "\n")


	fmt.Fprintf(fo, "Path,Hashed Data\n")

	
	fmt.Print("\033[s")
	
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		h := sha256.New224()
		if info.IsDir() {
			dirCount++
			return nil
		}
		
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if hash == 1 {
			h = sha256.New224()
		} else if hash == 2 {
			h = sha256.New()
		} else if hash == 3 {
			h = sha512.New384()
		} else {
			h = sha512.New()
		}

		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(fo, "%s,%x\n", path, h.Sum(nil))
		fileCount++
		progressBarIndex := (fileCount*30)/allFiles
		if progressBarIndex > 29 {
			progressBarIndex = 29
		}
		fmt.Print("\033[u\033[K")
		fmt.Printf("%s(%d/%d)", progressBar[progressBarIndex], fileCount, allFiles)
		
		
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	
	
	return dirCount, fileCount
	//	return list
}



func getCounts(rootpath string) (int, int) {

	//	list := make([]string, 0, 10)

	fileCount := 0
	dirCount := 0

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirCount++
			return nil
		} else {
			fileCount++
		}

		/*
			if filepath.Ext(path) == ".sh" {
				list = append(list, path)
			}
		*/
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return dirCount, fileCount
	//	return list
}





/*
func makeHashFile(filename string, hash int, fo *os.File) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	h := sha256.New224()

	if hash == 1 {
		h = sha256.New224()
	} else if hash == 2 {
		h = sha256.New()
	} else if hash == 3 {
		h = sha512.New384()
	} else {
		h = sha512.New()
	}

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(fo, "%s(withDate)\t%x\n", filename, h.Sum(nil))

//218
	return
}

func makeHashFile2(filename string, hash int, fo *os.File) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	h := sha256.New224()

	if hash == 1 {
		h = sha256.New224()
	} else if hash == 2 {
		h = sha256.New()
	} else if hash == 3 {
		h = sha512.New384()
	} else {
		h = sha512.New()
	}

	f.Seek(220, 0) //without Date
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(fo, "%s(withoutDate)\t%x\n", filename, h.Sum(nil))

//218
	return
}
*/
