package projectGenerator

import (
	"encoding/json"
	"fmt"
	"github.com/pepelazz/projectGenerator/templates"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	//"github.com/otiai10/copy"
)

type (
	// функция для модификации файлов при копировании из шаблона проекта в конечный проект
	copyFileModifyFunc func(path string, file []byte) []byte
)

var (
	project        types.ProjectType
	tmplMap        map[string]*template.Template
)

func readData(p types.ProjectType)  {
	project = p
	// проставляем localpath если он не заполнен
	project.FillLocalPath()
	// передаем project в папку types, чтобы иметь доступ из функций шаблонов к проекту
	templates.SetProject(&project)
	project.DistPath = "../src"
	project.FillDocTemplatesFields()
	project.GenerateGrid()
	project.FillVueFlds()
	// проставляем дефолтное время сервера, если не задано в настройках проекта
	if len(project.Config.Postgres.TimeZone) == 0 {
		project.Config.Postgres.TimeZone = "Europe/Moscow"
	}

	// передаем project в папку types, чтобы иметь доступ из функций шаблонов к проекту
	types.SetProject(&project)
}

func Start(p types.ProjectType, modifyFunc copyFileModifyFunc)  {
	// читаем данные для проекта
	readData(p)
	// читаем темплейты
	tmplMap = templates.ParseTemplates(project)

	// удаляем старые файлы
	removeOldFiles(project.DistPath)
	
	// генерим файлы для проекты
	templates.WriteProjectFiles(p, tmplMap)

	// генерим файлы для документов
	for _, d := range p.Docs {
		for _, dt := range d.Templates {
			err := templates.ExecuteToFile(dt.Tmpl, d, dt.DistPath, dt.DistFilename)
			utils.CheckErr(err, fmt.Sprintf("'%s' ExecuteToFile '%s'", d.Name, dt.DistFilename))
		}
	}

	// копируем файлы проекта (которые не шаблоны)
	err := copyFiles(p,"../../../pepelazz/projectGenerator/sourceFiles", "../", modifyFunc)
	utils.CheckErr(err, "Copy sourceFiles")

	templates.OtherTemplatesGenerate(project)
}

// функция для копирования файлов с возможностью модификаации содержимого файлов
func copyFiles(p types.ProjectType, source, dist string, modifyFunc copyFileModifyFunc) (err error)  {
	err = filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				file, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				// для windows заменяем слэши в пути на обратные
				dirPath := strings.TrimSuffix(strings.TrimPrefix(strings.Replace(path, "\\", "/", -1), source), info.Name())
				// создаем директории
				err = os.MkdirAll(dist + dirPath, os.ModePerm)
				if err != nil {
					return err
				}
				// заменяем ссылки в go файлах
				if strings.HasSuffix(info.Name(), ".go") {
					file = []byte(strings.Replace(string(file), "github.com/pepelazz/projectGenerator", p.Config.LocalProjectPath, -1))
				}
				// добавляем названия pg методов в файл apiCallPgFunc.go
				if strings.HasSuffix(info.Name(), "apiCallPgFunc.go") {
					file = []byte(strings.Replace(string(file), "// for codeGenerate ##pgFuncList_slot1", printApiCallPgFuncMethods(), -1))
				}
				// изменение config.js
				if strings.HasSuffix(path, "app"+string(os.PathSeparator)+"plugins"+string(os.PathSeparator)+"config.js") {
					file = configJsModify(p, file)
				}
				// изменение sidemenu/index.vue
				if strings.HasSuffix(path, "components"+string(os.PathSeparator)+"sidemenu"+string(os.PathSeparator)+"index.vue") {
					file = []byte(strings.Replace(string(file), "// for codeGenerate ##sidemenu_slot1", sidemenuJsModify(), -1))
				}
				// изменение routes.js
				if strings.HasSuffix(path, "src"+string(os.PathSeparator)+"router"+string(os.PathSeparator)+"routes.js") {
					file = []byte(strings.Replace(string(file), "// for codeGenerate ##routes_slot1", routesJsModify(), -1))
				}
				// изменение _Task/main.toml - дописываем дополнительные методы
				if strings.HasSuffix(path, "_Task"+string(os.PathSeparator)+"main.toml") {
					insertText := "# for codeGenerate task_methods_slot"
					if project.Sql.Methods != nil {
						isMethodsExist := false
						for _, v := range project.Sql.Methods["task"] {
							isMethodsExist = true
							insertText = fmt.Sprintf("%s\n\t\"%s\",", insertText, v.Name)
						}
						if isMethodsExist {
							file = []byte(strings.Replace(string(file), "# for codeGenerate task_methods_slot", insertText, -1))
						}
					}
				}
				// изменение index.template.html
				if strings.HasSuffix(path, "src"+string(os.PathSeparator)+"index.template.html") {
					file = []byte(strings.Replace(string(file), "[[appName]]", p.Name, -1))
				}
				// изменение loginPage.vue и home.vue
				if strings.HasSuffix(path, "loginPage.vue") || strings.HasSuffix(path, "home.vue") {
					file = []byte(strings.Replace(string(file), "[[appLogoSrc]]", p.Config.Logo, -1))
				}
				// проставляем Config.Postgres.TimeZone в sql файлах
				if strings.HasSuffix(path, ".sql") {
					file = []byte(strings.Replace(string(file), "[[Config.Postgres.TimeZone]]", p.Config.Postgres.TimeZone, -1))
				}
				// добавляем в триггер для задач дополнительные блоки
				if strings.HasSuffix(path, "trigger_task_update_table_name.sql") {
					insertText := "-- for codeGenerate #trigger_task_update_table_name_slot"
					if project.Sql.Methods != nil {
						isMethodsExist := false
						for _, v := range project.Sql.Methods["task"] {
							isMethodsExist = true
							if txt, ok := v.Params["trigger_task_update_table_name.sql"]; ok {
								insertText = fmt.Sprintf("%s\n%s", insertText, txt)
							}
						}
						if isMethodsExist {
							file = []byte(strings.Replace(string(file), "-- for codeGenerate #trigger_task_update_table_name_slot", insertText, -1))
						}
					}
				}
				// применяем модификатор для текста файла
				if modifyFunc != nil {
					file = modifyFunc(dirPath + info.Name(), file)
				}
				// записываем файл по новому пути
				err = ioutil.WriteFile(dist + dirPath + info.Name(), file, 0644)
				if err != nil {
					return err
				}
			}
			return nil
		})
	return
}

