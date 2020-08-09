package templates

import (
	"github.com/pepelazz/projectGenerator/utils"
	t "github.com/pepelazz/projectGenerator/types"
	"text/template"
)

func stateMachineReadTmplAction(funcMap template.FuncMap, path ...string) *template.Template {
	funcMap["tmplSqlActionPrintCaseBlock"] = t.DocSm{}.TmplSqlActionPrintCaseBlock
	funcMap["tmplSqlActionPrintRefUpdateBlock"] = t.DocSm{}.TmplSqlActionPrintRefUpdateBlock
	funcMap["tmplSqlActionPrintRefUpdateVarDeclare"] = t.DocSm{}.TmplSqlActionPrintRefUpdateVarDeclare
	funcMap["tmplSqlActionPrintAfterHook"] = t.DocSm{}.TmplSqlActionPrintAfterHook

	tmpls, err := template.New("").Funcs(funcMap).Delims("[[", "]]").ParseFiles(path...)
	utils.CheckErr(err, "stateMachineReadTmplAction")
	for _, tmpl := range tmpls.Templates() {
		return tmpl
	}
	return nil
}

func stateMachineReadTmplUpdate(funcMap template.FuncMap, path ...string) *template.Template {
	funcMap["tmplSqlUpdatePrintCaseBlock"] = t.DocSm{}.TmplSqlUpdatePrintCaseBlock

	tmpls, err := template.New("").Funcs(funcMap).Delims("[[", "]]").ParseFiles(path...)
	utils.CheckErr(err, "stateMachineReadTmplUpdate")
	for _, tmpl := range tmpls.Templates() {
		return tmpl
	}
	return nil
}

func stateMachineReadTmplWebclientItem(funcMap template.FuncMap, path ...string) *template.Template {
	tmpls, err := template.New("").Funcs(funcMap).Delims("[[", "]]").ParseFiles(path...)
	utils.CheckErr(err, "stateMachineReadTmplWebclientItem")
	for _, tmpl := range tmpls.Templates() {
		return tmpl
	}
	return nil
}
