package types

import (
	"encoding/json"
	"fmt"
	"github.com/pepelazz/projectGenerator/utils"
	"github.com/spf13/cast"
	"log"
	"strings"
)

func (d DocType) PrintListRowAvatar() string  {
	res := fmt.Sprintf(`
        <q-item-section avatar @click="$router.push(currentUrl + item.id)" style="cursor: pointer">
          <q-avatar rounded>
            <img src="%s" alt="">
          </q-avatar>
        </q-item-section>
	`, d.Vue.MenuIcon)
	// проверяем есть ли переопределение шаблона
	if d.Vue.TmplFuncs != nil {
		if f, ok := d.Vue.TmplFuncs["PrintListRowAvatar"]; ok {
			res = f(d)
		}
	}
	return res
}

func (d DocType) PrintListRowLabel() string  {
	isFolder := "" // признак, что является parent, в случае если рекурсия
	if d.IsRecursion {
		isFolder = "<q-item-label caption><q-icon name='folder' v-if='item.is_folder'/></q-item-label>"
	}
	clickOpenItem := ""
	if d.Vue.IsVueTitleClickable {
		clickOpenItem = fmt.Sprintf(` @click="$router.push(currentUrl + item.id)" style="cursor: pointer" `, )
	}
	res := fmt.Sprintf(`
        <q-item-section>
          <q-item-label lines="1" %s>{{item.title}}</q-item-label>
          %s
        </q-item-section>
	`, clickOpenItem, isFolder)
	// проверяем есть ли переопределение шаблона
	if d.Vue.TmplFuncs != nil {
		if f, ok := d.Vue.TmplFuncs["PrintListRowLabel"]; ok {
			res = f(d)
		}
	}
	return res
}

func (d *DocType) AddVueMethod(tmplName, methodName, method string)  {
	if d.Vue.Methods == nil {
		d.Vue.Methods = map[string]map[string]string{}
	}
	if _, ok := d.Vue.Methods[tmplName]; !ok {
		d.Vue.Methods[tmplName] = map[string]string{}
	}
	d.Vue.Methods[tmplName][methodName] = method
}

func (d *DocType) Filli18n() {
	if d.Vue.I18n == nil {
		d.Vue.I18n = map[string]string{}
	}
	if _, ok := d.Vue.I18n["listTitle"]; !ok {
		d.Vue.I18n["listTitle"] = utils.UpperCaseFirst(d.NameRu)
	}
	if _, ok := d.Vue.I18n["listDeletedTitle"]; !ok {
		d.Vue.I18n["listDeletedTitle"] = "Удаленные " + d.NameRu
	}
}

func (d DocType) PrintVueItemOptionsFld() string  {
	arr := []string{}
	for _, fld := range d.Flds {
		if fld.Sql.IsOptionFld {
			arr = append(arr, fld.Name)
		}
	}
	if len(arr) > 0 {
		return fmt.Sprintf("'%s'", strings.Join(arr, "', '"))
	}
	return ""
}

func (d DocType) PrintVueImport(tmplName string) string  {
	isLodashAdded := false
	res := []string{}
	// ссылки на миксины
	if d.Vue.Mixins != nil {
		if arr, ok := d.Vue.Mixins[tmplName]; ok {
			for _, m := range arr {
				res = append(res, fmt.Sprintf("\timport %s from '%s'", m.Title, m.Import))
			}
		}
	}
	// ссылки на компоненты
	if d.Vue.Components != nil {
		if m, ok := d.Vue.Components[tmplName]; ok {
			for name, path := range m {
				res = append(res, fmt.Sprintf("\timport %s from '%s'", name, path))
			}
		}
	}

	if tmplName == "docItem" {
		// если есть поле с типом multipleSelect то добавляем lodash
		for _, fld := range d.Flds {
			if fld.Vue.Type == FldVueTypeMultipleSelect && !isLodashAdded{
				res = append(res, "\timport _ from 'lodash'")
				isLodashAdded = true
				break
			}
		}
		if d.IsRecursion {
			parentPath := "."
			// в случае табов папка с компонентами на уровень выше
			if len(d.Vue.Tabs) > 0 {
				parentPath = "../.."
			}
			res = append(res, fmt.Sprintf("\timport compRecursiveChildList from '%s/comp/recursiveChildList'", parentPath))
		}
	}

	if tmplName == "docItemWithTabs" {
		for _, tab := range d.Vue.Tabs {
			res = append(res, fmt.Sprintf("\timport %[1]sTab from './tabs/%[1]s/index'", tab.Title))
		}
	}

	return strings.Join(res, "\n")
}

func (d DocType) PrintVueMethods(tmplName string) string  {
	// извлекаем map с методами из описания документа
	methods := map[string]string{}
	if d.Vue.Methods != nil {
		if m , ok := d.Vue.Methods[tmplName]; ok {
			// копируем методы в map, который указан выше
			for k, v := range m {
				methods[k] = v
			}
		}
	}
	// печатаем методы
	res := ""
	for mName, funcTxt := range methods {
		res = fmt.Sprintf("%s\t%s {\n\t\t\t\t%s\n\t\t\t\t\t\t},\n", res, mName, funcTxt)
	}
	return res
}

