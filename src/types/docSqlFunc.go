package types

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/src/utils"
	"github.com/serenize/snaker"
	"log"
	"strings"
)

// main.toml печать списка полей
func (d DocType) PrintSqlModelFlds() (res string) {
	arr := []string{"\t{name=\"id\",\t\t\ttype=\"serial\"}"}
	for _, fld := range d.Flds {
		if len(fld.Name) == 0 || fld.Sql.IsOptionFld {
			continue
		}
		arr = append(arr, fld.PrintPgModel())
	}
	if d.Sql.IsSearchText {
		arr = append(arr, "\t{name=\"search_text\",\t\t\ttype=\"text\",\tcomment=\"колонка для поиска\"}")
	}
	arr = append(arr, "\t{name=\"options\",\t\t\t\ttype=\"jsonb\",\tcomment=\"разные дополнительные параметры\"}")
	arr = append(arr, "\t{name=\"created_at\",\t\t\t\ttype=\"timestamp\",\text=\"with time zone\"}")
	arr = append(arr, "\t{name=\"updated_at\",\t\t\t\ttype=\"timestamp\",\text=\"with time zone\"}")
	arr = append(arr, "\t{name=\"deleted\",\t\t\t\ttype=\"bool\",\text=\"not null default false\"}")

	res = fmt.Sprintf("fields = [\n%s\n]", strings.Join(arr, ",\n"))
	return
}

// main.toml печать fk_constraints
func (d DocType) PrintSqlModelFkConstraints() (res string) {

	// формирование строчки для fkConstraints в случае если таблица является связью двух таблиц и эта связь уникальна
	printSqlUniqLinkConstraint := func(d DocType) string {
		flds := []FldType{}
		for _, fld := range d.Flds {
			if len(fld.Sql.Ref) > 0 {
				flds = append(flds, fld)
			}
		}
		if len(flds) > 1 {
			return fmt.Sprintf(`	{name="%s_already_exist", ext="UNIQUE (%s, %s)"},`, snaker.CamelToSnake(d.Name), flds[0].Name, flds[1].Name)
		}
		return ""
	}

	arr := []string{}
	for _, fld := range d.Flds {
		// поле ссылка на другую таблицу
		if len(fld.Sql.Ref) > 0 {
			if fld.Sql.Ref != "user" {
				arr = append(arr, fmt.Sprintf("\t{fld=\"%s\", ref=\"%s\", fk=\"id\"}", fld.Name, fld.Sql.Ref))
			} else {
				arr = append(arr, fmt.Sprintf(`{fld="%s", ref="\"%s\"", fk="id"}`, fld.Name, fld.Sql.Ref))
			}
		}
		// ограничение на уникальность
		if fld.Sql.IsUniq {
			arr = append(arr, fmt.Sprintf("\t{name=\"%s_%s_already_exist\", ext=\"UNIQUE (%s)\"}", d.Name, fld.Name, fld.Name))
		}
	}
	// ограничение на уникальность связи между таблицами
	if d.Sql.IsUniqLink {
		arr = append(arr, printSqlUniqLinkConstraint(d))
	}
	if len(arr) > 0 {
		res = fmt.Sprintf("fkConstraints = [\n%s\n]", strings.Join(arr, ",\n"))
	}
	return
}

// main.toml печать triggers
func (d DocType) PrintSqlModelTriggers() (res string) {
	arr := []string{fmt.Sprintf("\t{name=\"%s_created\", when=\"before insert or update\", ref=\"for each row\", funcName=\"builtin_fld_update\"}", d.Name)}
	if d.Sql.IsBeforeTrigger {
		arr = append(arr, fmt.Sprintf("\t{name=\"%s_trigger_before\", when=\"before insert or update\", ref=\"for each row\", funcName=\"%s_trigger_before\"}", d.Name, d.Name))
	}
	if d.Sql.IsAfterTrigger {
		arr = append(arr, fmt.Sprintf("\t{name=\"%s_trigger_after\", when=\"after insert or update\", ref=\"for each row\", funcName=\"%s_trigger_after\"}", d.Name, d.Name))
	}
	if len(arr) > 0 {
		res = fmt.Sprintf("triggers = [\n%s\n]", strings.Join(arr, ",\n"))
	}
	return
}

