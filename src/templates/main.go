package templates

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pepelazz/projectGenerator/src/types"
	"github.com/pepelazz/projectGenerator/src/utils"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)


var funcMap = template.FuncMap{
	"ToUpper":        strings.ToUpper,
	"UpperCaseFirst": utils.UpperCaseFirst,
	"ToLowerCamel": strcase.ToLowerCamel,
	"PrintVueFldTemplate": PrintVueFldTemplate,
}

func ParseTemplates(p types.ProjectType) map[string]*template.Template {
	// парсинг общих шаблонов
	res := map[string]*template.Template{}
	readFiles := func(prefix, delimLeft, delimRight string, path ...string) {
		tmpls, err := template.New("").Funcs(funcMap).Delims(delimLeft, delimRight).ParseFiles(path...)
		utils.CheckErr(err, "ParseFiles")
		for _, t := range tmpls.Templates() {
			res[prefix + t.Name()] = t
		}
	}

	// webClient
	path := "../../projectGenerator/src/templates/webClient/doc/"
	readFiles("webClient_", "[[", "]]", path + "index.vue", path + "item.vue")
	// sql
	path = "../../projectGenerator/src/templates/sql/"
	readFiles("sql_", "{{", "}}", path + "main.toml")

	// парсинг шаблонов для конкретного документа
	for i, d := range p.Docs {
		for tName, dt := range d.Templates {
			t, err := template.New(tName).Funcs(funcMap).Delims("[[", "]]").ParseFiles(dt.Source)
			utils.CheckErr(err, fmt.Sprintf("doc: %s tmpl: %s parse template error: %s", d.Name, tName, err))
			// сохраняем template в поле структуры
			dt.Tmpl = t
		}
		// дописываем стандартные шаблоны
		if d.IsBaseTemapltes {
			for _, tName := range []string{"webClient_item.vue", "webClient_index.vue", "sql_main.toml"} {
				// если шаблона с таким именем нет, то добавляем стандартный
				if _, ok := d.Templates[tName]; !ok {
					distPath, distFilename := utils.ParseDocTemplateFilename(d.Name, tName, p.DistPath, i)
					d.Templates[tName]= &types.DocTemplate{Tmpl: res[tName], DistPath: distPath, DistFilename: distFilename}
				}
			}
		}
	}

	return res
}

func ExecuteToFile(t *template.Template, d interface{}, path, filename string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path +"/" + filename, []byte(tpl.String()), 0644)
}

func PrintVueFldTemplate(fld types.FldType) string {
	name := fld.Vue.Name
	if len(name) == 0 {
		name = fld.Name
	}
	nameRu := fld.Vue.NameRu
	if len(nameRu) == 0 {
		nameRu = fld.NameRu
	}
	fldType := fld.Vue.Type
	if len(fldType) == 0 {
		fldType = fld.Type
	}
 	switch fldType {
	case "string":
		return fmt.Sprintf(`<q-input type='text' v-model="item.%s" label="%s" autogrow/>`, name, nameRu)
	default:
		return fmt.Sprintf("not found vueFldTemplate for type `%s`", fld.Type)
	}
}
