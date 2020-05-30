package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

var project *types.ProjectType

func SetProject(p *types.ProjectType) {
	project = p
}

var funcMap = template.FuncMap{
	"ToUpper":             strings.ToUpper,
	"ToLower":             strings.ToLower,
	"UpperCaseFirst":      utils.UpperCaseFirst,
	"ToLowerCamel":        strcase.ToLowerCamel,
	"PrintVueFldTemplate": PrintVueFldTemplate,
}

func ParseTemplates(p types.ProjectType) map[string]*template.Template {
	// парсинг общих шаблонов
	res := map[string]*template.Template{}
	readFiles := func(prefix, delimLeft, delimRight string, path ...string) {
		tmpls, err := template.New("").Funcs(funcMap).Delims(delimLeft, delimRight).ParseFiles(path...)
		utils.CheckErr(err, "ParseFiles")
		for _, t := range tmpls.Templates() {
			res[prefix+t.Name()] = t
		}
	}
	// project
	path := "../../../pepelazz/projectGenerator/templates/project/"
	readFiles("project_", "{{", "}}", path+"config.toml", path+"main.go", path+"docker-compose.yml", path+"docker-compose.dev.yml", path+"restoreDump.sh", path+"deploy.ps1")

	// webClient
	path = "../../../pepelazz/projectGenerator/templates/webClient/doc/"
	readFiles("webClient_", "[[", "]]", path+"index.vue", path+"item.vue", path+"itemWithTabs.vue", path+"tabInfo.vue", path+"tabTasks.vue")
	// sql
	path = "../../../pepelazz/projectGenerator/templates/sql/"
	readFiles("sql_", "{{", "}}", path+"main.toml")
	path = "../../../pepelazz/projectGenerator/templates/sql/function/"
	readFiles("sql_function_", "{{", "}}", path+"get_by_id.sql", path+"list.sql", path+"update.sql", path+"trigger_before.sql", path+"trigger_after.sql")

	// парсинг шаблонов для конкретного документа
	for i, d := range p.Docs {
		for tName, dt := range d.Templates {
			// возможность расширить функции для шаблона.
			// Если в документе определена FuncMap, то расширяем ее стандартными функциями FuncMap и передаем в шаблон
			fMap := funcMap
			if dt.FuncMap != nil {
				fMap = dt.FuncMap
				for k, v := range funcMap {
					fMap[k] = v
				}
			}
			// извлекаем имя файла шаблона, чтобы использовать его в качестве имени шабона. Иначе могут быть ошибки
			path := strings.Split(dt.Source, "/")
			fName := path[len(path)-1]
			t, err := template.New(fName).Funcs(fMap).Delims("[[", "]]").ParseFiles(dt.Source)
			utils.CheckErr(err, fmt.Sprintf("ParseTemplates doc: %s tmpl: %s parse template error: %s", d.Name, tName, err))
			// сохраняем template в поле структуры
			dt.Tmpl = t
		}
		// дописываем стандартные шаблоны
		baseTmplNames := []string{}
		if d.IsBaseTemapltes.Vue {
			baseTmplNames = append(baseTmplNames, "webClient_item.vue", "webClient_index.vue")
		}
		if d.IsBaseTemapltes.Sql {
			baseTmplNames = append(baseTmplNames, "sql_main.toml", "sql_function_get_by_id.sql", "sql_function_list.sql", "sql_function_update.sql", "sql_function_trigger_before.sql", "sql_function_trigger_after.sql")
		}
		// если документ отмечен свойством рекурсии, то дополнительные шаблоны
		if d.IsRecursion {
			docIsRecursionProccess(p, &d)
		}
		// в случае если указаны табы, то подбираем соответствующие шаблоны
		for _, tab := range d.Vue.Tabs {
			var t *template.Template
			// ищем в списке общих шаблонов
			if t1, ok := res["webClient_"+tab.TmplName]; ok {
				t = t1
			}
			// потом ищем в шаблонах конкретного документа
			if t1, ok := d.Templates["webClient_"+tab.TmplName]; ok {
				t = t1.Tmpl
			}
			if t == nil {
				log.Fatalf("ParseTemplates: Template not found for tab %s webClient_%s", d.Name, tab.TmplName)
			}

			tName := "webClient_tabs_" + tab.Title
			compPath := d.Name
			if len(d.Vue.Path) > 0 {
				compPath = d.Vue.Path // в случае если указан специальный путь к компоненте
			}
			distPath := fmt.Sprintf("%s/webClient/src/app/components/%s/tabs/%s", p.DistPath, compPath, tab.Title)
			d.Templates[tName] = &types.DocTemplate{Tmpl: t, DistPath: distPath, DistFilename: "index.vue"}
		}

		for _, fld := range d.Flds {
			// если в документе есть поле с типо тэг, то создаем sql метод для запроса списка тэгов.
			fldTagProccess(p, &d, &fld)
			// если в документе есть поле с типо jsonList, то создаем специальную компоненту
			fldJsonListProccess(p, &d, &fld)

		}

		for _, tName := range baseTmplNames {
			// если шаблона с таким именем нет, то добавляем стандартный
			if _, ok := d.Templates[tName]; !ok {
				if tName == "sql_function_trigger_before.sql" && !d.Sql.IsBeforeTrigger {
					continue
				}
				if tName == "sql_function_trigger_after.sql" && !d.Sql.IsAfterTrigger {
					continue
				}
				params := map[string]string{}
				if len(d.Vue.Path)> 0 {
					params["doc.Vue.Path"] = d.Vue.Path
				}
				distPath, distFilename := utils.ParseDocTemplateFilename(d.Name, tName, p.DistPath, i, params)
				tmpl := res[tName]
				// возможность переопределить шаблон
				// если указаны табы, то подменяем шаблон item.vue на itemWithTabs.vue
				if len(d.Vue.Tabs) > 0 {
					if strings.HasPrefix(distPath, "../src/webClient/src/app/components") && distFilename == "item.vue" {
						tmpl = res["webClient_itemWithTabs.vue"]
					}
				}
				// игнорируем шаблоны для табов, их добавляем по специальным путям, которые указаны в d.Vue.Tabs (см раздел выше)
				if strings.HasPrefix(tName, "webClient_tab") {
					continue
				}
				d.Templates[tName] = &types.DocTemplate{Tmpl: tmpl, DistPath: distPath, DistFilename: distFilename}
			}
		}

		p.Docs[i] = d
	}

	return res
}

