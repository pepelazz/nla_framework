package types

import (
	"fmt"
	"github.com/pepelazz/nla_framework/utils"
	"log"
	"path"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

// создание поля title
func GetFldTitle(params ...string) (fld FldType) {
	var classStr string
	nameRu := "название"
	if len(params)>0 {
		classStr= params[0]
	}
	for _, v := range params {
		if strings.HasPrefix(v, "name_ru:") {
			nameRu = strings.TrimSpace(strings.TrimPrefix(v, "name_ru:"))
		}
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType {Name: "title", NameRu: nameRu, Type: FldTypeString, Sql: FldSql{IsRequired: true, IsUniq: true, IsSearch:true, Size:150}, Vue: FldVue{RowCol: [][]int{{1, 1}}, Class: []string{classStr}}}
	return
}

// создание поля title, которое заполняется тригером
func GetFldTitleComputed(triggerSqlString string, params ...string) (fld FldType) {
	var classStr string
	if len(params)>0 {
		classStr= params[0]
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType {Name: "title", NameRu: "название", Type: FldTypeString, Sql: FldSql{IsSearch:true, FillValueInBeforeTrigger: triggerSqlString}, Vue: FldVue{RowCol: [][]int{{1, 1}}, Class: []string{classStr}}}
	return
}

func GetFldDouble(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	if len(params)>0 {
		classStr= params[0]
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeDouble, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля String
func GetFldString(name, nameRu string, size int, rowCol [][]int, params ...string) (fld FldType) {
	classArr := []string{getDefaultClassStr("")}
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classArr = []string{getDefaultClassStr(v)}
		} else {
			if strings.HasPrefix(v, "col-") {
				classArr = append(classArr, v)
			}
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Vue:FldVue{RowCol: rowCol, Class: classArr, Readonly:readonly}}
	if size > 0 {
		fld.Sql.Size = size
	}
	return
}

// создание простого Date
func GetFldDate(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	if len(params)>0 {
		classStr= params[0]
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeDate, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого DateTime
func GetFldDateTime(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	if len(params)>0 {
		classStr= params[0]
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeDatetime, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля Int
func GetFldInt(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	if len(params)>0 {
		classStr= params[0]
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeInt, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля Int64
func GetFldInt64(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	if len(params)>0 {
		classStr= params[0]
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeInt64, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание поля UUID
func GetFldUuid(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	if len(params)>0 {
		classStr= params[0]
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeUuid, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля Checkbox
func GetFldCheckbox(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeBool, Vue:FldVue{RowCol: rowCol, Type:FldVueTypeCheckbox, Class: []string{classStr}, Readonly:readonly}}
	return
}

// создание простого поля Radio
func GetFldRadioString(name, nameRu string, rowCol [][]int, options []FldVueOptionsItem, params ...string) (fld FldType) {
	var classStr string
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Sql:FldSql{Size:50}, Vue:FldVue{RowCol: rowCol, Type:FldVueTypeRadio, Options: options, Class: []string{classStr}, Readonly:readonly}}
	return
}

// создание простого поля Ref
// - isShowLink
// - isAddNew
// - isClearable
func GetFldRef(name, nameRu, refTable string, rowCol [][]int, params ...string) (fld FldType) {
	classArr := []string{getDefaultClassStr("")}
	for i, v := range params {
		if strings.HasPrefix(v, "col-") {
			if i == 0 {
				classArr = []string{getDefaultClassStr(v)}
			} else {
				classArr = append(classArr, v)
			}
		}

	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeInt,  Sql: FldSql{Ref: refTable, IsSearch:true}, Vue:FldVue{RowCol: rowCol, Ext: map[string]string{}, Class: classArr}}
	for _, v := range params {
		// добавляем аватарку с ссылкой на выбранный документ
		if v == "isShowLink" {
			// проставляем значение pathUrl и avatar на последнем шаге, после инициализации всех документов  в методе FillVueFlds
			fld.Vue.Ext["pathUrl"] = ""
			fld.Vue.Ext["avatar"] = ""
		}
		// добавляем возможность создание новой записи
		if v == "isAddNew" {
			// проставляем значение addNewUrl на последнем шаге, после инициализации всех документов  в методе FillVueFlds
			fld.Vue.Ext["addNewUrl"] = ""
		}
		if v == "isClearable" {
			fld.Vue.Ext["isClearable"] = "true"
		}
		if strings.HasPrefix(v, "ext:") {
			// записываем как есть, потом преобразуем при записи во vue
			fld.Vue.Ext["rawJsonExt"] = strings.TrimSpace(strings.TrimPrefix(v, "ext:"))
		}
	}
	return
}

// создание поля phone
func GetFldPhone(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Sql: FldSql{Size: 30}, Vue:FldVue{RowCol: rowCol, Type:FldVueTypePhone, Class: []string{classStr}, Readonly:readonly}}
	return
}

// создание поля email
func GetFldEmail(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	var classStr string
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Sql: FldSql{Size: 100}, Vue:FldVue{RowCol: rowCol, Type:FldVueTypeEmail, Class: []string{classStr}, Readonly:readonly}}
	return
}

// поле с кастомной композицией
func GetFldJsonbComposition(name, nameRu string, rowCol [][]int, classStr, compName string, params ...string) (fld FldType) {
	isOptionsFld := ""
	sqlType := FldTypeJsonb
	classArr := []string{getDefaultClassStr(classStr)}
	for i, v := range params {
		// IsOptionFld передаем отдельным параемтром, потому что SetIsOptionFld() срабатывает уже после того как строка с компонентой сформмирована
		if v == "IsOptionFld" {
			isOptionsFld = "options."
		}
		if strings.HasPrefix(v, "col-") {
			classArr = append(classArr, v)
		}
		if strings.HasPrefix(v, "sqlType:") {
			sqlType = strings.TrimSpace(strings.TrimPrefix(v, "sqlType:"))
			// убираем из params, потому что они дальше печатаются во vue компоненту
			params[i] = ""
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type: sqlType,  Vue:FldVue{RowCol: rowCol, Class: classArr, Composition: func(p ProjectType, d DocType, fld FldType) string {
		return fmt.Sprintf("<%[1]s :fld='item.%[5]s%[2]s' :item='item' @update='item.%[5]s%[2]s = $event' label='%[3]s' %[4]s/>", compName, name, nameRu, strings.Join(params, " "), isOptionsFld)
	}}}
	return
}

// поле с кастомной композицией
func GetFldJsonbCompositionWithoutFld(rowCol [][]int, classStr, compName string, params ...string) (fld FldType) {
	classArr := []string{getDefaultClassStr(classStr)}
	for i, v := range params {
		if strings.HasPrefix(v, "col-") {
			classArr = append(classArr, v)
			// стираем класс, чтобы не попал в атрибуты vue компоненты
			params[i] = ""
		}
	}
	fld = FldType{Type:FldTypeVueComposition,  Vue:FldVue{RowCol: rowCol, Class: classArr, Composition: func(p ProjectType, d DocType, fld FldType) string {
		return fmt.Sprintf("<%[1]s :item='item' %[2]s/>", compName, strings.Join(params, " "))
	}}}
	return
}

// простое html поле
func GetFldSimpleHtml(rowCol [][]int, classStr, htmlStr string) (fld FldType) {
	classStr = getDefaultClassStr(classStr)
	fld = FldType{Type:FldTypeVueComposition,  Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType, fld FldType) string {
		return htmlStr
	}}}
	return
}

// создание простого поля Select с типом string
func GetFldSelectString(name, nameRu string, size int, rowCol [][]int, options []FldVueOptionsItem, params ...string) (fld FldType) {
	readonly := "false"
	classStr := getDefaultClassStr("")
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeSelect, Ext: map[string]string{}, Class: []string{classStr}, Readonly:readonly, Options:options}}
	for i, v := range params {
		if i == 0 {
			fld.Vue.Class = []string{getDefaultClassStr(v)}
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				fld.Vue.Readonly = "true"
			}
			if strings.HasPrefix(v, "col-") {
				fld.Vue.Class = append(fld.Vue.Class, v)
			}
			if v == "isClearable" {
				fld.Vue.Ext["isClearable"] = "true"
			}
		}
	}
	if size > 0 {
		fld.Sql.Size = size
	}
	return
}

// создание простого поля MultipleSelect с типом string
func GetFldSelectMultiple(name, nameRu string, rowCol [][]int, options []FldVueOptionsItem, params ...string) (fld FldType) {
	readonly := "false"
	classStr := getDefaultClassStr("")
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeTextArray, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeMultipleSelect, Ext: map[string]string{}, Class: []string{classStr}, Readonly:readonly, Options:options}}
	for i, v := range params {
		if i == 0 {
			fld.Vue.Class = []string{getDefaultClassStr(v)}
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				fld.Vue.Readonly = "true"
			}
			if strings.HasPrefix(v, "col-") {
				fld.Vue.Class = append(fld.Vue.Class, v)
			}
			if v == "isClearable" {
				fld.Vue.Ext["isClearable"] = "true"
			}
		}
	}

	return
}

// создание простого поля Int
func GetFldTag(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classArr := []string{getDefaultClassStr("")}
	onlyExistTags := "false" // флаг для UI контрола, чтобы можно было только выбирать из существующих тэгов и нельзя было создавать новые
	for i, v := range params {
		if i == 0 {
			classArr = []string{getDefaultClassStr(v)}
		} else {
			if strings.HasPrefix(v, "col-") {
				fld.Vue.Class = append(fld.Vue.Class, v)
			}
			if strings.HasPrefix(v, "only_exist_tags") {
				onlyExistTags="true"
			}
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeTextArray, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeTags,  Class: classArr, Ext: map[string]string{"onlyExistTags": onlyExistTags}}}
	return
}

// создание поля-виджета со связями многие-к-многим
func GetFldLinkListWidget(linkTable string, rowCol [][]int, classStr string, opts map[string]interface{}) (fld FldType) {
	classStr = getDefaultClassStr(classStr)
	return FldType{Type: FldTypeVueComposition,  Vue: FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType, fld FldType) string {
		return GetVueCompLinkListWidget(p, d, linkTable, opts)
	}}}
}

// функция конвертации списка имен файлов с шаблонами в  map[string]*DocTemplate
func GetCustomTemplates(p ...string) map[string]*DocTemplate  {
	res := map[string]*DocTemplate{}
	for _, name := range p {
		res[name] = &DocTemplate{}
	}
	return res
}

// создание поля адрес с возможностью поиска через dadata
func GetFldDadataAddress(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := getDefaultClassStr("")
	if len(params)>0 {
		classStr= getDefaultClassStr(params[0])
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeDadataAddress, Class: []string{classStr}}}
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				fld.Vue.Readonly = "true"
			}
			if v == "isClearable" {
				if fld.Vue.Ext == nil {
					fld.Vue.Ext = map[string]string{}
				}
				fld.Vue.Ext["isClearable"] = "true"
			}
		}
	}
	return
}

