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

	// в случае коннекта к Битрикс генерим файлы
	if len(p.Config.Bitrix.ApiUrl) > 0 {
		sourcePath := "../../../pepelazz/projectGenerator/templates/integrations/bitrix/bitrixMain.go"
		t, err := template.New("bitrixMain.go").Funcs(funcMap).Delims("[[", "]]").ParseFiles(sourcePath)
		utils.CheckErr(err, "bitrixMain.go")
		distPath := fmt.Sprintf("%s/bitrix", p.DistPath)
		err = ExecuteToFile(t, p, distPath, "main.go")
		utils.CheckErr(err, fmt.Sprintf("'project' ExecuteToFile '%s'", "bitrix/main.go"))

	}
}

func OtherTemplatesGenerate(p types.ProjectType)  {
	tmplGenerateStep2.TasksTmpl(p)
	// добавляем функции в plugin/utils.js
	tmplGenerateStep2.PluginUtilsJs(p)
}

