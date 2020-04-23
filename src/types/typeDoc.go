package types

import (
	"fmt"
	"strings"
	"text/template"
)

type (
	DocType struct {
		Name string
		NameRu string
		Flds []FldType
		Vue DocVue
		Templates map[string]*DocTemplate
		IsBaseTemapltes bool // флаг что генерируем стандартные шаблоны для документа
	}

	DocVue struct {
		Route string
		Grid []VueGridDiv
		Mixins map[string][]string // название файла - название миксина. Для прописывания импорта
	}

	// специальное представление для сетки
	VueGridDiv struct {
		Class string
		Grid []VueGridDiv
		Fld FldType
	}

	DocTemplate struct {
		Source string
		DistPath string
		DistFilename string
		Tmpl *template.Template
	}
)

func (d DocVue) PrintImport(tmplName string) string  {
	res := []string{}
	if d.Mixins != nil {
		if arr, ok := d.Mixins[tmplName]; ok {
			for _, s := range arr {
				res = append(res, fmt.Sprintf("\timport %s from './mixins/%s'", s, s))
			}
		}
	}

	return strings.Join(res, "\n")
}

func (d DocVue) PrintMixins(tmplName string) string  {
	res := []string{}
	if d.Mixins != nil {
		if arr, ok := d.Mixins[tmplName]; ok {
			for _, s := range arr {
				res = append(res, fmt.Sprintf("%s", s))
			}
		}
	}

	return strings.Join(res, ", ")
}