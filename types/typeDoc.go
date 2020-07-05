package types

import (
	"github.com/pepelazz/projectGenerator/utils"
	"github.com/serenize/snaker"
	"log"
	"text/template"
)

const (
	DocTypeLinkTable = "linkTable"
)

type (
	DocType struct {
		Name                 string
		NameRu               string
		Type                 string
		Flds                 []FldType
		Vue                  DocVue
		Sql                  DocSql
		Templates            map[string]*DocTemplate
		TemplatePathOverride map[string]TmplPathOverride // map для переопределения источника шаблона по его названию
		IsBaseTemapltes      DocIsBaseTemapltes          // флаг что генерируем стандартные шаблоны для документа
		PathPrefix           string                      // префикс,если папка, в которой лежит папка с описанием документа находится не на одном уровне с main.go. Например 'docs', если docs/client/...
		IsTaskAllowed        bool                        // признак, что к таблице можно прикреплять задачи
		StateMachine         *DocSm
		IsRecursion          bool // признак, что документ имеет рекурсию. Есть parent_id - ссылка на самого себя
		Integrations         DocIntegrations
	}

	TmplPathOverride struct {
		Source string
		Dist   string
	}

	DocVue struct {
		RouteName      string
		Routes         [][]string // можно указать роуты, тогда они не формируются автоматически
		Path           string     // путь к папке с компонентами, если отличается от стандартного. Например client/deal... Используется для вложенных путей
		MenuIcon       string
		BreadcrumbIcon string
		Roles          []string
		Grid           []VueGridDiv
		Mixins         map[string][]VueMixin        // название файла - название миксина. Для прописывания импорта
		Components     map[string]map[string]string // название файла - название миксина: путь для импорта. Для прописывания импорта
		Methods        map[string]map[string]string // название файла - название метода - текст функции
		TmplFuncs      map[string]func(DocType) string
		I18n           map[string]string
		GloablI18n     map[string]map[string]string // для вынесение справочника в utils.js, чтобы потом можно было вызывать $util.i18n_<название функции>
		Tabs           []VueTab
		Hooks          DocVueHooks // куски vue кода
		Readonly       string
	}

	VueTab struct {
		Title      string
		TitleRu    string
		TmplName   string
		Icon       string
		HtmlParams string
		HtmlInner  string
	}

	VueMixin struct {
		Title  string
		Import string
	}

	// специальное представление для сетки
	VueGridDiv struct {
		Class string
		Grid  []VueGridDiv
		Fld   FldType
	}

	DocTemplate struct {
		Source       string
		DistPath     string
		DistFilename string
		Tmpl         *template.Template
		FuncMap      template.FuncMap // возможность добавлять для конкретного шаблона свои функции, которые затем можно использовать внутри шаблона вместе со стандартными функциями
	}

	DocSql struct {
		Methods         map[string]*DocSqlMethod
		IsUniqLink      bool        // флаг, что таблица является связью двух таблиц и связь между ними уникальная
		IsBeforeTrigger bool        // флаг что добавляем before триггер
		IsAfterTrigger  bool        // флаг что добавляем after триггер
		IsSearchText    bool        // флаг что добавляем поле search_text
		Indexes         []string    // индексы
		Hooks           DocSqlHooks // куски sql кода
	}

	DocIsBaseTemapltes struct {
		Vue bool
		Sql bool
	}

	DocSqlMethod struct {
		Name   string
		Roles  []string
		Params map[string]string
		Tmpl   DocSqlMethodTmpl
	}

	DocSqlMethodTmpl struct {
		Source  string
		Dist    string
		FuncMap template.FuncMap
	}

	DocSqlHooks struct {
		DeclareVars         map[string]string
		BeforeInsertUpdate  []string
		BeforeInsert        []string
		AfterInsertUpdate   []string
		BeforeTriggerBefore []string
	}

	DocIntegrations struct {
		Bitrix DocIntegrationsBitrix
		Odata  DocIntegrationsOdata
	}

	DocIntegrationsBitrix struct {
		Name        string
		UrlName     string // часть имени запроса. Например crm.company.list.json
		IsDebugMode bool   // показываем открытый get метод для тестирования импорта
		Result      struct {
			StructDesc string // описание вложенной структуры для маппинга json
			PathStr    string // путь до массива с данными. Например, Result.Tasks
		}
		UrlQuery       string
		IsNoPagination bool // признак, что все данные получаются за один запрос
	}

	DocIntegrationsOdata struct {
		Name        string
		Url         string // часть имени запроса. Например crm.company.list.json
		IsDebugMode bool   // показываем открытый get метод для тестирования импорта
		//Result struct {
		//	StructDesc string // описание вложенной структуры для маппинга json
		//	PathStr string // путь до массива с данными. Например, Result.Tasks
		//}
		//UrlQuery string
	}

	DocVueHooks struct {
		ItemModifyResult []string
		ItemBeforeSave   []string
		ItemForSave      []string
		ItemHtml         []string
	}
)

