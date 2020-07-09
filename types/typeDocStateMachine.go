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
	}

	DocSmState struct {
		Title      string
		TitleRu    string
		Actions    []DocSmAction
		UpdateFlds []FldType // поля, которые можно редактировать в этом стейте
		IconSrc    string
	}

	DocSmAction struct {
		From       string
		To         string
		Label      string
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

func (DocSm) TmplSqlActionPrintAfterHook(d DocType) string  {
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
	if len(res)>0 {
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
		fldArr := []string{}
		for _, f := range st.UpdateFlds {
			fldArr = append(fldArr, fmt.Sprintf("'%s'", f.Name))
		}
		res = fmt.Sprintf("%s\t\twhen '%s' then\n\t\t\tupdateFlds = ARRAY [%s]::text[];\n", res, st.Title, strings.Join(fldArr, ", "))
	}
	return res
}

func (st DocSmState) GetStateUpdateFldsGrid() func() [][]FldType {
	res := [][]FldType{}
	for _, f := range st.UpdateFlds {
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