func ExecuteToFile(t *template.Template, d interface{}, path, filename string) error {
	if t == nil {
		log.Fatalf("template is nil for path '%s/%s'\n", path, filename)
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path+"/"+filename, []byte(tpl.String()), 0644)
}

// печать vue темплейтов для
func PrintVueFldTemplate(fld types.FldType) string {
	name := fld.Vue.Name
	if len(name) == 0 {
		name = fld.Name
	}
	nameRu := fld.Vue.NameRu
	if len(nameRu) == 0 {
		nameRu = fld.NameRu
	}
	readonly := fld.Vue.Readonly
	if len(readonly) == 0 {
		readonly = "false"
	}
	fldType := fld.Vue.Type
	if len(fldType) == 0 {
		fldType = fld.Type
		// в случае ref поля
		if fld.Type == types.FldTypeInt && len(fld.Sql.Ref) > 0 {
			fldType = "ref"
		}
	}
	borderStyle := "outlined"
	if fld.Vue.IsBorderless {
		borderStyle = "borderless"
	}
	// если указана функция для композиции, то меняем тип на vueComposition
	if fld.Vue.Composition != nil {
		fldType = types.FldTypeVueComposition
	}
	if fld.Vue.Type == types.FldVueTypeTags {
		fldType = types.FldVueTypeTags
	}
	params := ""
	if len(fld.Vue.Class) > 0 {
		params = params + fmt.Sprintf(" class='%s' ", fld.Vue.ClassPrint())
	}
	if len(fld.Vue.Vif)>0 {
		params = params + fmt.Sprintf(" v-if='%s' ", fld.Vue.Vif)
	}
	switch fldType {
	case types.FldTypeString, types.FldTypeText:
		return fmt.Sprintf(`<q-input %s type='text' v-model="item.%s" label="%s" autogrow :readonly='%s' %s/>`, borderStyle, name, nameRu, readonly, params)
	case types.FldTypeInt, types.FldTypeDouble:
		return fmt.Sprintf(`<q-input %s type='number' v-model="item.%s" label="%s" :readonly='%s' %s/>`, borderStyle, name, nameRu, readonly, params)
	// дата
	case types.FldTypeDate:
		return fmt.Sprintf(`<comp-fld-date %s label="%s" :date-string="$utils.formatPgDate(item.%s)" @update="v=> item.%s = v" :readonly='%s' %s/>`, borderStyle, nameRu, name, name, readonly, params)
	// дата с временем
	case types.FldTypeDatetime:
		return fmt.Sprintf(`<comp-fld-date-time %s label="%s" :date-string="$utils.formatPgDateTime(item.%s)" @update="v=> item.%s = v" :readonly='%s' %s/>`, borderStyle, nameRu, name, name, readonly, params)
	// вариант ссылки на другую таблицу
	case "ref":
		// если map Ext не инициализирован, то создаем его, чтобы не было ошибки при json.Marshal
		if fld.Vue.Ext == nil {
			fld.Vue.Ext = map[string]string{}
		}
		// если специально не определено поле для ajaxSelectTitle, то формируем стандартное [ref_table_name]_title
		ajaxSelectTitle := strings.TrimSuffix(fld.Name, "_id") + "_title"
		if v, ok := fld.Vue.Ext["ajaxSelectTitle"]; ok {
			ajaxSelectTitle = v
		}
		extJsonStr, err := json.Marshal(fld.Vue.Ext)
		utils.CheckErr(err, fmt.Sprintf("json.Marshal(fld.Vue.Ext) fld %s", fld.Name))

		// заполняем название postgres метода для получения списка. По дефолту [ref_table_name]_list
		pgMethod := fld.Sql.Ref + "_list"
		if m, ok := fld.Vue.Ext["pgMethod"]; ok {
			pgMethod = m
		}
		return fmt.Sprintf(`<comp-fld-ref-search %s pgMethod="%s" label="%s" :item='item.%s' :itemId='item.%s' :ext='%s' @update="v=> item.%s = v.id" @clear="item.%s = null" :readonly='%s' %s/>`, borderStyle, pgMethod, nameRu, ajaxSelectTitle, name, extJsonStr, name, name, readonly, params)
	case types.FldVueTypeSelect, types.FldVueTypeMultipleSelect:
		options, err := json.Marshal(fld.Vue.Options)
		utils.CheckErr(err, fmt.Sprintf("'%s' json.Marshal(fld.Vue.Options)", fld.Name))
		multiple := ""
		if fldType == types.FldVueTypeMultipleSelect {
			multiple = "multiple"
		}
		isClearable := ""
		for key, _ := range fld.Vue.Ext {
			if key == "isClearable" {
				isClearable = "clearable"
			}
		}
		return fmt.Sprintf(`<q-select %s label="%s" v-model='item.%s' :options='%s' %s %s :readonly='%s' %s/>`, borderStyle, nameRu, name, options, multiple, isClearable, readonly, params)
	case types.FldTypeVueComposition:
		if fld.Vue.Composition == nil {
			log.Fatal(fmt.Sprintf("fld have type '%s', but fld.Vue.Composition function is nil", types.FldTypeVueComposition))
		}
		// возможен вариант что функция рендеринга поля шаблона вызываается до того как сам документ был инициализирован и соответственно была заполнена ссылка на него в поле fld.Doc
		// в таком случае в функуию передаем пустой документ. Если функция не использует ссылку на документ, то все ок. Но если в функции идет обращение к инфе о документе, то функция отработает некорректно.
		// возможное решение, чтобы вызов функции в шаблоне происходил уже после инициализации документа
		linkOnDoc := types.DocType{}
		if fld.Doc != nil {
			linkOnDoc = *fld.Doc
		}
		return fld.Vue.Composition(*project, linkOnDoc)
	case types.FldVueTypeTags:
		if fld.Vue.Ext["onlyExistTags"] == "true" {
			// вариант когда нельзя создавать новые тэги, только выбирать из существующих
			return fmt.Sprintf("<q-select %s label='%s' v-model='item.%s' use-chips multiple @filter='%[2]sFilterFn' :options='%[2]sFilterOptions' :readonly='%s'/>", borderStyle, nameRu, name, readonly, params)
		} else {
			// вариант когда можно создавать новые тэги
			return fmt.Sprintf("<q-select %s label='%s' v-model='item.%s' use-input use-chips multiple input-debounce='0' @new-value='%[2]sCreateValue' @filter='%[2]sFilterFn' :options='%[2]sFilterOptions' :readonly='%s'/>", borderStyle, nameRu, name, readonly, params)
		}
	case types.FldVueTypeCheckbox:
		return fmt.Sprintf("<q-checkbox label='%s' v-model='item.%s' :disabled='%s' :false-value='null' indeterminate-value='some' %s/>", nameRu, name, readonly, params)
	case types.FldVueTypeRadio:
		options := ""
		for _, v := range fld.Vue.Options {
			options = fmt.Sprintf("%s <q-radio size='xs' dense v-model='item.%s' val='%s' label='%s' :disable ='%s'/>\n", options, name, v.Value, v.Label, readonly)
		}
		return fmt.Sprintf(`<div class="row q-col-gutter-md q-pb-xs">
	<div class="col-4">%s</div>
	<div class="col-8 q-gutter-sm ">
	%s
	</div>
	</div>
`, nameRu, options)
	case types.FldVueTypeDadataAddress:
		return fmt.Sprintf(`<comp-fld-address %s label="%s" :fld='item.%s' @update="v=> item.%s = v" :readonly='%s' %s/>`, borderStyle, nameRu, name, name, readonly, params)
	case types.FldVueTypeJsonList:
		return fmt.Sprintf("<comp-fld-json-list-%s label='%s' :item='item' :fld='item.%s' @update='item.%s = $event' :readonly='%s' %s/>", name, nameRu, name, name, readonly, params)
	case types.FldVueTypeFiles:
		extStr := ""
		for k, v := range fld.Vue.Ext {
			extStr = extStr + fmt.Sprintf(", %s: \"%s\"", k, v)
		}
		return fmt.Sprintf("<comp-fld-files label='%s' :fld='item.%s' :ext = '{tableName: \"%s\", tableId: this.id%s}' :readonly='%s' %s/>", nameRu, name, fld.Doc.Name, extStr, readonly, params)
	case types.FldVueTypeImg:
		extStr := ""
		for k, v := range fld.Vue.Ext {
			extStr = extStr + fmt.Sprintf(", %s: \"%s\"", k, v)
		}
		return fmt.Sprintf("<comp-fld-img label='%s' :fld='item.%s' :ext = '{tableName: \"%s\", tableId: this.id, fldName: \"%s\"%s}' :readonly='%s' %s/>", nameRu, name, fld.Doc.Name, name, extStr, readonly, params)
	case types.FldVueTypeImgList:
		extStr := ""
		for k, v := range fld.Vue.Ext {
			extStr = extStr + fmt.Sprintf(", %s: \"%s\"", k, v)
		}
		return fmt.Sprintf("<comp-fld-img-list label='%s' :fld='item.%s' :ext = '{tableName: \"%s\", tableId: this.id, fldName: \"%s\"%s}' :readonly='%s' %s/>", nameRu, name, fld.Doc.Name, name, extStr, readonly, params)
	default:
		return fmt.Sprintf("not found vueFldTemplate for type `%s`", fldType)
	}
}