func (d DocType) Fld(fldName string) *FldType {
	for _, f := range d.Flds {
		if f.Name == fldName {
			return &f
		}
	}
	log.Fatalf("d.Fld: doc '%s' fld '%s' not found", d.Name, fldName)
	return nil
}

// место вызова разных доп функций для инициализации документа, после того как основные поля заполнены
func (d *DocType) Init() {
	// проверяем что есть поле title
	//isExist := false
	//for _, fld := range d.Flds {
	//	if fld.Name == "title" {
	//		isExist = true
	//	}
	//}
	//if !isExist {
	//	log.Fatalf("doc '%s' missed field 'title'", d.Name)
	//}

	d.Filli18n()
	for i := range d.Flds {
		d.Flds[i].Doc = d
	}
	// если есть табы и к документу можно присоединять задачи, то прописываем миксин
	if d.IsTaskAllowed && len(d.Vue.Tabs) > 0 {
		if d.Vue.Mixins == nil {
			d.Vue.Mixins = map[string][]VueMixin{}
		}
		if d.Vue.Mixins["docItemWithTabs"] == nil {
			d.Vue.Mixins["docItemWithTabs"] = []VueMixin{}
		}
		d.Vue.Mixins["docItemWithTabs"] = append(d.Vue.Mixins["docItemWithTabs"], VueMixin{Title: "taskList", Import: "../../mixins/taskList"})
	}
}

func (d DocType) PgName() string {
	return snaker.CamelToSnake(d.Name)
}

func (d DocType) NameCamelCase() string {
	return snaker.SnakeToCamel(d.Name)
}

func (d DocType) IsBitrixIntegration() bool {
	return len(d.Integrations.Bitrix.UrlName) > 0
}

func (d DocType) IsBitrixIntegrationDebugMode() bool {
	return d.Integrations.Bitrix.IsDebugMode
}

func (d DocType) IsOdataIntegration() bool {
	return len(d.Integrations.Odata.Name) > 0
}

func (d DocType) IsOdataIntegrationDebugMode() bool {
	return d.Integrations.Odata.IsDebugMode
}

func (d *DocType) AddFld(fld FldType)  {
	if d.Flds == nil {
		d.Flds = []FldType{}
	}
	d.Flds = append(d.Flds, fld)
}

// sugar для добавление компоненты во vue
// имя шаблона. Например, docItem
func (d *DocType) AddVueComposition(tmpName, compName string) {
	importName := "comp" + utils.UpperCaseFirst(compName)
	parentPath := "./comp/"
	if len(d.Vue.Tabs) > 0 {
		parentPath = "../../comp/"
	}
	importAddress := parentPath + compName + ".vue"
	dTemplateName := "webClient_comp_" + compName + ".vue"
	// добавляем в список компонент
	if d.Vue.Components == nil {
		d.Vue.Components = map[string]map[string]string{}
	}
	if d.Vue.Components[tmpName] == nil {
		d.Vue.Components[tmpName] = map[string]string{}
	}
	d.Vue.Components[tmpName][importName] = importAddress
	// добавляем в список шаблонов для загрузки
	if d.Templates == nil {
		d.Templates = map[string]*DocTemplate{}
	}
	d.Templates[dTemplateName] = &DocTemplate{}
}

// sugar для добавления табов и задач к документу
func (d *DocType) AddVueTaskAndTabs() {
	// в шаблон vue добавляем табы
	d.Vue.Tabs = []VueTab{
		{"info", "инфо", "tabInfo.vue", "assignment", "", ""},
		{"tasks", "задачи", "tabTasks.vue", "alarm", ":list='taskListForRender'", "<q-badge v-if='taskListForRender.length>0' color='red' floating>{{taskListForRender.length}}</q-badge>"},
	}
	// указываем признак, что к документу можно прикреплять задачи
	d.IsTaskAllowed = true
}

// добаление свойства рекурсии - добавляются поля и проставляется признак
func (d *DocType) SetIsRecursion(title string) {
	d.IsRecursion = true
	d.Flds = append(d.Flds,
		FldType{Name: "parent_id", NameRu: "родитель", Type: FldTypeInt, Sql: FldSql{Ref: d.Name, IsSearch: true, IsNotUpdatable: true}},
		FldType{Name: "is_folder", NameRu: "признак, что является группой", Type: FldTypeBool, Sql: FldSql{IsNotUpdatable: true}}, )
	d.Vue.I18nAdd("recursiveListTitle", title)
}

func (dv *DocVue) I18nAdd(titleEn, titleRu string) {
	if dv.I18n == nil {
		dv.I18n = map[string]string{}
	}
	dv.I18n[titleEn] = titleRu
}
