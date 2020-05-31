package templates

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"github.com/serenize/snaker"
	"log"
	"text/template"
)

func docIsIntegrationProccess(p types.ProjectType, d *types.DocType)  {
	if len(d.Integrations.Bitrix.Name) > 0 {
		docIsIntegrationBitrixProccess(p, d)
	}
}

func docIsIntegrationBitrixProccess(p types.ProjectType, d *types.DocType) {
	sourcePath := "../../../pepelazz/projectGenerator/templates/integrations/bitrix/bitrixDoc.go"
	// проверяем возможность того, что путь к шаблону был переопределен внутри документа
	if d.TemplatePathOverride != nil {
		if tmpl, ok := d.TemplatePathOverride["bitrixDoc.go"]; ok {
			if len(tmpl.Source)> 0 {
				sourcePath = tmpl.Source
			}
		}
	}
	docName := d.Name
	localFuncMap := template.FuncMap{
		"LocalProjectPath": func() string{ return p.Config.LocalProjectPath},
		"DocNameCamel": func() string{ return snaker.SnakeToCamel(docName)},
		"IsBtxFld": func(fld types.FldType) bool {
			return len(getBtxFld(*d, fld).Name) > 0
		},
		"GetBtxFldName": func(fld types.FldType) string {
			return getBtxFld(*d, fld).Name
		},
		"GetBtxFldType": func(fld types.FldType) string {
			t := getBtxFld(*d, fld).Type
			if len(t) == 0 {
				return "interface{}"
			}
			return t
		},
		"CastToGoType": func(fld types.FldType) string {
			fName := strcase.ToCamel(fld.Name)
			// если в описании поля указан способ приведения к типу, то используем его
			if len(getBtxFld(*d, fld).CastToGoType) > 0 {
				return getBtxFld(*d, fld).CastToGoType
			}
			switch fld.Type {
			case types.FldTypeText, types.FldTypeString:
				return fmt.Sprintf("res.%[1]s = cast.ToString(btxDoc.%[1]s)", fName)
			case types.FldTypeInt:
				return fmt.Sprintf("res.%[1]s = cast.ToInt(btxDoc.%[1]s)", fName)
			case types.FldTypeInt64:
				return fmt.Sprintf("res.%[1]s = cast.ToInt64(btxDoc.%[1]s)", fName)
			case types.FldTypeIntArray:
				return fmt.Sprintf(`res.%[1]s = []int{}
				intSlice%[1]s, err := cast.ToIntSliceE(btxDoc.%[1]s)
				if err == nil {
					res.%[1]s = intSlice%[1]s
				}`, fName)
			case types.FldTypeTextArray:
				return fmt.Sprintf(`res.%[1]s = []string{}
				txtSlice%[1]s, err := cast.ToStringSliceE(btxDoc.%[1]s)
				if err == nil {
					res.%[1]s = txtSlice%[1]s
				}`, fName)
			}
			return "`!!! CastToGoType not found for type: " + fld.Type + " fld: " + fld.Name + "`"
		},
	}
	for k, v := range funcMap {
		localFuncMap[k] = v
	}
	t, err := template.New("bitrixDoc.go").Funcs(localFuncMap).Delims("[[", "]]").ParseFiles(sourcePath)
	utils.CheckErr(err, "bitrixDoc.go")
	distPath := fmt.Sprintf("%s/bitrix", p.DistPath)
	d.Templates["webClient_comp_bitrixDoc.go"] = &types.DocTemplate{Tmpl: t, DistPath: distPath, DistFilename: snaker.SnakeToCamelLower(d.Name) + ".go"}
}

func getBtxFld(d types.DocType, fld types.FldType) types.BitrixFld  {
	if btxFldInt, ok := fld.IntegrationData["bitrix"]; ok {
		if btxFld, ok := btxFldInt.(types.BitrixFld); ok {
			return btxFld
		} else {
			log.Fatalf("docIsIntegrationBitrixProccess doc: '%s' fld: '%s' not BitrixFld", d.Name, fld.Name)
		}
	}
	return types.BitrixFld{}
}