// main.toml печать methods
func (d DocType) PrintSqlModelMethods() (res string) {
	arr := []string{
		fmt.Sprintf("\t\"%s_update\"", d.Name),
		fmt.Sprintf("\t\"%s_list\"", d.Name),
		fmt.Sprintf("\t\"%s_get_by_id\"", d.Name),
	}
	if d.Sql.IsBeforeTrigger {
		arr = append(arr, fmt.Sprintf("\t\"%s_trigger_before\"", d.Name))
	}
	if d.Sql.IsAfterTrigger {
		arr = append(arr, fmt.Sprintf("\t\"%s_trigger_after\"", d.Name))
	}

	if len(arr) > 0 {
		res = fmt.Sprintf("methods = [\n%s\n]", strings.Join(arr, ",\n"))
	}
	return
}

// main.toml печать methods
func (d DocType) PrintSqlModelAlterScripts() (res string) {
	// подбираем postgres тип для alter script
	getType := func(fld FldType) string {
		if fld.Type == "string" {
			if fld.Sql.Size>0 {
				return fmt.Sprintf("CHARACTER VARYING(%v)", fld.Sql.Size)
			} else {
				return "text"
			}
		}
		if fld.Type == FldTypeDouble {
			return "double precision"
		}
		if utils.CheckContainsSliceStr(fld.Type, FldTypeDate, FldTypeDatetime) {
			return "timestamp"
		}
		return fld.Type
	}
	arr := []string{}

	for _, fld := range d.Flds {
		if len(fld.Name) == 0 ||  fld.Sql.IsOptionFld || utils.CheckContainsSliceStr(fld.Name, "id", "created_at", "updated_at", "deleted") {
			continue
		}
		arr = append(arr, fmt.Sprintf("\t\"alter table %s add column if not exists %s %s;\"", d.PgName(), fld.Name, getType(fld)))

	}
	if d.Sql.IsSearchText {
		arr = append(arr, fmt.Sprintf("\t\"alter table %s add column if not exists search_text text;\"", d.PgName()))
	}

	if len(arr) > 0 {
		res = fmt.Sprintf("alterScripts = [\n%s\n]", strings.Join(arr, ",\n"))
	}
	return
}


// get_by_id.sql функиця по добавлению join
func (d DocType) PrintSqlFuncGetById() (res string) {
	cnt := 1
	arr := []string{fmt.Sprintf("with t%v as (select * from %s where id = (params ->> 'id')::int)", cnt, d.PgName())}
	for _, f := range d.Flds {
		if len(f.Sql.Ref) > 0 {
			cnt++
			refTable := f.Sql.Ref
			if refTable == "user" {
				refTable = `"user"`
			}
			arr = append(arr, fmt.Sprintf("\t\tt%v as (select t%[2]v.*, c.title as %[3]s_title from t%[2]v left join %[4]s c on c.id = t%[2]v.%[5]s)", cnt, cnt-1, f.Sql.Ref, refTable, f.Name))
		}
	}
	res = fmt.Sprintf("%s\n \tselect row_to_json(t%v.*)::jsonb into result from t%v;", strings.Join(arr, ",\n"), cnt, cnt)
	return
}

func (d DocType) PrintSqlFuncListWhereCond() string  {
	arr := []string{"['ilike', 'search_text', 'search_text']"}
	for _, fld := range d.Flds {
		if len(fld.Sql.Ref)>0 {
			arr = append(arr, fmt.Sprintf("\t\t['notQuoted', '%[1]s', 'doc.%[1]s']", fld.Name))
		}
	}
	return strings.Join(arr, ",\n")
}

// get_by_id.sql функиця по добавлению join
func (d DocType) PrintSqlFuncList() (res string) {
	cnt := 1
	arr := []string{fmt.Sprintf("EXECUTE ('\n\twith t%v as (select * from %s as doc ' || condQueryStr || ')", cnt, d.PgName())}
	for _, f := range d.Flds {
		if len(f.Sql.Ref) > 0 {
			// в случае если у документа есть флаг IsSearchText, то в options прописывается объект title, в который при обновлении записывается инфа необходимая для поиска
			// если поле отмечено как участвующее в поиске, то инфа по нему записывается в options.title. Соответственно не надо через join находить значение title
			if d.Sql.IsSearchText && f.Sql.IsSearch {
				continue
			}
			cnt++
			refTable := f.Sql.Ref
			if refTable == "user" {
				refTable = `"user"`
			}
			arr = append(arr, fmt.Sprintf("\t\tt%v as (select t%[2]v.*, c.title as %[3]s_title from t%[2]v left join %[4]s c on c.id = t%[2]v.%[5]s)", cnt, cnt-1, f.Sql.Ref, refTable, f.Name))
		}
	}
	res = fmt.Sprintf("%s\n \tselect array_to_json(array_agg(t%v.*)) from t%v') into result;", strings.Join(arr, ",\n"), cnt, cnt)

	return
}

