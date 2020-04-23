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

func ParseTemplateFilename(docName, filename, globalDistPath string) (distPath, distFilename string) {
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

	return
}