func (d DocType) PrintVueItemHookBeforeSave() string  {
	res := ""
	for _, v := range d.Vue.Hooks.ItemBeforeSave {
		res = res + fmt.Sprintf("%s\n", v)
	}
	return res
}

func (d DocType) PrintVueItemForSave() string {
	res := ""
	for _, fld := range d.Flds {
		if fld.Vue.Type == FldVueTypeSelect {
			res = fmt.Sprintf("%s%[2]s: this.item.%[2]s ? this.item.%[2]s.value : null,\n", res, fld.Name)
		}
		if fld.Vue.Type == FldVueTypeMultipleSelect {
			res = fmt.Sprintf("%s%[2]s: this.item.%[2]s ? this.item.%[2]s.map(({value}) => value).filter(v => v)  : [],\n", res, fld.Name)
		}
	}
	for _, v := range d.Vue.Hooks.ItemForSave {
		res = fmt.Sprintf("%s%s,\n", res, v)
	}
	if d.IsRecursion {
		res = fmt.Sprintf("%sparent_id: this.parent_id ? +this.parent_id : null,\n", res)
	}
	return res
}

func (d DocType) PrintVueItemResultModify() string {
	res := ""
	for _, v := range d.Vue.Hooks.ItemModifyResult {
		res = res + fmt.Sprintf("%s\n", v)
	}
	for _, fld := range d.Flds {
		// single select - преобразуем v -> {label: label, value: v}
		if fld.Vue.Type == FldVueTypeSelect {
			options, err := json.Marshal(fld.Vue.Options)
			utils.CheckErr(err, fmt.Sprintf("'%s' json.Marshal(fld.Vue.Options)", fld.Name))
			funcStr := fmt.Sprintf(`
				if (res.%[1]s) {
                    let arr = %[2]s
                    let %[1]s_item = arr.find(v => v.value === res.%[1]s)
                    if (%[1]s_item) res.%[1]s = {value: res.%[1]s, label: %[1]s_item.label}
                    }
			`, fld.Name, options)
			res = fmt.Sprintf("%s%s", res, funcStr)
		}
		// multiple select - преобразуем v -> {label: label, value: v}
		if fld.Vue.Type == FldVueTypeMultipleSelect {
			options, err := json.Marshal(fld.Vue.Options)
			utils.CheckErr(err, fmt.Sprintf("'%s' json.Marshal(fld.Vue.Options)", fld.Name))
			funcStr := fmt.Sprintf(`
				if (res.%[1]s) {
                    let arr = %[2]s
					res.%[1]s = res.%[1]s.map(name => _.find(arr, {value: name})).filter(v => v)
                    }
			`, fld.Name, options)
			res = fmt.Sprintf("%s%s", res, funcStr)
		}
	}
	// в случае рекурсии добавляем вычисление форматирование parentProductBreadcrumb
	if d.IsRecursion {
		str := fmt.Sprintf("if (res.parent_title) this.parentProductBreadcrumb = [{label: res.parent_title, to: `${res.parent_id}`, docType: '%s'}]", d.Name)
		res = fmt.Sprintf("%s\n%s", res, str)
	}
	return res
}

// функция преобразования select полей внтури карточки в state machine
func (d DocType) PrintVueItemStateMachineCardMounted() string {
	res := ""
	for _, fld := range d.Flds {
		// single select - преобразуем v -> {label: label, value: v}
		if fld.Vue.Type == FldVueTypeSelect {
			funcStr := fmt.Sprintf(`
				if (this.item.%[1]s && this.$utils._.isString(this.item.%[1]s) && !this.is_current_state) {
					if (this.item.%[1]s) this.item.%[1]s  = {label: this.$utils.i18n_%[2]s_%[1]s(this.item.%[1]s), value: this.item.%[1]s}
                    }
			`, fld.Name, d.Name)
			res = fmt.Sprintf("%s%s", res, funcStr)
		}
		// multiple select - преобразуем v -> {label: label, value: v}
		if fld.Vue.Type == FldVueTypeMultipleSelect {
			options, err := json.Marshal(fld.Vue.Options)
			utils.CheckErr(err, fmt.Sprintf("'%s' json.Marshal(fld.Vue.Options)", fld.Name))
			funcStr := fmt.Sprintf(`
				if (res.%[1]s)) {
                    let arr = %[2]s
					res.%[1]s = res.%[1]s.map(name => _.find(arr, {value: name})).filter(v => v)
                    }
			`, fld.Name, options)
			res = fmt.Sprintf("%s%s", res, funcStr)
		}
	}
	return res
}

func (d DocVue) PrintMixins(tmplName string) string  {
	res := []string{}
	if d.Mixins != nil {
		if arr, ok := d.Mixins[tmplName]; ok {
			for _, s := range arr {
				res = append(res, fmt.Sprintf("%s", s.Title))
			}
		}
	}

	return strings.Join(res, ", ")
}

