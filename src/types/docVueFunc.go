package types

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