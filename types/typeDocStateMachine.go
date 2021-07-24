package types

import (
	"fmt"
	"github.com/serenize/snaker"
	"strings"
)

type (
	// State machine
	DocSm struct {
		States []*DocSmState
		Tmpls  DocSmTmpls
	}

	DocSmState struct {
		Title            string
		TitleRu          string
		Actions          []DocSmAction
		UpdateFlds       []FldType // поля, которые можно редактировать в этом стейте
		IconSrc          string
		FuncMapForCard   map[string]interface{}
		FuncMapForAction map[string]interface{}
	}

	DocSmAction struct {
		From       string
		To         string
		Label      string
		Icon       string
		UpdateFlds []FldType              // поля, которые заполняются при смене стейта
		Conditions []DocSmActionCondition // условия выполнения экшена
		Hooks      DocSmActionlHooks
	}

	DocSmActionlHooks struct {
		DeclareVars []string
		Before      []string
		After       []string
	}

	DocSmActionCondition struct {
		SqlText string
		VueIf   string
	}

	DocSmTmpls struct {
		ItemStateHeader string
		IsShowChat      bool
		Hooks DocSmTmplsHooks
	}
	DocSmTmplsHooks struct {
		AfterActionBtns []string
		ItemMethods []string
		BeforeChat []string
	}
)

func (DocSm) TmplSqlActionPrintCaseBlock(d DocType) string {
	res := ""
	for _, st := range d.StateMachine.States {
		for _, actn := range st.Actions {
			// собираем строку - copyToParamsFlds = '{amount}'::text[]
			copyToParamsFlds := []string{}
			for _, v := range st.UpdateFlds {
				copyToParamsFlds = append(copyToParamsFlds, v.Name)
			}
			// собираем строку - updateFlds = ARRAY ['amount', 'sum', 'state', 'comment']
			updateFlds := []string{"'state'"}
			for _, v := range actn.UpdateFlds {
				updateFlds = append(updateFlds, fmt.Sprintf("'%s'", v.Name))
			}
			conditions := ""
			for _, cond := range actn.Conditions {
				conditions = conditions + cond.SqlText + "\n"
			}
			beforeHooks := ""
			for _, hook := range actn.Hooks.Before {
				beforeHooks = beforeHooks + hook + "\n"
			}
			res = fmt.Sprintf(`%s
		when '%s_to_%s' then
			newStateName = '%[3]s';
			allowedStates = '{%[2]s}'::text[];
			copyToParamsFlds = '{%[4]s}'::text[];
			updateFlds = ARRAY [%[5]s]::text[];
			%[6]s 
			%[7]s
`, res, st.Title, actn.To, strings.Join(copyToParamsFlds, ", "), strings.Join(updateFlds, ", "), conditions, beforeHooks)
		}
	}
	return res
}

func (DocSm) TmplSqlActionPrintAfterHook(d DocType) string {
	res := ""
	for _, st := range d.StateMachine.States {
		for _, actn := range st.Actions {
			afterHooks := ""
			for _, hook := range actn.Hooks.After {
				afterHooks = afterHooks + hook + "\n"
			}
			if len(afterHooks) > 0 {
				res = fmt.Sprintf("%s\n\t\twhen '%s_to_%s' then \n\t\t\t\t%s", res, st.Title, actn.To, afterHooks)
			}
		}
	}
	if len(res) > 0 {
		res = fmt.Sprintf("case params->>'action_name'\n%s\n\t\telse\n\tend case;", res)
	}
	return res
}

func (DocSm) TmplSqlActionPrintRefUpdateBlock(d DocType) string {
	res := ""
	for _, fld := range d.Flds {
		if len(fld.Sql.Ref) > 0 {
			varName := snaker.SnakeToCamelLower(strings.TrimSuffix(fld.Name, "_id")) + "Title"
			refTableName := fld.Sql.Ref
			if refTableName == "user" {
				refTableName = "\"user\""
			}
			res = res + fmt.Sprintf(`
			-- в случае обновления ссылки добавляем название
			if copyFldName = '%[1]s' AND (rJson ->> copyFldName)::int notnull then
				select title into %[2]s from %[3]s where id = (rJson ->> copyFldName)::int;
				params = params || jsonb_set(params, '{options, states, 0}'::text[] || '{%[4]s_title}'::text[], to_jsonb(%[2]s));
			end if;
		`, fld.Name, varName, refTableName, strings.TrimSuffix(fld.Name, "_id"))
		}
	}
	return res
}

func (DocSm) TmplSqlActionPrintRefUpdateVarDeclare(d DocType) string {
	res := ""
	for _, fld := range d.Flds {
		if len(fld.Sql.Ref) > 0 {
			res = res + snaker.SnakeToCamelLower(strings.TrimSuffix(fld.Name, "_id")) + "Title TEXT;\n\t"
		}
	}
	return res
}

func (DocSm) TmplSqlUpdatePrintCaseBlock(d DocType) string {
	res := ""
	for _, st := range d.StateMachine.States {
		fldArr := []string{"'deleted'"}
		for _, f := range st.UpdateFlds {
			fldArr = append(fldArr, fmt.Sprintf("'%s'", f.Name))
		}
		res = fmt.Sprintf("%s\t\twhen '%s' then\n\t\t\tupdateFlds = ARRAY [%s]::text[];\n", res, st.Title, strings.Join(fldArr, ", "))
	}
	return res
}

