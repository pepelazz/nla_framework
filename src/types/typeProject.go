package types

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/src/utils"
)

type (
	ProjectType struct {
		Name string
		Docs []DocType
		DistPath string
	}
)

func (p *ProjectType) GetDocByName(docName string) *DocType  {
	for _, d := range p.Docs {
		if d.Name == docName {
			return &d
		}
	}
	return nil
}

// заполняем поля темплейтов - из короткой формы записи в полную
func (p *ProjectType) FillDocTemplatesFields()  {
	for i, d := range p.Docs {
		if d.Templates == nil {
			d.Templates = map[string]*DocTemplate{}
		}
		for tName, t := range d.Templates {
			// прописываем полный путь к файлу шаблона
			if len(t.Source) == 0 {
				t.Source = fmt.Sprintf("%s/tmpl/%s", d.Name, tName)
			}
			distPath, distFilename := utils.ParseDocTemplateFilename(d.Name, tName, p.DistPath, i)
			t.DistFilename = distFilename
			t.DistPath = distPath
		}
		p.Docs[i] = d
	}
}

// генерим сетку для Vue
func (p *ProjectType) GenerateGrid()  {
	for i, d := range p.Docs {
		d.Vue.Grid = makeGrid(d)
		p.Docs[i] = d
	}
}
