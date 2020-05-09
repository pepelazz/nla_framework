package types

import (
	"encoding/json"
	"fmt"
	"github.com/pepelazz/projectGenerator/src/utils"
	"github.com/spf13/cast"
	"log"
	"strings"
)

func (d DocType) PrintListRowLabel() string  {
	res := `
        <q-item-section>
          <q-item-label lines="1">{{item.title}}</q-item-label>
        </q-item-section>
	`
	// проверяем есть ли переопределение шаблона
	if d.Vue.TmplFuncs != nil {
		if f, ok := d.Vue.TmplFuncs["PrintListRowLabel"]; ok {
			res = f(d)
		}
	}
	return res
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
			for _, s := range arr {
				res = append(res, fmt.Sprintf("\timport %s from '../../mixins/%s'", s, s))
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
	}

	if tmplName == "docItemWithTabs" {
		for _, tab := range d.Vue.Tabs {
			res = append(res, fmt.Sprintf("\timport %[1]sTab from './tabs/%[1]s'", tab.Title))
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

func (d DocType) PrintVueItemForSave() string {
	res := ""
	for _, fld := range d.Flds {
		if fld.Vue.Type == FldVueTypeSelect {
			res = fmt.Sprintf("%s%[2]s: this.item.%[2]s ? this.item.%[2]s.value : undefined,\n", res, fld.Name)
		}
		if fld.Vue.Type == FldVueTypeMultipleSelect {
			res = fmt.Sprintf("%s%[2]s: this.item.%[2]s ? this.item.%[2]s.map(({value}) => value).filter(v => v)  : [],\n", res, fld.Name)
		}
	}
	return res
}

func (d DocType) PrintVueItemResultModify() string {
	res := ""
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
	return res
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

func (d DocVue) PrintComponents(tmplName string) string  {
	res := []string{}
	if d.Components != nil {
		if m, ok := d.Components[tmplName]; ok {
			for name := range m {
				res = append(res, name)
			}
		}
	}

	if tmplName ==  "docItemWithTabs" {
		for _, t := range d.Tabs {
			res = append(res, t.Title + "Tab")
		}
	}

	return strings.Join(res, ", ")
}

func GetVueCompLinkListWidget (p ProjectType, d DocType, tableName string, opts map[string]interface{}) string {
	var tableIdFldName, tableDependName, tableDependFldName, tableDependRoute, label, avatarSrc string
	tableIdName := d.Name
	linkTableName := tableName
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
				} else {
					tableIdFldName = f.Name
				}
			}
		}
	}

	return fmt.Sprintf("<comp-link-list-widget label='%s' :id='id' tableIdName='%s' tableIdFldName='%s' tableDependName='%s' tableDependFldName='%s' tableDependRoute='/%s' linkTableName='%s' avatarSrc='%s'/>", label, tableIdName, tableIdFldName, tableDependName, tableDependFldName, tableDependRoute, linkTableName, avatarSrc)
}

// заголовки табов
func (d DocVue) PrintItemTabs()  string{
	res := []string{}
	sep := "\n\t\t\t\t\t\t\t\t"
	for _, tab := range d.Tabs {
		if len(tab.HtmlInner) == 0 {
			res = append(res, fmt.Sprintf("<q-tab name='%s'  icon='%s' label='%s'/>", tab.Title, tab.Icon, tab.TitleRu))
		} else {
			res = append(res, fmt.Sprintf("<q-tab name='%s'  icon='%s' label='%s'>%[5]s\t%[4]s%[5]s</q-tab>", tab.Title, tab.Icon, tab.TitleRu, tab.HtmlInner, sep))
		}
	}
	return strings.Join(res, sep)
}

// список компонентов для отображения содержимого табов
func (d DocVue) PrintItemTabPanels()  string{
	res := []string{}
	for _, tab := range d.Tabs {
		res = append(res, fmt.Sprintf("<!-- %s       -->\n\t\t\t\t\t\t\t\t<q-tab-panel name='%[2]s'><%[2]s-tab :id='id' %[3]s/></q-tab-panel>", tab.TitleRu, tab.Title, tab.HtmlParams))
	}
	return strings.Join(res, "\n\t\t\t\t\t\t\t\t")
}
