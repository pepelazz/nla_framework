package tmplGenerateStep2

import (
	"fmt"
	"github.com/pepelazz/nla_framework/types"
	"github.com/pepelazz/nla_framework/utils"
	"strings"
	"text/template"
)

// шаблоны для сообщений о задачах, i18n и пр
func PluginUtilsJs(p types.ProjectType)  {
	distPath := fmt.Sprintf("%s/webClient/src/app/plugins", p.DistPath)

	funcNames, funcBodyes := getI18nForSelectFlds(p)

	funcMap := template.FuncMap{
		"ExportDefaultList": func() string {
			return funcNames
		},
		"FunctionsList": func() string {
			return funcBodyes
		},
	}
	path := fmt.Sprintf("../../../pepelazz/nla_framework/templates/project/webClient/quasar_%v/app/plugins/utils.js", p.GetQuasarVersion())
	t, err := template.New("utils.js").Funcs(funcMap).Delims("[[", "]]").ParseFiles(path)
	utils.CheckErr(err, "OverriteCopiedFiles ParseFiles")

	err = executeToFile(t, "", distPath, "utils.js")
	utils.CheckErr(err, "OverriteCopiedFiles ExecuteToFile")
}

func getI18nForSelectFlds(p types.ProjectType) (funcNames, funcBodyes string) {
	for _, d := range p.Docs {
		for _, fld := range d.Flds {
			if fld.Vue.Type == types.FldVueTypeSelect || fld.Vue.Type == types.FldVueTypeMultipleSelect  || fld.Vue.Type == types.FldVueTypeRadio {
				// флаг, для определения указан ли цвет
				isColor := false
				fNname := fmt.Sprintf("i18n_%s_%s", d.Name, fld.Name)
				// название функции
				funcNames = fmt.Sprintf("%s%s,\n\t", funcNames, fNname)
				arr := []string{}
				for _, v := range fld.Vue.Options {
					arr = append(arr, fmt.Sprintf("%s: '%s'", v.Value, v.Label))
					if len(v.Color)>0{
						isColor = true
					}
				}
				funcBodyes = fmt.Sprintf(`%s
const %s = (v) => {
	const d = {
		%s
	}
	return Array.isArray(v) ? v.map(v1 => d[v1]) : d[v]
}
				`, funcBodyes, fNname, strings.Join(arr, ",\n\t\t"))

				// если указан цвет, то делаем доп функцию по переводу названия в цвет
				if isColor {
					fNname = fmt.Sprintf("color_%s_%s", d.Name, fld.Name)
					// название функции
					funcNames = fmt.Sprintf("%s%s,\n\t", funcNames, fNname)
					arr := []string{}
					for _, v := range fld.Vue.Options {
						arr = append(arr, fmt.Sprintf("%s: '%s'", v.Value, v.Color))
					}
					funcBodyes = fmt.Sprintf(`%s
const %s = (v) => {
	const d = {
		%s
	}
	return Array.isArray(v) ? v.map(v1 => d[v1]) : d[v]
}
				`, funcBodyes, fNname, strings.Join(arr, ",\n\t\t"))
				}
			}
		}
		// в документе может быть прописан дополнительный глобальный справочник
		for fNname, m := range d.Vue.GloablI18n {
			// название функции
			funcNames = fmt.Sprintf("%s%s,\n\t", funcNames, fNname)
			arr := []string{}
			for val, label := range m {
				arr = append(arr, fmt.Sprintf("%s: '%s'", val, label))
			}
			funcBodyes = fmt.Sprintf(`%s
const %s = (v) => {
	const d = {
		%s
	}
	return Array.isArray(v) ? v.map(v1 => d[v1]) : d[v]
}
				`, funcBodyes, fNname, strings.Join(arr, ",\n\t\t"))
		}
	}
	return
}