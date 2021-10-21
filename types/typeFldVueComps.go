package types

import (
	"fmt"
	"github.com/serenize/snaker"
	"log"
	"strings"
)

// пример
//doc.AddFld(t.GetFldVueCompositionTable(&doc, t.FldVueCompositionTable{
//	FldName: "parcel_table",
//	TableName: "Список посылок",
//	Columns: []t.FldVueCompositionTableColumn{
//	{Name: "id", Label: "ID"},
//	{Name: "title", Label: "название", Sortable: true},
//	},
//	PgMethod: `{method: 'parcel_list', params: {}}`,
//	Pagination: t.FldVueCompositionTablePagination{
//	RowsPerPage: 10,
//	},
//	Separator: "cell",
//}, [][]int{{2, 1}}, "col-8"))

// таблица
type (

	FldVueCompositionTable struct {
		FldName string
		TableName string
		Columns []FldVueCompositionTableColumn
		PgMethod string
		Pagination FldVueCompositionTablePagination
		Separator string // horizontal (default), vertical, cell, none
	}
	FldVueCompositionTableColumn struct {
		Name string
		Label string
		Align string
		Field string
		Sortable bool
	}
	FldVueCompositionTablePagination struct {
		RowsPerPage int
	}

)

func GetFldVueCompositionTable(d *DocType, tbl FldVueCompositionTable, rowCol [][]int, params... string) (fld FldType) {
	if len(tbl.FldName) == 0 {
		log.Fatalf("doc: '%s'. Missed FldName in FldVueCompositionTable", d.Name)
	}
	if len(tbl.PgMethod) == 0 {
		log.Fatalf("doc: '%s'. Missed PgMethod in FldVueCompositionTable", d.Name)
	}
	// если в snake стиле название, то переводим в camel
	if strings.Contains(tbl.FldName, "_") {
		tbl.FldName = snaker.SnakeToCamelLower(tbl.FldName)
	}
	// проставляем дефолты в columns
	for i, col := range tbl.Columns {
		if len(col.Name) == 0 {
			log.Fatalf("doc: '%s'. Missed column name in FldVueCompositionTable", d.Name)
		}
		if len(col.Field) == 0 {
			col.Field = col.Name
		}
		if len(col.Align) == 0 {
			col.Align = "left"
		}
		tbl.Columns[i] = col
	}
	if tbl.Pagination.RowsPerPage == 0 {
		tbl.Pagination.RowsPerPage = 5
	}
	if len(tbl.Separator)== 0{
		tbl.Separator = "horizontal"
	}
	// прописываем пути шаблона
	if d.Templates == nil {
		d.Templates = map[string]*DocTemplate{}
	}
	d.Templates[fmt.Sprintf("%s_common_table", tbl.FldName)] = &DocTemplate{
		Source:       fmt.Sprintf("%s/templates/webClient/quasar_%v/doc/comp/commonTable.vue", getRootDirPath(), d.GetProject().GetQuasarVersion()),
		DistPath:     "../src/webClient/src/app/components/partner/comp",
		DistFilename: tbl.FldName + "CommonTable.vue",
		FuncMap: map[string]interface{}{
			"GetTableTitle": func() string {return tbl.TableName},
			"GetColumns": func() []FldVueCompositionTableColumn { return tbl.Columns},
			"GetPgMethod": func() string {return tbl.PgMethod},
			"GetRowsPerPage": func() int {return tbl.Pagination.RowsPerPage },
			"GetSeparator": func() string {return tbl.Separator },
		},
	}

	// добавляем в список компонент
	parentPath := "./comp/"
	if len(d.Vue.Tabs) > 0 {
		parentPath = "../../comp/"
	}
	importAddress := parentPath + tbl.FldName + "CommonTable.vue"
	if d.Vue.Components == nil {
		d.Vue.Components = map[string]map[string]string{}
	}
	if d.Vue.Components["docItem"] == nil {
		d.Vue.Components["docItem"] = map[string]string{}
	}
	d.Vue.Components["docItem"][tbl.FldName + "CommonTable"] = importAddress

	// параметры самого поля
	classStr := "col-md-4 col-xs-6"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Type:FldTypeVueComposition,  Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType, fld FldType) string {
		return fmt.Sprintf("<%[1]s :item='item'/>", strings.Replace(snaker.CamelToSnake(tbl.FldName + "CommonTable"), "_", "-", -1))
	}}}

	return
}