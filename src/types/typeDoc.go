package types

import (
	"github.com/serenize/snaker"
	"text/template"
)

const (
	DocTypeLinkTable = "linkTable"
)

type (
	DocType struct {
		Name            string
		NameRu          string
		Type            string
		Flds            []FldType
		Vue             DocVue
		Sql             DocSql
		Templates       map[string]*DocTemplate
		IsBaseTemapltes DocIsBaseTemapltes // флаг что генерируем стандартные шаблоны для документа
	}

	DocVue struct {
		RouteName string
		MenuIcon  string
		Roles     []string
		Grid      []VueGridDiv
		Mixins    map[string][]string          // название файла - название миксина. Для прописывания импорта
		Methods   map[string]map[string]string // название файла - название метода - текст функции
		TmplFuncs map[string]func(DocType) string
		I18n      map[string]string
	}

	// специальное представление для сетки
	VueGridDiv struct {
		Class string
		Grid  []VueGridDiv
		Fld   FldType
	}

	DocTemplate struct {
		Source       string
		DistPath     string
		DistFilename string
		Tmpl         *template.Template
	}

	DocSql struct {
		Methods         map[string]*DocSqlMethod
		IsUniqLink      bool // флаг, что таблица является связью двух таблиц и связь между ними уникальная
		IsBeforeTrigger bool // флаг что добавляем before триггер
		IsAfterTrigger  bool // флаг что добавляем after триггер
		IsSearchText    bool // флаг что добавляем поле search_text
	}

	DocIsBaseTemapltes struct {
		Vue bool
		Sql bool
	}

	DocSqlMethod struct {
		Name  string
		Roles []string
	}
)

// место вызова разных доп функций для инициализации документа, после того как основные поля заполнены
func (d *DocType) Init() {
	d.Filli18n()
}

func (d DocType) PgName() string {
	return snaker.CamelToSnake(d.Name)
}

func (d DocType) NameCamelCase() string  {
	return snaker.SnakeToCamel(d.Name)
}
