package types

import (
	"fmt"
	"github.com/serenize/snaker"
	"log"
	"strings"
)

// таблица
type (

	FldVueCompositionTable struct {
		FldName string
		TableName string
		Columns []FldVueCompositionTableColumn
		PgMethod string
	}
	FldVueCompositionTableColumn struct {
		Name string
		Label string
		Align string
		Field string
		Sortable bool
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
	// прописываем пути шаблона
	if d.Templates == nil {
		d.Templates = map[string]*DocTemplate{}
	}
	d.Templates[fmt.Sprintf("%s_common_table", tbl.FldName)] = &DocTemplate{
		Source:       "../../../pepelazz/projectGenerator/templates/webClient/doc/comp/commonTable.vue",
		DistPath:     "../src/webClient/src/app/components/partner/comp",
		DistFilename: tbl.FldName + "CommonTable.vue",
		FuncMap: map[string]interface{}{
			"GetTableTitle": func() string {return tbl.TableName},
			"GetColumns": func() []FldVueCompositionTableColumn { return tbl.Columns},
			"GetPgMethod": func() string {return tbl.PgMethod},
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
	fld = FldType{Type:FldTypeVueComposition,  Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType) string {
		return fmt.Sprintf("<%[1]s :item='item'/>", strings.Replace(snaker.CamelToSnake(tbl.FldName + "CommonTable"), "_", "-", -1))
	}}}

	return
}

//func (fld FldType) SetFromConfigTable(d *DocType, fldName string) FldType {
//	if d.Sql.Hooks.BeforeInsertUpdate == nil {
//		d.Sql.Hooks.BeforeInsertUpdate = []string{}
//	}
//	triggerStr := fmt.Sprintf(`
//		params = params || jsonb_build_object('%[1]s', (select %[2]s from config limit 1));
//		if params->>'%[1]s' isnull then
//			return jsonb_build_object('ok', false, 'message', 'missed %[2]s in "config" table');
//		end if;
//	`, fld.Name, fldName)
//	d.Sql.Hooks.BeforeInsertUpdate = append(d.Sql.Hooks.BeforeInsertUpdate, triggerStr)
//	return fld
//}

//func GetFldVueCompTable(d *DocType, comp FldVueCompositionTable, rowCol [][]int, params... string) (fld FldType) {
//	return getFldVueComposition(d, comp, rowCol, params...)
//}
//
//func getFldVueComposition(d *DocType, comp FldVueCompositionTmp, rowCol [][]int, params... string) (fld FldType) {
//	classStr := "col-md-4 col-xs-6"
//	if len(params)>0 {
//		classStr= params[0]
//	}
//	fld = FldType{Type:FldTypeVueComposition,  Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType) string {
//		return fmt.Sprintf("<%[1]s :item='item'/>", comp.GetName())
//	}}}
//	return
//}
