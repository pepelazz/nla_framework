package types

import "github.com/pepelazz/projectGenerator/src/utils"

func (d DocType) PrintListRowLabel() string  {
	res := `
        <q-item-section>
          <q-item-label lines="1">{{item.title}}</q-item-label>
        </q-item-section>
	`
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
