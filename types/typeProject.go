package types

import (
	"fmt"
	"github.com/pepelazz/nla_framework/utils"
	"github.com/serenize/snaker"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type (
	ProjectType struct {
		Name                     string
		Docs                     []DocType
		DistPath                 string
		Config                   ProjectConfig
		Vue                      ProjectVue
		Sql                      ProjectSql
		Go                       ProjectGo
		Roles                    []ProjectRole // список ролей в проекте
		IsDebugMode              bool
		OverridePathForTemplates map[string]string // map для замены путей к исходным файлам. Ключ - путь к генерируемому файлу, значение - новый путь к исходному файлу.
		I18n I18nType
	}
	ProjectConfig struct {
		Logo             string
		LocalProjectPath string
		Auth             AuthConfig
		Postgres         PostrgesConfig
		WebServer        WebServerConfig
		Email            EmailConfig
		DevMode          DevModeConfig
		Vue              VueConfig
		Bitrix           BitrixConfig
		Telegram         TelegramConfig
		Odata            OdataConfig
		Yandex           YandexConfig
		User             UserConfig
		Backup           BackupConfig
		Docker 			 DockerConfig
		Graylog 		 GraylogConfig
	}
	AuthConfig struct {
		ByEmail               bool // дефолт - авторизация по email
		ByPhone               bool // авторизация по номеру телефона
		SqlHooks              AuthConfigSqlHooks
		IsPassStepWaitingAuth bool // возможность отключить статус waiting_auth для вновь зарегестрированных пользователей
		SmsService            AuthConfigSmsService
		UserSqlFunction []string // дополнительные sql функции для таблицы User
	}
	AuthConfigSqlHooks struct {
		CheckIsUserExist []string
	}
	AuthConfigSmsService struct {
		Url      string // url для отправки sms,в которой первый параметр телефон, второй токен.  https://smsc.ru/sys/send.php?login=&psw=&phones=+%s&mes=%s
		CheckErr string
	}
	// дополнительные настройки для таблицы users
	UserConfig struct {
		Roles UserConfigRolesForMethods
	}
	UserConfigRolesForMethods struct {
		UserList []string
		UserUpdate []string
	}
	PostrgesConfig struct {
		DbName   string
		Port     int64
		Password string
		Host string
		TimeZone string // время для сервера default 'Europe/Moscow' (например 'Asia/Novosibirsk')
		Version string // версия Postgres, по дефолту 12
	}
	WebServerConfig struct {
		Port     int64
		Url      string
		Path     string
		Ip       string
		Username string // root или ...
		SshPort  int64
	}
	DevModeConfig struct {
		IsDocker bool
	}
	EmailConfig struct {
		Sender     string
		Password   string
		Host       string
		Port       int64
		SenderName string
		IsSendWithEmptySender bool
	}
	ProjectVue struct {
		UiAppName string
		UiAppLogoOnly string
		Routes    [][]string
		Menu      []VueMenu
		Hooks     ProjectVueHooks
		// кастомные шаблоны сообщений, которые выводятся в правом боковом списке
		MessageTmpls             []ProjectVueMessageTmpl
		IsHideTaskToolbar        bool // не показывааем боковое меню с задачами
		IsHideMessageToolbar     bool // не показывааем боковое меню с сообщениями
		IsHideUserAvatarUploader bool // не даем возможность пользователям загружать аватарки
	}
	VueMenu struct {
		DocName  string // если указано docName, то url и иконка копируются из описания документа
		Icon     string
		Text     string
		Url      string
		IsFolder bool
		LinkList []VueMenu
		Roles    []string
	}

	VueConfig struct {
		DadataToken string
		QuasarVersion int // версия quasar-framework 1, 2
	}

	I18nType struct {
		IsExist bool
		LangList []string
		Data  map[string]map[string]map[string]string //RU : message : save: 'сохранить'
	}


	BitrixConfig struct {
		ApiUrl       string
		UserId       string
		WebhookToken string
	}

	TelegramConfig struct {
		BotName string
		Token   string
	}

	OdataConfig struct {
		Url              string
		Login            string
		Password         string
		ExchangePlanName string
		ExchangePlanGuid string
	}

	YandexConfig struct {
		MetrikaId string
	}

	ProjectSql struct {
		Methods     map[string][]DocSqlMethod // имя документа и список методов. Например "task": []{"task_by_deal"}
		InitialData []string                  // данные при первоначальной загрузке
	}
	ProjectGo struct {
		JobList []string // список job'ов
		Routes  ProjectGoRoutes
	}
	ProjectRole struct {
		Name   string
		NameRu string
	}
	ProjectGoRoutes struct {
		Imports []string
		NotAuth []string // роуты вне блока, требующего авторизации
		Api []string // роуты в блоке, требующего авторизации
		Static []string // роуты в static блоке
	}
	ProjectVueHooks struct {
		Profile ProjectVueHooksProfile
	}
	ProjectVueHooksProfile struct {
		Flds string
	}
	// компонента - кастомный шаблон сообщения
	ProjectVueMessageTmpl struct {
		CompName string // название компоненты
		CompPath string // путь к файлу компоненты
	}
	BackupConfig struct {
		ToYandexDisk BackupConfigYandexDisk
	}
	BackupConfigYandexDisk struct {
		Token              string
		Path               string
		PostgresDockerName string
		Period             int // периодичность в минутах
		FilesCount         int // количество последних файлов, которое остается на яндекс сервере. Остальные удаляются
	}
	DockerConfig struct {
		AfterCopy []string // дополнительные строки для копирования в Dockerfile
		Volumes []string // маппинг директорий
	}

	GraylogConfig struct {
		Host string
		Port int
	}
)

func (p *ProjectType) GetDocByName(docName string) *DocType {
	for _, d := range p.Docs {
		if d.Name == docName {
			return &d
		}
	}
	return nil
}

// заполняем поля темплейтов - из короткой формы записи в полную
func (p *ProjectType) FillDocTemplatesFields() {
	for i, d := range p.Docs {
		if d.Templates == nil {
			d.Templates = map[string]*DocTemplate{}
		}
		for tName, t := range d.Templates {
			// прописываем полный путь к файлу шаблона
			if len(t.Source) == 0 {
				// учитывааем что возможен префикс, если папка с документом вложена в другую папку
				pathPrefix := ""
				if len(d.PathPrefix) > 0 {
					pathPrefix = d.PathPrefix + "/"
				}
				t.Source = fmt.Sprintf("%s%s/tmpl/%s", pathPrefix, snaker.SnakeToCamelLower(d.Name), tName)
			}
			// если не указан конечный путь, то формируем его исходя из ключа шаблона (например webClient_comp_...)
			if len(t.DistPath) == 0 {
				params := map[string]string{}
				if len(d.Vue.Path) > 0 {
					params["doc.Vue.Path"] = d.Vue.Path
				}
				distPath, distFilename := utils.ParseDocTemplateFilename(d.Name, tName, p.DistPath, i, params)
				t.DistFilename = distFilename
				t.DistPath = distPath
			}
		}
		p.Docs[i] = d
	}
}

// заполняем незаполненные поля для Vue
func (p *ProjectType) FillVueFlds() {
	for i, d := range p.Docs {
		for j, fld := range d.Flds {
			// если NameRu не заполнено, то копируем из fld
			if len(fld.Vue.NameRu) == 0 {
				p.Docs[i].Flds[j].Vue.NameRu = fld.NameRu
			}
			// заполняем IsRequired
			if fld.Sql.IsRequired {
				p.Docs[i].Flds[j].Vue.IsRequired = fld.Sql.IsRequired
			}
			// заполняем незаполненные поля в extension
			for k, _ := range fld.Vue.Ext {
				// если в параметрах есть pathUrl и поле является Ref, это значит надо заполнить route к доккументу, на который идет ссылка + ссылка на аватарку
				if k == "pathUrl" && len(fld.Sql.Ref) > 0 {
					for _, dRef := range p.Docs {
						if dRef.Name == fld.Sql.Ref {
							fld.Vue.Ext["pathUrl"] = "/" + dRef.Vue.RouteName
							fld.Vue.Ext["avatar"] = dRef.Vue.MenuIcon
						}
					}
					// для ссылки на таблицу user проставляем иконки отдельно
					if fld.Sql.Ref == "user" {
						fld.Vue.Ext["pathUrl"] = "/users"
						fld.Vue.Ext["avatar"] = "image/users.svg"
					}
				}
				if k == "addNewUrl" && len(fld.Sql.Ref) > 0 {
					for _, dRef := range p.Docs {
						if dRef.Name == fld.Sql.Ref {
							fld.Vue.Ext["addNewUrl"] = "/" + dRef.Vue.RouteName + "/new"
						}
					}
				}
			}
		}
	}
}

// заполняем боковое меню для Vue
func (p *ProjectType) FillSideMenu() {
	if p.Vue.Menu == nil {
		log.Fatalf("ProjectType.FillSideMenu p.Vue.Menu == nil")
	}
	for i, v := range p.Vue.Menu {
		if len(v.DocName) > 0 {
			d := p.GetDocByName(v.DocName)
			if d == nil {
				log.Fatalf("ProjectType.FillSideMenu p.GetDocByName doc '%s' not found", v.DocName)
			}
			if len(v.Icon) == 0 {
				p.Vue.Menu[i].Icon = d.Vue.MenuIcon
			}
			if len(v.Url) == 0 {
				p.Vue.Menu[i].Url = d.Vue.RouteName
			}
			if len(v.Text) == 0 {
				// i18n_ признак, чтобы различать текст и локализацию в js меню
				p.Vue.Menu[i].Text = fmt.Sprintf("i18n_menu.%s", d.Name)
				// если есть локализованное название для списка, то используем его (там множественное число). Если нет, то название документа
				//if title, ok := d.Vue.I18n["listTitle"]; ok {
				//	p.Vue.Menu[i].Text = title
				//} else {
				//	p.Vue.Menu[i].Text = utils.UpperCaseFirst(d.NameRu)
				//}
			}
			if len(v.Roles) == 0 {
				p.Vue.Menu[i].Roles = d.Vue.Roles
			}
		}
		if v.IsFolder {
			for j, v1 := range v.LinkList {
				if len(v1.DocName) > 0 {
					d := p.GetDocByName(v1.DocName)
					if d == nil {
						log.Fatalf("ProjectType.FillSideMenu p.GetDocByName doc '%s' not found", v1.DocName)
					}
					if len(v1.Icon) == 0 {
						p.Vue.Menu[i].LinkList[j].Icon = d.Vue.MenuIcon
					}
					if len(v1.Url) == 0 {
						p.Vue.Menu[i].LinkList[j].Url = d.Vue.RouteName
					}
					if len(v1.Text) == 0 {
						p.Vue.Menu[i].LinkList[j].Text = fmt.Sprintf("i18n_menu.%s", d.Name)
						// если есть локализованное название для списка, то используем его (там множественное число). Если нет, то название документа
						//if title, ok := d.Vue.I18n["listTitle"]; ok {
						//	p.Vue.Menu[i].LinkList[j].Text = title
						//} else {
						//	p.Vue.Menu[i].LinkList[j].Text = utils.UpperCaseFirst(d.NameRu)
						//}
					}
					if len(v1.Roles) == 0 {
						p.Vue.Menu[i].LinkList[j].Roles = d.Vue.Roles
					}
				}
			}
		}
	}
}

func (p *ProjectType) AddVueRoute(urlName, compPath string) {
	if p.Vue.Routes == nil {
		p.Vue.Routes = [][]string{}
	}
	p.Vue.Routes = append(p.Vue.Routes, []string{urlName, compPath})
}

// генерим сетку для Vue
func (p *ProjectType) GenerateGrid() {
	for i, d := range p.Docs {
		d.Vue.Grid = makeGrid(d)
		p.Docs[i] = d
	}
}

// признак что есть интеграция с Битрикс
func (p ProjectType) IsBitrixIntegration() bool {
	return len(p.Config.Bitrix.ApiUrl) > 0
}

// признак что есть интеграция с Telegram
func (p ProjectType) IsTelegramIntegration() bool {
	return len(p.Config.Telegram.Token) > 0
}

// признак что есть интеграция с Odata
func (p ProjectType) IsOdataIntegration() bool {
	return len(p.Config.Odata.Url) > 0
}

// признак что настроен бэкап на яндекс диск
func (p ProjectType) IsBackupOnYandexDisk() bool {
	return len(p.Config.Backup.ToYandexDisk.Token) > 0
}

// печать списка go jobs
func (p ProjectType) PrintGoJobList() string {
	res := ""
	if len(p.Go.JobList) > 0 {
		return strings.Join(p.Go.JobList, "\n")
	}
	return res
}

// печать списка go jobs
func (p ProjectType) PrintJsRoles() string {
	res := "{label: i18n.global.t('user.role_admin'), value: 'admin'},\n"
	isOverrideStudent := false
	for _, r := range p.Roles {
		if r.Name == "student" {
			isOverrideStudent = true
			res = res + fmt.Sprintf("\t\t{label: i18n.global.t('user.role_%s'), value: '%s'},\n", r.Name, r.Name)
		}
	}
	// отдельно рассматриваем возможность переопределения дефолтной роли student
	if !isOverrideStudent {
		res = res + fmt.Sprintf("\t\t{label: i18n.global.t('user.role_student'), value: 'student'},\n")
	}
	for _, r := range p.Roles {
		if r.Name != "student" {
			res = res + fmt.Sprintf("\t\t{label: i18n.global.t('user.role_%s'), value: '%s'},\n", r.Name, r.Name)
		}
	}
	return res
}

// FillLocalPath Если не указан путь к локальному проекту, то вычисляем его автоматически
func (p *ProjectType) FillLocalPath() string {
	// для старого варианта, когда проект находится в директории GOPATH
	if len(p.Config.LocalProjectPath) == 0 {
		if _, err := os.Stat("../go.mod"); os.IsNotExist(err) {
			// путь к локальной директории
			path, _ := filepath.Abs("./")
			// находим gopath
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				gopath = build.Default.GOPATH
			}
			// убираем из начала пути gopath
			path = strings.TrimPrefix(path, gopath)
			// приводим разделитель пути к unix стилю
			path = strings.Replace(path, string(os.PathSeparator), "/", -1)
			// убираем из начала еще src
			path = strings.TrimPrefix(path, "/src/")
			// убираем из конца projectTemplate и добавляем src
			path = strings.TrimSuffix(path, "/projectTemplate") + "/src"
			p.Config.LocalProjectPath = path
		} else {
			// для нового варианта, когда проект создается с модулями
			path, _ := os.Getwd()
			path = strings.Replace(path, string(os.PathSeparator), "/", -1)
			arr := strings.Split(strings.TrimSuffix(path, "/projectTemplate"), "/")
			p.Config.LocalProjectPath = arr[len(arr)-1] + "/src"
		}

	}
	return p.Config.LocalProjectPath
}