func removeOldFiles(distPath string)  {
	// удаляем модели в sql, потому что могула изменится нумерация файлов и тогдда риск дублирования
	err := os.RemoveAll(distPath + "/sql/model")
	utils.CheckErr(err, "removeOldFiles")
}

func printApiCallPgFuncMethods() (res string)  {
	res = "// for codeGenerate ##pgFuncList_slot1"
	printPgMethod := func(m types.DocSqlMethod) {
		var roles string
		if len(m.Roles) > 0{
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
	return
}

func configJsModify(p types.ProjectType, file []byte) (res []byte)  {
	jsTablesForTask := func() string  {
		res := map[string]string{}
		for _, d := range project.Docs {
		if d.IsTaskAllowed {
		res[d.Name] = d.NameRu
	}
	}
		jsonStr, _ := json.Marshal(res)
		return string(jsonStr)
	}
	breadcrumbIcons := []string{}
	for _, d := range p.Docs {
		if len(d.Vue.BreadcrumbIcon)>0 {
			breadcrumbIcons = append(breadcrumbIcons, fmt.Sprintf("%s: '%s'", d.Name, d.Vue.BreadcrumbIcon))
		}
	}
	fileStr := string(file)
	fileStr = strings.Replace(fileStr, "[[appName]]", p.Name, -1)
	fileStr = strings.Replace(fileStr, "[[uiAppName]]", p.Vue.UiAppName, -1)
	fileStr = strings.Replace(fileStr, "[[webPort]]", fmt.Sprintf("%v", p.Config.WebServer.Port), -1)
	fileStr = strings.Replace(fileStr, "[[url]]",  strings.TrimPrefix(p.Config.WebServer.Url, "https://"), -1)
	fileStr = strings.Replace(fileStr, "[[logoSrc]]",  p.Config.Logo, -1)
	fileStr = strings.Replace(fileStr, "[[breadcrumbIcons]]", strings.Join(breadcrumbIcons, ",\n"), -1)
	// проставляем список таблиц, к которым можно прикреплять задачи
	fileStr = strings.Replace(fileStr, "[[codoGeneratedTablesForTask]]",  jsTablesForTask(), -1)
	return []byte(fileStr)
}

// функция по добавлению routes
func routesJsModify() string  {
	res := "// for codeGenerate ##routes_slot1"
	for _, r := range project.Vue.Routes {
		if len(r)<2 {
			log.Fatalf("routesJsModify project.Vue.Route route array %v length < 2", r)
		}
		res = fmt.Sprintf("%s\n\t{path: '/%s', component: () => import(`../app/components/%s`), props: true},", res, r[0], r[1])
		//{path: '/users/:id', component: () => import(`../app/components/users/item.vue`), props: true},
	}
	return res
}

// функция для построения бокового меню во Vue
func sidemenuJsModify() string  {
	res := "// for codeGenerate ##sidemenu_slot1"
	printMenuItem := func(m types.VueMenu) string{
		roles := ""
		if m.Roles != nil && len(m.Roles)>0 {
			roles = fmt.Sprintf(`'%s'`, strings.Join(m.Roles, `', '`))
		}
		return fmt.Sprintf("{icon: '%s', text: '%s', url: '/%s', roles: [%s]},\n", m.Icon, m.Text, m.Url, roles)
	}
	for _, m := range project.Vue.Menu {
		if !m.IsFolder {
			res = fmt.Sprintf("%s\n\t\t\t\t\t\t\t\t\t%s", res, printMenuItem(m))
			// {icon: 'people', text: 'Пользователи', url: '/users', role: ['admin']},
		} else {
			linkList := "\t\t\t\t\t\t\t\t[\n"
			for _, m1 := range m.LinkList {
				linkList = fmt.Sprintf("%s\t\t\t\t\t\t\t\t%s", linkList, printMenuItem(m1))
			}
			linkList = linkList + "],"
			roles := ""
			if m.Roles != nil && len(m.Roles)>0 {
				roles = fmt.Sprintf(`'%s'`, strings.Join(m.Roles, `', '`))
			}
			res = fmt.Sprintf("%s{isFolder: true, icon: '%s', text: '%s', roles: [%s], linkList: %s},\n", res, m.Icon, m.Text, roles, linkList)
		}
	}
	return res
}