func (st DocSm) GetFirstState() DocSmState {
	if len(st.States) > 0 {
		return *(st.States[0])
	}
	return DocSmState{}
}

func (st DocSmState) GetStateUpdateFldsGrid() func() [][]FldType {
	res := [][]FldType{}
	for _, f := range st.UpdateFlds {
		if f.Name == "state" {
			continue
		}
		rowNum := f.Vue.RowCol[0][0]
		// проставляем дефолтное значение readonly
		if len(f.Vue.Readonly) == 0 {
			f.Vue.Readonly = "isReadonly"
		}
		// автоматически увеличиваем массив в зависимости от количества строк
		for {
			if len(res) > rowNum {
				break
			}
			res = append(res, []FldType{})
		}
		res[rowNum-1] = append(res[rowNum-1], f)
	}
	return func() [][]FldType {
		return res
	}
}

func (action DocSmAction) GetUpdateFldsGrid() func() [][]FldType {
	res := [][]FldType{}
	for _, f := range action.UpdateFlds {
		if f.Name == "state" {
			continue
		}
		rowNum := f.Vue.RowCol[0][0]
		// автоматически увеличиваем массив в зависимости от количества строк
		for {
			if len(res) > rowNum {
				break
			}
			res = append(res, []FldType{})
		}
		res[rowNum-1] = append(res[rowNum-1], f)
	}
	return func() [][]FldType {
		return res
	}
}

func (sm *DocSm) GenerateTmpls(doc *DocType, params map[string]interface{}) {
	path := fmt.Sprintf("../../../pepelazz/nla_framework/templates/webClient/quasar_%v/doc/comp/stateMachine", doc.GetProject().GetQuasarVersion())
	cardTmplPath := path + "/cardTmpl.vue"
	actionBtnPath := path + "/actionBtn.vue"
	if params != nil {
		if p, ok := params["cardTmplPath"]; ok {
			cardTmplPath = p.(string)
		}
		if p, ok := params["actionBtnPath"]; ok {
			actionBtnPath = p.(string)
		}
	}
	for _, st := range sm.States {
		// шабоны cardTmpl для карточек состояний
		fileName := "state_" + st.Title + "_card.vue"
		stTitle := st.Title
		stTitleRu := st.TitleRu
		iconSrc := st.IconSrc
		if doc.Templates == nil {
			doc.Templates = map[string]*DocTemplate{}
		}
		doc.Templates[fileName] = &DocTemplate{
			Source:       cardTmplPath,
			DistPath:     fmt.Sprintf("../src/webClient/src/app/components/%s/comp", doc.Name),
			DistFilename: fileName,
			FuncMap: map[string]interface{}{
				"GetStateName": func() string { return stTitle },
				"GetLabel":     func() string { return stTitleRu },
				"GetIconSrc":   func() string { return iconSrc },
				//"GetUpdateFlds": func() []t.FldType {return updateFlds},
				"GetStateUpdateFldsGrid": st.GetStateUpdateFldsGrid(),
			},
		}
		// расширяем FuncMap функциями, которые указаны в шаблоне
		for k, v := range st.FuncMapForCard {
			doc.Templates[fileName].FuncMap[k] = v
		}
		// шабоны actionBtn для кнопок по переходу в новое состояние
		for _, actn := range st.Actions {
			fileName := st.Title + "_to_" + actn.To + "_btn.vue"
			actnLabel := actn.Label
			actnName := st.Title + "_to_" + actn.To
			actnIconSrc := actn.Icon
			updateFlds := []FldType{}
			// собираем цепочку условий vif
			vifArr := []string{}
			for _, v := range actn.Conditions {
				if len(v.VueIf) > 0 {
					vifArr = append(vifArr, v.VueIf)
				}
			}
			// собираем строчку с условиями v-if только в том случае если условия существуют
			vif := ""
			if len(vifArr) > 0 {
				vif = fmt.Sprintf("v-if=\"%s\"", strings.Join(vifArr, " && "))
			}

			for _, fld := range actn.UpdateFlds {
				if fld.Name == "state" {
					continue
				}
				updateFlds = append(updateFlds, fld)
			}
			doc.Templates[fileName] = &DocTemplate{
				Source:       actionBtnPath,
				DistPath:     fmt.Sprintf("../src/webClient/src/app/components/%s/comp", doc.Name),
				DistFilename: fileName,
				FuncMap: map[string]interface{}{
					"GetLabel":          func() string { return actnLabel },
					"GetIconSrc":        func() string { return actnIconSrc },
					"GetActionName":     func() string { return actnName },
					"GetUpdateFlds":     func() []FldType { return updateFlds },
					"GetUpdateFldsGrid": actn.GetUpdateFldsGrid(),
					"Vif":               func() string { return vif },
				},
			}
			// расширяем FuncMap функциями, которые указаны в шаблоне
			for k, v := range st.FuncMapForAction {
				doc.Templates[fileName].FuncMap[k] = v
			}
		}
	}
}