func (p ProjectType) PrintApiCallPgFuncMethods() string {
	res := ""
	printPgMethod := func(m DocSqlMethod) {
		var roles string
		if len(m.Roles) > 0 {
			roles = fmt.Sprintf(`"%s"`, strings.Join(m.Roles, `", "`))
		}
		res = fmt.Sprintf("%s\n\t\tPgMethod{\"%s\", []string{%s}, nil, BeforeHookAddUserId},", res, m.Name, roles)
	}
	if project.Sql.Methods != nil {
		for _, v := range project.Sql.Methods {
			for _, m := range v {
				printPgMethod(m)
			}
		}
	}
	for _, d := range project.Docs {
		for _, m := range d.Sql.Methods {
			printPgMethod(*m)
		}
	}
	return res
}

// PrintProcessPgErrorMsgs печать перевода сообщений из postgres
// например `violates unique constraint "day_already_exist"` -> "отчет на данную дату уже существует"
func (p ProjectType) PrintProcessPgErrorMsgs() string {
	//if strings.Contains(err.Error(), `violates unique constraint "day_already_exist"`) {
	//	return "отчет на данную дату уже существует"
	//}
	var res []string
	for _, d := range project.Docs {
		for _, m := range d.Sql.UniqConstrains {
			if len(m.Message)>0 {
				res = append(res, fmt.Sprintf("\tif strings.Contains(err.Error(), `violates unique constraint \"%s\"`) {\n" +
					"\t\treturn `%s`\n" +
					"\t}", m.Name, m.Message))
			}
		}
	}
	return strings.Join(res, "\n")
}

func (p *ProjectType) GetQuasarVersion() int  {
	if p.Config.Vue.QuasarVersion == 0 {
		return 1
	}
	return p.Config.Vue.QuasarVersion
}

func (p *ProjectType) AddI18n(lang, prefix, key, value string)  {
	if len(p.I18n.Data)==0 {
		p.I18n.Data = map[string]map[string]map[string]string{}
	}
	if len(p.I18n.Data[lang]) == 0 {
		p.I18n.Data[lang] = map[string]map[string]string{}
	}
	if len(p.I18n.Data[lang][prefix]) == 0 {
		p.I18n.Data[lang][prefix] = map[string]string{}
	}
	p.I18n.Data[lang][prefix][key] = value
}

