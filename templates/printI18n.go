package templates

import (
	"fmt"
	"github.com/pepelazz/nla_framework/types"
	"github.com/pepelazz/nla_framework/utils"
	"strings"
	"text/template"
)


// заполняем словарь локализации
func FillDocI18n(p types.ProjectType) {
	for i := range p.Docs {
		for _, lang := range p.I18n.LangList {
			if len(p.Docs[i].I18n)==0 {
				p.Docs[i].I18n = map[string]map[string]string{}
			}
			if len(p.Docs[i].I18n[lang]) == 0 {
				p.Docs[i].I18n[lang] = map[string]string{}
			}
			//NAME
			if _, ok := p.Docs[i].I18n[lang]["name"]; !ok {
				if lang == "ru" {
					p.Docs[i].I18n[lang]["name"] = p.Docs[i].NameRu
				} else {
					p.Docs[i].I18n[lang]["name"] = strings.ReplaceAll(p.Docs[i].Name, "_", " ")
				}
			}
			//NAME_PLURAL
			if _, ok := p.Docs[i].I18n[lang]["name_plural"]; !ok {
				if lang == "ru" {
					p.Docs[i].I18n[lang]["name_plural"] = p.Docs[i].Vue.I18n["listTitle"]
				} else {
					p.Docs[i].I18n[lang]["name_plural"] = strings.ReplaceAll(p.Docs[i].Name, "_", " ")
				}
			}
			//NAME_PLURAL
			if _, ok := p.Docs[i].I18n[lang]["name_plural_deleted"]; !ok {
				if lang == "ru" {
					p.Docs[i].I18n[lang]["name_plural_deleted"] = "удаленные " + p.Docs[i].Vue.I18n["listTitle"]
				} else {
					p.Docs[i].I18n[lang]["name_plural_deleted"] = "deleted " + strings.ReplaceAll(p.Docs[i].Name, "_", " ")
				}
			}
			// MENU
			if _, ok := p.I18n.Data[lang]["menu"][p.Docs[i].Name]; !ok {
				if lang == "en" {
					p.I18n.Data[lang]["menu"][p.Docs[i].Name] = strings.ReplaceAll(p.Docs[i].Name, "_", " ")
				} else {
					p.I18n.Data[lang]["menu"][p.Docs[i].Name] = p.Docs[i].Vue.I18n["listTitle"]
				}
			}
			for _, fld := range p.Docs[i].Flds {
				// проверяем, чтобы не перезаписать уже существующий ключ
				if _, ok := p.Docs[i].I18n[lang][fld.Name]; !ok {
					if len(fld.Name) > 0 && len(fld.NameRu) > 0 {
						if lang == "ru" {
								p.Docs[i].I18n[lang][fld.Name] = fld.NameRu
								if fld.Name == "title" {
									p.Docs[i].I18n[lang][fld.Name] = "название"
								}
						} else {
							// все остальные языки кроме русского
							name := strings.ReplaceAll(fld.Name, "_", " ")
							name = strings.ReplaceAll(name, "id", "")
							p.Docs[i].I18n[lang][fld.Name] = name
						}
					}
				}
			}
		}
	}
}

// ПЕЧАТЬ i18n/index.js
func PrintI18nJs(p types.ProjectType)  {
	resStr := ""
	for _, lang := range p.I18n.LangList {
		dirLang := lang
		if lang == "en" {
			dirLang = "en-US"
		}
		resStr = fmt.Sprintf("%simport %s from './%s'\n", resStr, lang, dirLang)
	}
	resStr = resStr + "\nexport default {"
	for _, lang := range p.I18n.LangList {
		langKey := lang
		if lang == "en" {
			langKey = "en-US"
		} else {
			langKey = lang + "-" + strings.ToUpper(lang)
		}
		resStr = fmt.Sprintf("%s\n\t\t'%s': %s,", resStr, langKey, strings.ReplaceAll(lang, "-", ""))
	}
	resStr = resStr + "\n}\n"

	t, _ := template.New("").Parse(resStr)
	err := ExecuteToFile(t, p,  p.DistPath + "/webClient/src/i18n", "index.js")
	utils.CheckErr(err, "ExecuteToFile template /i18n/index.js")
}

func PrintDocI18nJs(p types.ProjectType, lang string)  {
	// ПЕЧАТЬ doc/index.js
	resStr := ""
	for _, d := range p.Docs {
		resStr = fmt.Sprintf("%simport %s from './%s'\n", resStr, d.Name, d.Name)
	}
	resStr = resStr + "\nexport default {"
	// печатаем сообщения на уровне проекта
	if _, ok := p.I18n.Data[lang]; ok {
		for m, list := range p.I18n.Data[lang] {
			resStr = fmt.Sprintf("%s\n	%s: {", resStr, m)
			for k, v := range list {
				resStr = fmt.Sprintf("%s\n 		%s: '%s',", resStr, k, v)
			}
			resStr = resStr + "\n	},"
		}
	}
	// ссылки на документы
	for _, d := range p.Docs {
		resStr = fmt.Sprintf("%s\n	%s,", resStr, d.Name)
	}
	resStr = resStr + "\n}\n"

	t, _ := template.New("").Parse(resStr)
	langDirName := lang
	if langDirName == "en" {
		langDirName = "en-US"
	}
	err := ExecuteToFile(t, p,  p.DistPath + "/webClient/src/i18n/" + langDirName, "index.js")
	utils.CheckErr(err, fmt.Sprintf("ExecuteToFile template /i18n/%s/index.js", langDirName))

	// ПЕЧАТЬ user.js и прочих docs
	for _, doc := range p.Docs {
		resStr := ""
		resStr = resStr + "\nexport default {"
		if _, ok := doc.I18n[lang]; ok {
			for k, v := range doc.I18n[lang] {
				resStr = fmt.Sprintf("%s\n 		%s: '%s',", resStr, k, v)
			}
		}
		resStr = resStr + "\n}\n"

		t, _ := template.New("").Parse(resStr)
		err = ExecuteToFile(t, p,  p.DistPath + "/webClient/src/i18n/" + langDirName, doc.Name + ".js")
		utils.CheckErr(err, fmt.Sprintf("ExecuteToFile template /i18n/%s/%s.js", langDirName, doc.Name))
	}


}