// создание поля json c редактируемым массивом элементов
func GetFldJsonList(name, nameRu string, rowCol [][]int, listParams FldVueJsonList, params ...string) (fld FldType) {
	classStr := getDefaultClassStr("")
	if len(params)>0 {
		classStr= getDefaultClassStr(params[0])
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeJsonList, JsonList: listParams, Class: []string{classStr}}}
	return
}

// создание поля для загрузки файлов
func GetFldFiles(name, nameRu string, rowCol [][]int, fileParams FldVueFilesParams, params ...string) (fld FldType) {
	classStr := getDefaultClassStr("")
	if len(params)>0 {
		classStr= getDefaultClassStr(params[0])
	}
	// заполняем параметры для ограничений по загрузке файлов
	ext := map[string]string{}
	if len(fileParams.Accept)>0{
		ext["accept"] = fileParams.Accept
	}
	if fileParams.MaxFileSize>0{
		ext["maxFileSize"] = strconv.FormatInt(fileParams.MaxFileSize, 10)
	}

	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeFiles, Ext: ext, Class: []string{classStr}}}
	return
}

// создание поля для загрузки списка изображений
func GetFldImgList(name, nameRu string, rowCol [][]int, fileParams FldVueImgParams, params ...string) (fld FldType) {
	classStr := getDefaultClassStr("")
	if len(params)>0 {
		classStr= getDefaultClassStr(params[0])
	}
	// заполняем параметры для ограничений
	ext := map[string]string{}
	if len(fileParams.Accept) > 0{
		ext["accept"] = fileParams.Accept
	}
	if fileParams.MaxFileSize > 0{
		ext["maxFileSize"] = strconv.FormatInt(fileParams.MaxFileSize, 10)
	}
	if fileParams.CanAddUrls {
		ext["canAddUrls"] = "true"
	}
	if len(fileParams.Crop) > 0 {
		// проверка что crop имеет формат 300x400
		arr := strings.Split(fileParams.Crop, "x")
		if len(arr) != 2 {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. You write this: %s", name, fileParams.Crop)
		}
		if _, err := strconv.Atoi(arr[0]); err != nil {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. %s not number", name, arr[0])
		}
		if _, err := strconv.Atoi(arr[1]); err != nil {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. %s not number", name, arr[1])
		}
		ext["crop"] = fileParams.Crop
	}
	if fileParams.Width > 0 {
		ext["width"] = strconv.Itoa(fileParams.Width)
	}

	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeImgList, Ext: ext, Class: []string{classStr}}}
	return
}

