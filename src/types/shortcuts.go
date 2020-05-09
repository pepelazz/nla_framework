package types

import (
	"fmt"
	"strings"
)

// создание простого поля Double
func GetFldTitle(params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType {Name: "title", NameRu: "название", Type: FldTypeString, Sql: FldSql{IsRequired: true, IsUniq: true, IsSearch:true, Size:150}, Vue: FldVue{RowCol: [][]int{{1, 1}}, Class: []string{classStr}}}
	return
}

func GetFldDouble(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeDouble, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля String
func GetFldString(name, nameRu string, size int, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Readonly:readonly}}
	if size > 0 {
		fld.Sql.Size = size
	}
	return
}

// создание простого поля String
func GetFldDate(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeDate, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля Int
func GetFldInt(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeInt, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля Ref
func GetFldRef(name, nameRu, refTable string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeInt,  Sql: FldSql{Ref: refTable, IsSearch:true}, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// поле с кастомной композицией
func GetFldJsonbComposition(name, nameRu string, rowCol [][]int, classStr, compName string, params ...string) (fld FldType) {
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb,  Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType) string {
		return fmt.Sprintf("<%[1]s :fld='item.%[2]s' @update='item.%[2]s = $event' label='%[3]s' %[4]s/>", compName, name, nameRu, strings.Join(params, " "))
	}}}
	return
}

// простое html поле
func GetFldSimpleHtml(rowCol [][]int, classStr, htmlStr string) (fld FldType) {
	fld = FldType{Type:FldTypeVueComposition,  Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType) string {
		return htmlStr
	}}}
	return
}

// создание простого поля Select с типом string
func GetFldSelectString(name, nameRu string, size int, rowCol [][]int, options []FldVueOptionsItem, params ...string) (fld FldType) {
	classStr := "col-4"
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeSelect, Class: []string{classStr}, Readonly:readonly, Options:options}}
	if size > 0 {
		fld.Sql.Size = size
	}
	return
}

// создание поля-виджета со связями многие-к-многим
func GetFldLinkListWidget(linkTable string, rowCol [][]int, classStr string, opts map[string]interface{}) (fld FldType) {
	return FldType{Type: FldTypeVueComposition,  Vue: FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType) string {
		return GetVueCompLinkListWidget(p, d, linkTable, opts)
	}}}
}


// функция конвертации списка имен файлов с шаблонами в  map[string]*DocTemplate
func GetCustomTemplates(p ...string) map[string]*DocTemplate  {
	res := map[string]*DocTemplate{}
	for _, name := range p {
		res[name] = &DocTemplate{}
	}
	return res
}
