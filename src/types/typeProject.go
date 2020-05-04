package types

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/src/utils"
	"github.com/serenize/snaker"
	"log"
)

type (
	ProjectType struct {
		Name     string
		Docs     []DocType
		DistPath string
		Config   ProjectConfig
		Vue ProjectVue
	}
	ProjectConfig struct {
		Logo string
		LocalProjectPath string
		Postgres  PostrgesConfig
		WebServer WebServerConfig
		Email     EmailConfig
	}
	PostrgesConfig struct {
		DbName   string
		Port     int64
		Password string
	}
	WebServerConfig struct {
		Port     int64
		Url   	string
		Path   	string
		Ip   	string
		Username   	string // root или ...
	}
	EmailConfig struct {
		SenderName     string
	}
	ProjectVue struct {
		UiAppName string
		Routes [][]string
		Menu []VueMenu
	}
	VueMenu struct {
		DocName string // если указано docName, то url и иконка копируются из описания документа
		Icon string
		Text string
		Url string
		IsFolder bool
		LinkList []VueMenu
		Roles []string
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
				t.Source = fmt.Sprintf("%s/tmpl/%s", snaker.SnakeToCamel(d.Name), tName)
			}
			distPath, distFilename := utils.ParseDocTemplateFilename(d.Name, tName, p.DistPath, i)
			t.DistFilename = distFilename
			t.DistPath = distPath
		}
		p.Docs[i] = d
	}
}

// заполняем незаполненные поля для Vue
func (p *ProjectType) FillVueFlds()  {
	for i, d := range p.Docs {
		for j, fld := range d.Flds {
			// если NameRu не заполнено, то копируем из fld
			if len(fld.Vue.NameRu) == 0 {
				p.Docs[i].Flds[j].Vue.NameRu = fld.NameRu
			}
			// заполняем IsRequired
			if fld.Sql.IsRequired {
				p.Docs[i].Flds[j].Vue.IsRequred = fld.Sql.IsRequired
			}
		}
	}
}

// заполняем боковое меню для Vue
func (p *ProjectType) FillSideMenu() {
	if (p.Vue.Menu ==  nil) {
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
				// если есть локализованное название для списка, то используем его (там множественное число). Если нет, то название документа
				if title, ok := d.Vue.I18n["listTitle"]; ok {
					p.Vue.Menu[i].Text = title
				} else {
					p.Vue.Menu[i].Text = utils.UpperCaseFirst(d.NameRu)
				}
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
						p.Vue.Menu[i].LinkList[j].Text = d.NameRu
					}
					if len(v1.Roles) == 0 {
						p.Vue.Menu[i].LinkList[j].Roles = d.Vue.Roles
					}
				}
			}
		}
	}
}

// генерим сетку для Vue
func (p *ProjectType) GenerateGrid() {
	for i, d := range p.Docs {
		d.Vue.Grid = makeGrid(d)
		p.Docs[i] = d
	}
}