// создание поля для загрузки одного
func GetFldImg(name, nameRu string, rowCol [][]int, fileParams FldVueImgParams, params ...string) (fld FldType) {
	classStr := getDefaultClassStr("")
	if len(params)>0 {
		classStr= getDefaultClassStr(params[0])
	}
	// заполняем параметры для ограничений
	ext := map[string]string{}
	if len(fileParams.Accept) > 0{
		ext["accept"] = fileParams.Accept
	}
	if fileParams.MaxFileSize > 0{
		ext["maxFileSize"] = strconv.FormatInt(fileParams.MaxFileSize, 10)
	}
	if fileParams.CanAddUrls {
		ext["canAddUrls"] = "true"
	}
	if len(fileParams.Crop) > 0 {
		// проверка что crop имеет формат 300x400
		arr := strings.Split(fileParams.Crop, "x")
		if len(arr) != 2 {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. You write this: %s", name, fileParams.Crop)
		}
		if _, err := strconv.Atoi(arr[0]); err != nil {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. %s not number", name, arr[0])
		}
		if _, err := strconv.Atoi(arr[1]); err != nil {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. %s not number", name, arr[1])
		}
		ext["crop"] = fileParams.Crop
	}
	if fileParams.Width > 0 {
		ext["width"] = strconv.Itoa(fileParams.Width)
	}

	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Sql:FldSql{Size:500}, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeImg, Ext: ext, Class: []string{classStr}}}
	return
}