// сборка search_text для list функции
func (d DocType) SearchTxt() string {
	if d.Sql.IsSearchText {
		return "search_text"
	}
	arr := []string{}
	for _, fld := range d.Flds {
		if fld.Sql.IsSearch {
			arr = append(arr, fld.Name)
		}
	}
	if len(arr) == 0 {
		log.Fatalf("missed isSearch fld: %s", d.Name)
	}
	if len(arr) > 1 {
		//concat(doc.description, ' ', doc.supplyer)
		return fmt.Sprintf("concat(%s)", strings.Join(arr, `, '' '', `))
	} else {
		return fmt.Sprintf("doc.%s", arr[0])
	}
}

// update функиця по добавлению execute insert
func (d DocType) PrintSqlFuncInsertNew() (res string) {

	//  индекс поля options для printLinkOnConflict
	optionsFldIndex := 1

	// формирование строчки для update в случае если таблица является связью двух таблиц и эта связь уникальна
	printLinkOnConflict := func() string {
		flds := []FldType{}
		for _, fld := range d.Flds {
			if len(fld.Sql.Ref) > 0 {
				flds = append(flds, fld)
			}
		}
		if len(flds) > 1 {
			return fmt.Sprintf(` ON CONFLICT (%s, %s) DO UPDATE SET options=$%v, deleted=false `, flds[0].Name, flds[1].Name, optionsFldIndex)
		}
		return ""
	}

	arr1 := []string{}
	arr2 := []string{}
	arr3 := []string{}
	for i, f := range d.Flds {
		if len(f.Name) == 0 {
			continue
		}
		if f.Sql.IsOptionFld || utils.CheckContainsSliceStr(f.Name, "id", "created_at", "updated_at", "deleted") {
			continue
		}
		arr1 = append(arr1, f.Name)
		arr2 = append(arr2, fmt.Sprintf("$%v", i+1))
		arrow := "->>"
		if utils.CheckContainsSliceStr(f.Type, "jsonb", FldTypeTextArray)  {
			arrow = "->"
		}
		paramStr := fmt.Sprintf("\t\t\t(params %s '%s')::%s", arrow, f.Name, f.PgInsertType())
		// для text[] своя форма записи
		if f.Type == FldTypeTextArray {
			paramStr = fmt.Sprintf("\t\t\ttext_array_from_json(params %s '%s')", arrow, f.Name)
			//text_array_from_json(params -> 'role')
		}
		arr3 = append(arr3, paramStr)
		// options добавляем последним, поэтому optionsFldIndex увеличиваем на единицу с каждым новым полем, которое будем добавлять
		optionsFldIndex = i+2
	}
	// отдельно добавляем options
	arr1 = append(arr1, "options")
	arr2 = append(arr2, fmt.Sprintf("$%v", optionsFldIndex))
	arr3 = append(arr3, "\t\t\tcoalesce(params -> 'options', '{}')::jsonb")

	res = fmt.Sprintf("EXECUTE ('INSERT INTO %s (%s) VALUES (%s) %s RETURNING *;')", d.Name, strings.Join(arr1, ", "), strings.Join(arr2, ", "), printLinkOnConflict())
	res = fmt.Sprintf("%s\n\t\tINTO %sRow\n\t\tUSING\n%s;", res, d.Name, strings.Join(arr3, ",\n"))
	return
}

// update функиця по добавлению execute insert
func (d DocType) PrintSqlFuncUpdateFlds() (res string) {
	arr := []string{}
	for _, f := range d.Flds {
		if len(f.Name) == 0 {
			continue
		}
		if f.Sql.IsOptionFld || utils.CheckContainsSliceStr(f.Name, "id", "created_at", "updated_at", "deleted") {
			continue
		}
		arr = append(arr, fmt.Sprintf("\t\t\t['%[1]s', '%[1]s', '%[2]s'],", f.Name, f.PgUpdateType()))
	}
	res = strings.Join(arr, "\n")
	return
}

// для BEFORE TRIGGER
// формирование строки из полей для search_txt
func (d DocType) GetSearchTextString() string  {
	arr := []string{}
	for _, fld := range d.Flds {
		if fld.Sql.IsSearch {
			if (len(fld.Sql.Ref) == 0) {
				arr = append(arr, "new." + fld.Name)
			} else {
				arr = append(arr, snaker.SnakeToCamelLower(strings.TrimSuffix(fld.Name, "_id")) + "Title")
			}
		}
	}
	return strings.Join(arr, ", ' ', ")
}

