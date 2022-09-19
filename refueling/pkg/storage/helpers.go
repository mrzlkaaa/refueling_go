package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	basePathTo string = "pkg/storage/data" //* will be joined with filepath.Dir(os.Getwd())
)

//! is not for use in prod
func FormatterCoreConfig(coreConfig *[][]string) *[]byte {
	var str string
	sliceLen := len(*coreConfig)
	for i, vv := range *coreConfig {
		for _, elem := range vv {
			str += elem + ","
		}
		if sliceLen == i+1 {
			str = str[:len(str)-1]
		}
	}
	formattedConfig := []byte(str)

	return &formattedConfig
}

//! is not for use in prod
func BackFormatterCoreConfig(coreConfig *[]byte) *[][]string {
	var str string
	var firstArr []string
	var arr2D [][]string

	for _, v := range *coreConfig { //*conv to str
		str += string(v)
	}
	array := strings.Split(str, ",")
	for i, v := range array {
		i += 1
		firstArr = append(firstArr, v)
		if i%4 == 0 {
			arr2D = append(arr2D, firstArr)
			firstArr = []string{}
		}
	}
	return &arr2D
}

//! is not for use in prod
func FormatterPDC(pdc *[]string) *[]byte {
	joined := strings.Join(*pdc, "")
	formattedPDC := []byte(joined)
	return &formattedPDC
}

//! is not for use in prod
func BackFormatterPDC(pdc *[]byte) *[]string {
	var str string
	var arr []string
	// for _, v := range *pdc { //*conv to str
	// 	str += string(v)
	// }
	str = string(*pdc)
	arr = strings.Split(str, "\n")
	for i := 0; i < len(arr); i++ {
		arr[i] += "\n"
	}
	// arr = append(arr, str)
	return &arr
}

func CheckPathToStoredFiles(refuelName int) ([]fs.FileInfo, string, error) {
	var rootFolder string = "data"
	path := filepath.Join(rootFolder, strconv.Itoa(refuelName))
	//* check if folder exists
	if _, err := os.Stat(path); err != nil {
		return []fs.FileInfo{}, "", err
	}
	files, err := ioutil.ReadDir(path) //* check files in folder
	if err != nil || len(files) == 0 {
		return files, "", errors.New("while checking path to stored files an error has occured")
		// return []string{}
	}
	return files, path, nil
	// return nil
}

func GetFileName(refuel int, name, extension string) string {
	return fmt.Sprintf("%v_%v.%v", strconv.Itoa(refuel), name, extension)
}

func OpenReadStoredConfig(filePath string, fileData chan<- [][]string) {
	wg.Add(1)
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("openning file...", file.Name())
	var arr2D [][]string
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		arr2D = append(arr2D, strings.Split(fileScanner.Text(), "    ")) //* split by 4 spaces
		// fmt.Println(arr)
	}
	fmt.Println("this is 2d arr", arr2D)
	fileData <- arr2D
	fmt.Printf("\n")
	//todo add text to 2d array --> return
}

func OpenReadStoredPDC(filePath string, lineData chan<- []string) {
	wg.Add(1)
	defer wg.Done()
	file, _ := os.Open(filePath)
	fmt.Println("openning file...", file.Name())
	// if err != nil {
	// 	//* do error handler
	// }
	var arr []string
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		arr = append(arr, fileScanner.Text()+"\n")
	}
	lineData <- arr
	//todo add text to 2d array --> return
}
