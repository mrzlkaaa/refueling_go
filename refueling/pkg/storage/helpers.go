package storage

import "strings"

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

func FormatterPDC(pdc *[]string) *[]byte {
	joined := strings.Join(*pdc, "")
	formattedPDC := []byte(joined)
	return &formattedPDC
}

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
