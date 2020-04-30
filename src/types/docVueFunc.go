package types

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/src/utils"
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
	res := []string{}
	// ссылки на миксины
	if d.Vue.Mixins != nil {
		if arr, ok := d.Vue.Mixins[tmplName]; ok {
			for _, s := range arr {
				res = append(res, fmt.Sprintf("\timport %s from './mixins/%s'", s, s))
			}
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
