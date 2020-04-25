package projectGenerator

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/src/templates"
	"github.com/pepelazz/projectGenerator/src/types"
	"github.com/pepelazz/projectGenerator/src/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	//"github.com/otiai10/copy"
)

type (
	// функция для модификации файлов при копировании из шаблона проекта в конечный проект
	copyFileModifyFunc func(path string, file []byte) []byte
)

var (
	project        types.ProjectType
	tmplMap        map[string]*template.Template
)

func readData(p types.ProjectType)  {
	project = p
	project.DistPath = "../src"
	project.FillDocTemplatesFields()
	project.GenerateGrid()
}

func Start(p types.ProjectType, modifyFunc copyFileModifyFunc)  {
	// читаем данные для проекта
	readData(p)
	// читаем темплейты
	tmplMap = templates.ParseTemplates(project)


	// генерим файлы для проекты
	templates.WriteProjectFiles(p, tmplMap)

	// генерим файлы для документов
	for _, d := range p.Docs {
		for _, dt := range d.Templates {
			err := templates.ExecuteToFile(dt.Tmpl, d, dt.DistPath, dt.DistFilename)
			utils.CheckErr(err, fmt.Sprintf("'%s' ExecuteToFile '%s'", d.Name, dt.DistFilename))
		}
	}

	// копируем файлы проекта (которые не шаблоны)
	err := copyFiles(p,"../../projectGenerator/src/sourceFiles", "../", modifyFunc)
	utils.CheckErr(err, "Copy sourceFiles")
}

// функция для копирования файлов с возможностью модификаации содержимого файлов
func copyFiles(p types.ProjectType, source, dist string, modifyFunc copyFileModifyFunc) (err error)  {
	err = filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				file, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				// для windows заменяем слэши в пути на обратные
				dirPath := strings.TrimSuffix(strings.TrimPrefix(strings.Replace(path, "\\", "/", -1), source), info.Name())
				// создаем директории
				err = os.MkdirAll(dist + dirPath, os.ModePerm)
				if err != nil {
					return err
				}
				// заменяем ссылки в go файлах
				if strings.HasSuffix(info.Name(), ".go") {
					file = []byte(strings.Replace(string(file), "github.com/pepelazz/projectGenerator", p.Config.LocalProjectPath, -1))
				}
				// применяем модификатор для текста файла
				if modifyFunc != nil {
					file = modifyFunc(dirPath + info.Name(), file)
				}
				// записываем файл по новому пути
				err = ioutil.WriteFile(dist + dirPath + info.Name(), file, 0644)
				if err != nil {
					return err
				}
			}
			return nil
		})
	return
}