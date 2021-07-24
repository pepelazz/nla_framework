package templates

import (
	"fmt"
	"github.com/pepelazz/nla_framework/types"
	"github.com/pepelazz/nla_framework/utils"
	"text/template"
)

// если в документе есть поле с типо jsonList, то создаем специальную компоненту
func fldJsonListProccess(p types.ProjectType, d *types.DocType, fld *types.FldType)  {
	if fld.Vue.Type == types.FldVueTypeJsonList {
		docName := d.Name
		fldName := fld.Name
		JsonList := fld.Vue.JsonList
		funcMap := map[string]interface{}{
			"GetDoc": func() string {return docName},
			"PrintVueFldTemplate": PrintVueFldTemplate,
			"GetJsonList": func() types.FldVueJsonList {return JsonList},
		}
		sourcePath := fmt.Sprintf("../../../pepelazz/nla_framework/templates/webClient/quasar_%v/doc/comp/fldJsonList.vue", p.GetQuasarVersion())
		// проверяем возможность того, что путь к шаблону был переопределен внутри документа
		if d.TemplatePathOverride != nil {
			if tmpl, ok := d.TemplatePathOverride["fldJsonList.vue"]; ok {
				if len(tmpl.Source)> 0 {
					sourcePath = tmpl.Source
				}
			}
		}
		t, err := template.New("fldJsonList.vue").Funcs(funcMap).Delims("[[", "]]").ParseFiles(sourcePath)
		utils.CheckErr(err, "fldJsonList.vue")
		dPath := d.Name
		if len(d.Vue.Path) > 0 {
			dPath = d.Vue.Path
		}
		distPath := fmt.Sprintf("%s/webClient/src/app/components/%s/comp", p.DistPath, dPath)
		// в случае табов изменяем path
		if len(d.Vue.Tabs) > 0{
			distPath = fmt.Sprintf("%s/webClient/src/app/components/%s/tabs/info/comp", p.DistPath, d.Name)
		}
		// проверяем возможность того, что путь для шаблону был переопределен внутри документа
		if d.TemplatePathOverride != nil {
			if tmpl, ok := d.TemplatePathOverride["fldJsonList.vue"]; ok {
				if len(tmpl.Dist)> 0 {
					distPath = tmpl.Dist
				}
			}
		}
		d.Templates["webClient_comp_" + fld.Name +"_fldJsonList.vue"] = &types.DocTemplate{Tmpl: t, DistPath: distPath, DistFilename: "compFldJsonList" + utils.UpperCaseFirst(fldName) + ".vue"}

		// добавляем в список компонентов
		importName := "compFldJsonList" + utils.UpperCaseFirst(fldName)
		importAddress := "./comp/" + importName + ".vue"
		if d.Vue.Components == nil {
			d.Vue.Components = map[string]map[string]string{}
		}
		if d.Vue.Components["docItem"] == nil {
			d.Vue.Components["docItem"] = map[string]string{}
		}
		d.Vue.Components["docItem"][importName] = importAddress
	}
}
