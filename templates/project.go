package templates

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/templates/tmplGenerateStep2"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"strings"
	"text/template"
)

func WriteProjectFiles(p types.ProjectType, tmplMap map[string]*template.Template)  {
	for name, t := range tmplMap {
		if strings.HasPrefix(name, "project_") {
			filename := strings.TrimPrefix(name, "project_")
			path := ".."
			if filename == "config.toml" || filename == "main.go" {
				path = "../src"
			}
			err := ExecuteToFile(t, p, path, filename)
			utils.CheckErr(err, fmt.Sprintf("'project' ExecuteToFile '%s'", name))
		}
	}

	// генерим шаблоны, которые указаны дополнительно на уровне проекта. Без относительно конкретных документов
	for _, m := range p.Sql.Methods {
		for _, v := range m {
			if len(v.Tmpl.Source) > 0 && len(v.Tmpl.Dist) > 0 {
				distPath, filename := utils.PathExtractFilename(v.Tmpl.Dist)
				distPath = "../src" + distPath
				t, err := template.New(filename).Delims("[[", "]]").ParseFiles(v.Tmpl.Source)
				utils.CheckErr(err, "p.Sql.Methods")

				err = ExecuteToFile(t, p, distPath, filename)
				utils.CheckErr(err, fmt.Sprintf("'project' ExecuteToFile '%s'", filename))
			}
		}
	}

	projectTmplPath := "../../../pepelazz/projectGenerator/templates/project"
	readTmplAndPrint(p, projectTmplPath + "/types/main.go", "/types",  "main.go", nil)
	readTmplAndPrint(p, projectTmplPath + "/types/config.go", "/types",  "config.go", nil)
	readTmplAndPrint(p, projectTmplPath + "/webServer/main.go", "/webServer",  "main.go", nil)
	readTmplAndPrint(p, projectTmplPath + "/sql/initialData.sql", "/sql/template/function/",  "initialData.sql", nil)
	readTmplAndPrint(p, projectTmplPath + "/sql/user_trigger_after.sql", "/sql/template/function/_User/",  "user_trigger_after.sql", template.FuncMap{"PrintUserAfterTriggerUpdateLinkedRecords": types.PrintUserAfterTriggerUpdateLinkedRecords})
	readTmplAndPrint(p, projectTmplPath + "/sql/01_User/main.toml", "/sql/model/01_User",  "main.toml", nil)
	readTmplAndPrint(p, projectTmplPath + "/jobs/main.go", "/jobs",  "main.go", nil)
	readTmplAndPrint(p, projectTmplPath + "/pg/pgListener.go", "/pg",  "pgListener.go", nil)
	readTmplAndPrint(p, projectTmplPath + "/webClient/app/components/users/roles.js", "/webClient/src/app/components/users",  "roles.js", nil)
	readTmplAndPrint(p, projectTmplPath + "/webClient/app/components/users/item.vue", "/webClient/src/app/components/users",  "item.vue", nil)

	if p.IsTelegramIntegration() {
		readTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/telegram/telegramAuth.go", "/webServer", "telegramAuth.go", nil)
		readTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/telegram/user_telegram_auth.sql", "/sql/template/function/_User", "user_telegram_auth.sql", nil)
		readTmplAndPrint(p, projectTmplPath + "/tgBot/main.go", "/tgBot", "main.go", nil)
	}

	// в случае коннекта к Битрикс генерим файлы
	if p.IsBitrixIntegration() {
		readTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/bitrix/bitrixMain.go", "/bitrix", "main.go", nil)
		//sourcePath := "../../../pepelazz/projectGenerator/templates/integrations/bitrix/bitrixMain.go"
		//t, err := template.New("bitrixMain.go").Funcs(funcMap).Delims("[[", "]]").ParseFiles(sourcePath)
		//utils.CheckErr(err, "bitrixMain.go")
		//distPath := fmt.Sprintf("%s/bitrix", p.DistPath)
		//err = ExecuteToFile(t, p, distPath, "main.go")
		//utils.CheckErr(err, fmt.Sprintf("'project' ExecuteToFile '%s'", "bitrix/main.go"))
	}

	// в случае коннекта к 1 Odata генерим файлы
	if p.IsOdataIntegration() {
		readTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/odata/main.go", "/odata", "main.go", nil)
		readTmplAndPrint(p, "../../../pepelazz/projectGenerator/templates/integrations/odata/odataQueryType.go", "/odata", "odataQueryType.go", nil)
	}
}

func OtherTemplatesGenerate(p types.ProjectType)  {
	tmplGenerateStep2.TasksTmpl(p)
	// добавляем функции в plugin/utils.js
	tmplGenerateStep2.PluginUtilsJs(p)
}

func readTmplAndPrint(p types.ProjectType, sourcePath, distPath, filename string, addFuncMap template.FuncMap) {
	fMap := template.FuncMap{}
	for k, v := range funcMap {
		fMap[k] = v
	}
	for k, v := range addFuncMap {
		fMap[k] = v
	}
	_, sourceFilename := utils.PathExtractFilename(sourcePath)
	t, err := template.New(sourceFilename).Funcs(fMap).Delims("[[", "]]").ParseFiles(sourcePath)
	utils.CheckErr(err, "readFileWithDist")
	err = ExecuteToFile(t, p, p.DistPath + distPath, filename)
	utils.CheckErr(err, fmt.Sprintf("readTmplAndPrint ExecuteToFile '%s/%s'", distPath, filename))
}

