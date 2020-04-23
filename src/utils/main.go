package utils

import (
	"fmt"
	"log"
	"strings"
)

func CheckErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func UpperCaseFirst(str string) string {
	return strings.Title(strings.ToLower(str))
}

func ParseDocTemplateFilename(docName, filename, globalDistPath string, docIndex int) (distPath, distFilename string) {
	// разбираем имя шаблона на части
	arr := strings.Split(filename, "_")
	if len(arr) < 2 {
		log.Fatalf("'%s' wrong template name %s. Must be at least two parts separete bay '_'\n", docName, filename)
	}
	// имя итогового файла это последний элемент в массиве
	distFilename = arr[len(arr)-1]
	// шаблоны для webClient
	if arr[0] == "webClient" {
		// собираем путь для генерации файла
		path := fmt.Sprintf("%s/webClient/src/app/components/%s", globalDistPath, docName)
		if arr[1] == "comp" {
			path = path + "/comp"
		}
		distPath = path
	}
	// шаблоны для sql
	if arr[0] == "sql" {
		// собираем путь для генерации файла
		var path string
		if arr[1] == "main.toml" {
			// формируем числовой префикс для названия файла (для устойчивой сортировки)
			var docIndexStr string
			if docIndex < 9 {
				docIndexStr = fmt.Sprintf("0%v", docIndex+1)
			} else {
				docIndexStr = fmt.Sprintf("%v", docIndex+1)
			}
			path = fmt.Sprintf("%s/sql/model/%s_%s", globalDistPath, docIndexStr, UpperCaseFirst(docName))
		} else {
			path = fmt.Sprintf("%s/sql/template/function/_%s", globalDistPath, UpperCaseFirst(docName))
		}
		distPath = path
	}

	return
}