// формирование json из полей для search_txt
func (d DocType) GetSearchTextJson() string  {
	arr := []string{}
	for _, fld := range d.Flds {
		if fld.Sql.IsSearch {
			if (len(fld.Sql.Ref) == 0) {
				arr = append(arr, fmt.Sprintf("'%[1]s', new.%[1]s", fld.Name))
			} else {
				fldName := strings.TrimSuffix(fld.Name, "_id")
				// переменная %sTitle заполняется внутри pg функции. Это title из таблицы, на которую ссылаются
				arr = append(arr, fmt.Sprintf("'%s_title', %sTitle", fldName, snaker.SnakeToCamelLower(fldName)))
			}
		}
	}
	return strings.Join(arr, ", ")
}

// формирование списка переменных для before триггера
func (d DocType) GetBeforeTriggerDeclareVars() string  {
	if !d.Sql.IsSearchText {
		return ""
	}
	res := ""
	for _, fld := range d.Flds {
		if fld.Sql.IsSearch && len(fld.Sql.Ref) > 0 {
			res = fmt.Sprintf("%s\n	%sTitle TEXT;", res, snaker.SnakeToCamelLower(strings.TrimSuffix(fld.Name, "_id")))
		}
	}
	return res
}
func (d DocType) GetBeforeTriggerFillRefVars() string {
	if !d.Sql.IsSearchText {
		return ""
	}
	res := ""
	for _, fld := range d.Flds {
		if fld.Sql.IsSearch && len(fld.Sql.Ref) > 0 {
			refName := fld.Sql.Ref
			if refName == "user" {
				refName = `"user"`
			}
			res = fmt.Sprintf("%s\n		select title into %sTitle from %s where id = new.%s;", res, snaker.SnakeToCamelLower(strings.TrimSuffix(fld.Name, "_id")), refName, fld.Name)
		}
	}
	if len(res) > 0 {
		res = fmt.Sprintf("-- заполняем ref поля%s", res)
	}
	return res
}

func (d DocType) RequiredFldsString() string {
	arr := []string{}
	for _, f := range d.Flds {
		if f.Sql.IsRequired {
			arr = append(arr, fmt.Sprintf("'%s'", f.Name))
		}
	}
	return strings.Join(arr, ", ")
}

// прописываем в модели документа список стандартных sql методов с указанными ролями
func (ds *DocSql) FillBaseMethods(docName string, roles ...string)  {
	if ds.Methods == nil {
		ds.Methods = map[string]*DocSqlMethod{}
	}
	if roles == nil {
		roles = []string{}
	}
	for _, name := range []string{"list", "update", "get_by_id"} {
		name := docName+"_"+name
		ds.Methods[name] = &DocSqlMethod{Name:name, Roles:roles}
	}
}

func (d DocType) PrintAfterTriggerUpdateLinkedRecords() string {
	res := ""
	// ищем таблицы, которые ссылаются на эту и если такие есть, то прописываем триггер, чтобы при обновлении записи, обновляем связанные записи чтобы обновились ссылки
	linkedDocs := [][]string{}
	for _, doc := range project.Docs {
		for _, f := range doc.Flds {
			if f.Sql.Ref == d.Name {
				linkedDocs = append(linkedDocs, []string{doc.Name, f.Name})
			}
		}
	}
	if len(linkedDocs) > 0 {
		// проверка что у документа есть поле title - по нему фиксируем изменения. Если нет, то выдаем ошибку.
		var isFldTitleExist bool
		for _, f := range d.Flds {
			if f.Name == "title" {
				isFldTitleExist = true
				break
			}
		}
		if !isFldTitleExist {
			log.Fatal(fmt.Sprintf("PrintAfterTriggerUpdateLinkedRecords '%s' missed field 'title'", d.Name))
		}
		res1:= "IF (TG_OP = 'UPDATE') THEN\n-- при смене названия обновляем все ссылающиеся записи, чтобы там переписалось новое название\nif new.title != old.title then\n"
		for _, arr := range linkedDocs {
			res1 = fmt.Sprintf("%s for r in select * from %s where %s = new.id loop\n update %s set updated_at=now() where id = r.id;\n end loop;\n end if;\n end if;", res1, arr[0], arr[1], arr[0])
		}
		res = fmt.Sprintf("%s\n%s", res, res1)
	}
	return res
}