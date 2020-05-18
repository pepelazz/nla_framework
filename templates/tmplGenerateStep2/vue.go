package tmplGenerateStep2

import (
	"bytes"
	"fmt"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

// шаблоны для сообщений о задачах
func TasksTmpl(p types.ProjectType)  {
	distPath := fmt.Sprintf("%s/webClient/src/app/components/currentUser/tasks", p.DistPath)
	// находим список файлов компонент в директории
	files, err := ioutil.ReadDir(distPath + "/taskTemplates")
	utils.CheckErr(err, "TasksTmpl")

	funcMap := template.FuncMap{
		"PrintComps": func() string {
			arr := []string{}
			for _, f := range files {
				arr = append(arr, strings.TrimSuffix(f.Name(), ".vue"))
			}
			return strings.Join(arr, ", ")
		},
		"PrintImports": func() (res string) {
			//import defaultTmpl from './taskTemplates/default'
			for _, f := range files {
				res = res + fmt.Sprintf("\n\timport %[1]s from './taskTemplates/%[1]s'	", strings.TrimSuffix(f.Name(), ".vue"))
			}
			return
		},
	}
	path := "../../../pepelazz/projectGenerator/sourceFiles/src/webClient/src/app/components/currentUser/tasks/list.vue"
	t, err := template.New("list.vue").Funcs(funcMap).Delims("[[", "]]").ParseFiles(path)
	utils.CheckErr(err, "OverriteCopiedFiles ParseFiles")

	err = executeToFile(t, "", distPath, "list.vue")
	utils.CheckErr(err, "OverriteCopiedFiles ExecuteToFile")
}

func executeToFile(t *template.Template, d interface{}, path, filename string) error {
	if t == nil {
		log.Fatalf("template is nil for path '%s/%s'\n", path, filename)
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path+"/"+filename, []byte(tpl.String()), 0644)
}