// добавление для таба функциональности счетчика
// добавляется миксин, чтобы в основном табе при открытии загружался список, длина которого и является счетчиком
func (vt VueTab) AddCounter(d *DocType, tabName, pgMethod, pgParams string)  VueTab {
	tabName = utils.UpperCaseFirst(tabName)
	if d.Vue.Mixins == nil {
		d.Vue.Mixins = map[string][]VueMixin{}
	}
	if d.Vue.Mixins["docItemWithTabs"] == nil {
		d.Vue.Mixins["docItemWithTabs"] = []VueMixin{}
	}
	d.Vue.Mixins["docItemWithTabs"] = append(d.Vue.Mixins["docItemWithTabs"], VueMixin{"tabCounter"+tabName, "./mixins/tabCounter"+tabName+".js"})
	sourcePath := fmt.Sprintf("%s/templates/webClient/quasar_%v/doc/mixins/tabCounter.js", getRootDirPath(), d.GetProject().GetQuasarVersion())
	funcMap := template.FuncMap{
		"VarName": func() string {return "tabCounter"+tabName},
		"PgMethod": func() string {return pgMethod},
		"PgParams": func() string {return pgParams},
	}

	docRouteName := d.Name
	if len(d.Vue.Path) > 0 {
		docRouteName = d.Vue.Path
	}
	distPath := fmt.Sprintf("../src/webClient/src/app/components/%s/mixins", docRouteName)
	d.Templates["webClient_mixin_tabCounter"+tabName+".js"] = &DocTemplate{
		Source: sourcePath,
		DistPath: distPath,
		FuncMap: funcMap,
		DistFilename: "tabCounter"+tabName+".js",
	}
	// добавляем параметры в html разметку таба
	vt.HtmlParams = vt.HtmlParams + " @updateCount='v => tabCounter"+tabName+" = v'"
	vt.HtmlInner = vt.HtmlInner + " <q-badge v-if='tabCounter"+tabName+">0' color='red' floating>{{tabCounter"+tabName+"}}</q-badge>"
	return vt
}

func getDefaultClassStr(v string) string  {
	if len(v) == 0 || v == "col-4" {
		return "col-md-4 col-sm-6 col-xs-12"
	}
	if v == "col-1" {
		return "col-md-1 col-sm-2 col-xs-6"
	}
	if v == "col-2" {
		return "col-md-2 col-sm-3 col-xs-6"
	}
	if v == "col-8" {
		return "col-md-8 col-sm-12 col-xs-12"
	}
	return v
}

func getRootDirPath() string  {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("ParseTemplates runtime.Caller: No caller information")
	}
	return strings.TrimSuffix(path.Dir(filename), "/types")
}
