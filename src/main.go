package projectGenerator

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/src/exampleProject"
	"github.com/pepelazz/projectGenerator/src/templates"
	"github.com/pepelazz/projectGenerator/src/types"
	"github.com/pepelazz/projectGenerator/src/utils"
	"text/template"
)

var (
	project        types.ProjectType
	tmplMap        map[string]*template.Template
)

func main() {
	// читаем данные для проекта
	readData(exampleProject.GetProject())
	// читаем темплейты
	tmplMap = templates.ParseTemplates(project)

	//if tmpls, ok := tmplMap["webClient"]; ok {
	//	t := tmpls.Lookup("docItem.vue")
	//	if t == nil {
	//		log.Fatalf("no found template: docItem.vue")
	//	}
	//	err := templates.ExecuteToFile(t, project.GetDocByName("city"), fmt.Sprintf("%s/webClient/src/app/components/city", globalDistPath), "item.vue")
	//	utils.CheckErr(err, "ExecuteToFile")
	//	err = templates.ExecuteToFile(t, project.GetDocByName("client"), fmt.Sprintf("%s/webClient/src/app/components/city", globalDistPath), "item.vue")
	//	utils.CheckErr(err, "ExecuteToFile")
	//}

}

func readData(p types.ProjectType)  {
	project = p
	project.DistPath = "../src"
	project.FillDocTemplatesFields()
	project.GenerateGrid()
}

func Start(p types.ProjectType)  {
	// читаем данные для проекта
	readData(p)
	// читаем темплейты
	tmplMap = templates.ParseTemplates(project)

	// генерим файлы для документов
	for _, d := range p.Docs {
		for _, dt := range d.Templates {
			err := templates.ExecuteToFile(dt.Tmpl, d, dt.DistPath, dt.DistFilename)
			utils.CheckErr(err, fmt.Sprintf("'%s' ExecuteToFile '%s'", d.Name, dt.DistFilename))
		}
	}


	//if tmpls, ok := tmplMap["webClient"]; ok {
	//	t := tmpls.Lookup("docItem.vue")
	//	if t == nil {
	//		log.Fatalf("no found template: docItem.vue")
	//	}
	//	for _, d := range project.Docs {
	//		err := templates.ExecuteToFile(t, d, fmt.Sprintf("%s/webClient/src/app/components/%s", p.DistPath, d.Name), "item.vue")
	//		utils.CheckErr(err, "ExecuteToFile")
	//	}
	//}
}