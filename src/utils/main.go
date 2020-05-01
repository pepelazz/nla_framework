package utils

import (
	"fmt"
	"github.com/serenize/snaker"
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
			// отсчитываем не от 1, потому что есть файлы, которые копируются, а не генерятся из модели. Например, 01_User, 02_UserAuth и пр
			docIndexStr := fmt.Sprintf("%v", docIndex+10)
			path = fmt.Sprintf("%s/sql/model/%s_%s", globalDistPath, docIndexStr, snaker.SnakeToCamel(docName))
		} else {
			path = fmt.Sprintf("%s/sql/template/function/_%s", globalDistPath, snaker.SnakeToCamel(docName))
			// в данном случае имя шаблона извлекаем по другому, потому что в имени функии используется нижние подчеркивания
			distFilename = snaker.CamelToSnake(docName) + strings.TrimPrefix(filename, "sql_function")
		}
		distPath = path
	}

	return
}

func CheckContainsSliceStr(str string, arr ...string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}