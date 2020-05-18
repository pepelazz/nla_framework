package templates

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"github.com/serenize/snaker"
	"text/template"
)

// если в документе есть поле с типо тэг, то создаем sql метод для запроса списка тэгов.
// При формировании шаблона передаем в него функцию GetFld для получения названия поля, для которого создана функция
func fldTagProccess(p types.ProjectType, d *types.DocType, fld *types.FldType)  {
	if fld.Vue.Type == types.FldVueTypeTags {
		docName := d.Name
		fldName := fld.Name
		methodName := d.Name + "_" + fld.Name +"_list"
		funcMap := map[string]interface{}{
			"GetDoc": func() string {return docName},
			"GetFld": func() string {return fldName},
		}
		sourcePath := "../../../pepelazz/projectGenerator/templates/sql/function/tag_list.sql"
		// проверяем возможность того, что путь к шаблону был переопределен внутри документа
		if d.TemplatePathOverride != nil {
			if tmpl, ok := d.TemplatePathOverride["tag_list.sql"]; ok {
				if len(tmpl.Source)> 0 {
					sourcePath = tmpl.Source
				}
			}
		}
		t, err := template.New("tag_list.sql").Funcs(funcMap).ParseFiles(sourcePath)
		utils.CheckErr(err, "tag_list.sql")
		distPath := fmt.Sprintf("%s/sql/template/function/_%s", p.DistPath, snaker.SnakeToCamel(d.Name))
		d.Templates["sql_function_" + fld.Name +"_tag_list.sql"] = &types.DocTemplate{Tmpl: t, DistPath: distPath, DistFilename: methodName + ".sql"}
		// добавляем в список sql методов
		if d.Sql.Methods == nil {
			d.Sql.Methods = map[string]*types.DocSqlMethod{}
		}
		d.Sql.Methods[methodName] = &types.DocSqlMethod{Name: methodName}
		// читаем шаблон и генерим файл с mixin
		t, err = template.New("mixinTag.js").Funcs(funcMap).Delims("[[", "]]").ParseFiles("../../../pepelazz/projectGenerator/templates/webClient/doc/mixinTag.js")
		utils.CheckErr(err, "mixinTag.js")
		distPath = fmt.Sprintf("%s/webClient/src/app/components/%s/mixins", p.DistPath, d.Name)
		// в случае табов изменяем path
		if len(d.Vue.Tabs) > 0{
			distPath = fmt.Sprintf("%s/webClient/src/app/components/%s/tabs/info/mixins", p.DistPath, d.Name)
		}
		// проверяем возможность того, что путь для шаблону был переопределен внутри документа
		if d.TemplatePathOverride != nil {
			if tmpl, ok := d.TemplatePathOverride["tag_list.js"]; ok {
				if len(tmpl.Dist)> 0 {
					distPath = tmpl.Dist
				}
			}
		}
		d.Templates["webClient_mixin_" + fld.Name +"_tag_list.js"] = &types.DocTemplate{Tmpl: t, DistPath: distPath, DistFilename: fld.Name +"_tag_list.js"}
		// добавляем в список миксинов
		if d.Vue.Mixins == nil {
			d.Vue.Mixins = map[string][]types.VueMixin{}
		}
		if d.Vue.Mixins["docItem"] == nil {
			d.Vue.Mixins["docItem"] = []types.VueMixin{}
		}
		d.Vue.Mixins["docItem"] = append(d.Vue.Mixins["docItem"], types.VueMixin{fld.Name +"_tag_list", "./mixins/" + fld.Name +"_tag_list"})
	}
}