func (d DocType) PrintComponents(tmplName string) string  {
	res := []string{}
	if d.Vue.Components != nil {
		if m, ok := d.Vue.Components[tmplName]; ok {
			for name := range m {
				res = append(res, name)
			}
		}
	}

	if tmplName ==  "docItemWithTabs" {
		for _, t := range d.Vue.Tabs {
			res = append(res, t.Title + "Tab")
		}
	}

	if d.IsRecursion && tmplName != "docItemWithTabs" {
		res = append(res, "compRecursiveChildList")
	}

	return strings.Join(res, ", ")
}

func GetVueCompLinkListWidget (p ProjectType, d DocType, tableName string, opts map[string]interface{}) string {
	var tableIdFldName, tableDependName, tableDependFldName, tableDependRoute, label, avatarSrc string
	tableIdName := d.Name
	linkTableName := tableName
	hideCreateNew := false
	// находим документы, на которые идет ссылка
	for _, doc := range p.Docs {
		if doc.Name == tableName {
			for _, f := range doc.Flds {
				if len(f.Sql.Ref) > 0 && f.Sql.Ref != d.Name {
					tableDependFldName = f.Name
					// ссылку на таблицу user рассматриваем отдельно, потому что ее нет в списке p.Docs
					if f.Sql.Ref == "user" {
						tableDependName = "user"
						tableDependRoute = "users"
						label = "сотрудники"
						if opts != nil {
							if v, ok := opts["listTitle"]; ok {
								label = cast.ToString(v)
							}
						}
					} else {
						depDoc := p.GetDocByName(f.Sql.Ref)
						if depDoc == nil {
							log.Fatalf(fmt.Sprintf("GetVueCompLinkListWidget not found '%s'", f.Sql.Ref))
						}
						tableDependName = depDoc.Name
						tableDependRoute = depDoc.Vue.RouteName
						avatarSrc = depDoc.Vue.MenuIcon
						label = depDoc.Vue.I18n["listTitle"]
					}
				} else if len(f.Sql.Ref) > 0 {
					tableIdFldName = f.Name
				}
			}
		}
	}
	fldsProp := ""
	slotOtherFlds := ""
	readonly := "false"
	if opts != nil {
		// убираем кнопку 'создать'
		if v, ok := opts["hideCreateNew"]; ok {
			hideCreateNew = cast.ToBool(v)
		}
		if v, ok := opts["flds"]; ok {
			fldsProp = fmt.Sprintf("\n\t\t\t\t\t\t:flds= \"%s\"", v)
		}
		if v, ok := opts["slotOtherFlds"]; ok {
			slotOtherFlds = fmt.Sprintf("\n\t\t\t\t\t\t<template v-slot:otherFlds='slotProps'>\n\t\t\t\t\t\t\t%s\n\t\t\t\t\t\t\t</template>", v)
		}
		if v, ok := opts["tableDependRoute"]; ok {
			tableDependRoute = cast.ToString(v)
		}
		if v, ok := opts["readonly"]; ok {
			readonly = cast.ToString(v)
		}
	}

	return fmt.Sprintf("<comp-link-list-widget label='%s' :id='id' tableIdName='%s' tableIdFldName='%s' tableDependName='%s' tableDependFldName='%s' tableDependRoute='/%s' linkTableName='%s' avatarSrc='%s' :hideCreateNew='%v' :readonly='%s' %s>%s</comp-link-list-widget>", label, tableIdName, tableIdFldName, tableDependName, tableDependFldName, tableDependRoute, linkTableName, avatarSrc, hideCreateNew, readonly, fldsProp, slotOtherFlds)
}

// заголовки табов
func (d DocType) PrintVueItemTabs()  string{
	res := []string{}
	sep := "\n\t\t\t\t\t\t\t\t"
	for _, tab := range d.Vue.Tabs {
		if len(tab.HtmlInner) == 0 {
			res = append(res, fmt.Sprintf("<q-tab name='%s'  icon='%s' label='%s'/>", tab.Title, tab.Icon, tab.TitleRu))
		} else {
			res = append(res, fmt.Sprintf("<q-tab name='%s'  icon='%s' label='%s'>%[5]s\t%[4]s%[5]s</q-tab>", tab.Title, tab.Icon, tab.TitleRu, tab.HtmlInner, sep))
		}
	}
	return strings.Join(res, sep)
}

// список компонентов для отображения содержимого табов
func (d DocType) PrintVueItemTabPanels()  string{
	res := []string{}
	params := ":id='id' :isOpenInDialog='isOpenInDialog' @updated='v=>$emit(`updated`, v)'"
	if d.IsRecursion {
		params = ":id='id' :isOpenInDialog='isOpenInDialog' :parent_id='parent_id'"
	}
	for _, tab := range d.Vue.Tabs {
		res = append(res, fmt.Sprintf("<!-- %s       -->\n\t\t\t\t\t\t\t\t<q-tab-panel name='%[2]s'><%[2]s-tab %[3]s %[4]s/></q-tab-panel>", tab.TitleRu, tab.Title, params, tab.HtmlParams))
	}
	return strings.Join(res, "\n\t\t\t\t\t\t\t\t")
}
