package types

import (
	"fmt"
	"github.com/serenize/snaker"
	"strings"
	"text/template"
)

type (
	DocType struct {
		Name string
		NameRu string
		Flds []FldType
		Vue DocVue
		Sql DocSql
		Templates map[string]*DocTemplate
		IsBaseTemapltes DocIsBaseTemapltes // флаг что генерируем стандартные шаблоны для документа
	}

	DocVue struct {
		RouteName string
		Grid      []VueGridDiv
		Mixins    map[string][]string // название файла - название миксина. Для прописывания импорта
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

	DocSql struct {
		Methods map[string]*DocSqlMethod
		IsUniqLink 		bool // флаг, что таблица является связью двух таблиц и связь между ними уникальная
		IsBeforeTrigger	bool // флаг что добавляем before триггер
		IsAfterTrigger	bool // флаг что добавляем after триггер
		IsSearchText	bool // флаг что добавляем поле search_text
	}

	DocIsBaseTemapltes struct {
		Vue bool
		Sql bool
	}

	DocSqlMethod struct {
		Name string
		Roles []string
	}
)

func (d DocType) PgName() string  {
	return snaker.CamelToSnake(d.Name)
}

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