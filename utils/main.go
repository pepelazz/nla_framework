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
	arr := strings.Split(str, " ")
	if len(arr) > 1 {
		arr[0] = strings.Title(arr[0])
		return strings.Join(arr, " ")
	}
	return strings.Title(str)
}

func ParseDocTemplateFilename(docName, filename, globalDistPath string, docIndex int, params map[string]string) (distPath, distFilename string) {
	// разбираем имя шаблона на части
	arr := strings.Split(filename, "_")
	if len(arr) < 2 {
		log.Fatalf("'%s' wrong template name %s. Must be at least two parts separete bay '_'\n", docName, filename)
	}
	// имя итогового файла это последний элемент в массиве
	distFilename = arr[len(arr)-1]
	// шаблоны для webClient
	if arr[0] == "webClient" {
		path := fmt.Sprintf("%s/webClient/src/app/components/%s", globalDistPath, docName)
		// если в параметрах передан webClientPath, то значит перезаписываем стандартный вариант. Обычно это для сулчаев вложенности. Например, cleint/deal
		if params != nil {
			if p, ok := params["doc.Vue.Path"]; ok {
				path = fmt.Sprintf("%s/webClient/src/app/components/%s", globalDistPath, p)
			}
		}
		// собираем путь для генерации файла
		if arr[1] == "comp" {
			path = path + "/comp"
		}
		// если шаблон вида webClient_taskTmpl_ то это шабон для отображения задачи в списке. Копируем его в папку components/currentUser/tasks/taskTemplates
		if arr[1] == "taskTmpl" {
			path = fmt.Sprintf("%s/webClient/src/app/components/currentUser/tasks/taskTemplates", globalDistPath)
		}
		// в случае i18n пишем в отдельную директорию в зависимости от языка
		if strings.HasPrefix(arr[1], "i18n") {
			// проверка что суффикс для языка указан
			if len(arr) < 3 {
				log.Printf("ParseDocTemplateFilename %v length < 3. Missed language suffix for i18n", arr)
			}
			// извлекаем из названия шаблона префикс языка локализации
			langDir := strings.TrimSuffix(arr[2], ".js")
			if langDir == "en" {
				langDir = "en-US"
			}
			path = fmt.Sprintf("%s/webClient/src/i18n/%s", globalDistPath, langDir)
			distFilename = fmt.Sprintf("%s.js", docName)
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

func PathExtractFilename(p string) (path, filename string) {
	arr:= strings.Split(p, "/")
	filename = arr[len(arr)-1]
	path = strings.TrimSuffix(p, "/"+filename)
	return
}

func ByteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}