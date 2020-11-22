package types

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/utils"
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
	// check constraints
	for _, cnstr := range d.Sql.CheckConstrains {
		arr = append(arr, fmt.Sprintf("\t{name=\"%s\", ext=\"CHECK(%s)\"}", cnstr.Name, cnstr.CheckConditions))
	}
	// uniq constraints
	for _, cnstr := range d.Sql.UniqConstrains {
		arr = append(arr, fmt.Sprintf("\t{name=\"%s\", ext=\"UNIQUE(%s)\"}", cnstr.Name, cnstr.UniqConditions))
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

// main.toml печать fk_constraints
func (d DocType) PrintSqlModelIndexes() string {
	res := []string{}
	if d.Sql.Indexes != nil {
		for _, v := range d.Sql.Indexes {
			res = append(res, v)
		}
	}
	if len(res) > 0 {
		return fmt.Sprintf("indexes = [\n\t%s\n]", strings.Join(res, ",\n"))
	}
	return ""
}

// main.toml печать methods
func (d DocType) PrintSqlModelMethods() (res string) {
	arr := []string{}
	if d.Sql.Methods != nil {
		for mName := range d.Sql.Methods {
			arr = append(arr, fmt.Sprintf("\t\"%s\"", mName))
		}
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
			if fld.Sql.Size > 0 {
				return fmt.Sprintf("CHARACTER VARYING(%v)", fld.Sql.Size)
			} else {
				return "text"
			}
		}
		if fld.Type == FldTypeInt64 {
			return "int"
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
		if len(fld.Name) == 0 || fld.Sql.IsOptionFld || utils.CheckContainsSliceStr(fld.Name, "id", "created_at", "updated_at", "deleted") {
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
			// формируем имя для title. Нужно для тех случаев когда имя столба отличается от имени таблицы, на которую идет ссылка.
			// например, from_location_id как ссылка на таблицу location. Тогда формируем поле from_location_title
			fldNameWithTitle := strings.TrimSuffix(f.Name, "_id") + "_title"
			arr = append(arr, fmt.Sprintf("\t\tt%v as (select t%[2]v.*, c.title as %[6]s from t%[2]v left join %[4]s c on c.id = t%[2]v.%[5]s)", cnt, cnt-1, f.Sql.Ref, refTable, f.Name, fldNameWithTitle))
		}
	}
	res = fmt.Sprintf("%s\n \tselect row_to_json(t%v.*)::jsonb into result from t%v;", strings.Join(arr, ",\n"), cnt, cnt)
	return
}

func (d DocType) PrintSqlFuncListRoleConditions() string {
	res := ""
	res = fmt.Sprintf(`if is_user_role((params->>'user_id')::int, '{"admin"}') is not true then
        params = params || jsonb_build_object('manager_id', params->>'user_id');
    end if;`, )

	return res
}


func (d DocType) PrintSqlFuncListWhereCond() string {
	arr := []string{"['ilike', 'search_text', 'search_text']"}
	for _, fld := range d.Flds {
		if fld.Name == "title" {
			continue
		}
		if len(fld.Sql.Ref) > 0 {
			arr = append(arr, fmt.Sprintf("\t\t['notQuoted', '%[1]s', 'doc.%[1]s']", fld.Name))
			continue
		}
		if fld.Sql.IsSearch {
			typeStr := "text"
			if fld.Type == FldTypeBool {
				typeStr = "notQuoted"
			}
			arr = append(arr, fmt.Sprintf("\t\t['%[1]s', '%[2]s', 'doc.%[2]s']", typeStr, fld.Name))
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

func (d DocType) PrintSqlFuncUpdateCheckParams() string {
	str := `
    -- проверика наличия id
    checkMsg = check_required_params(params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;
	`
	if d.IsBitrixIntegration() {
		str = fmt.Sprintf(`  
	  checkMsg = check_required_params(params, ARRAY ['btx_id']);
	  IF checkMsg IS NOT NULL
	  THEN
		RETURN checkMsg;
	  END IF;
      -- ищем запись по btx_id, если не находим, значит это новая запись
	  SELECT *
	  INTO %[1]sRow
	  FROM %[1]s
	  WHERE btx_id = (params ->> 'btx_id')::int;
		`, d.Name)
	}
	if d.IsOdataIntegration() {
		str = fmt.Sprintf(`  
	  checkMsg = check_required_params(params, ARRAY ['uuid']);
	  IF checkMsg IS NOT NULL
	  THEN
		RETURN checkMsg;
	  END IF;
      -- ищем запись по uuid, если не находим, значит это новая запись
	  SELECT *
	  INTO %[1]sRow
	  FROM %[1]s
	  WHERE uuid = (params ->> 'uuid')::uuid;
		`, d.Name)
	}
	return str
}

func (d DocType) PrintSqlFuncUpdateCheckIsNew() string {
	str := `if (params ->> 'id')::int = -1 then`
	if (d.IsBitrixIntegration() || d.IsOdataIntegration()) {
		str = fmt.Sprintf(`IF %sRow.id ISNULL THEN`, d.Name)
	}
	return str
}

func (d DocType) PrintSqlFuncUpdateQueryStr() string {
	str := fmt.Sprintf(`concat('UPDATE %s SET ', updateValue, ' WHERE id=', params ->> 'id', ' RETURNING *;')`, d.Name)
	if (d.IsBitrixIntegration()) {
		str = fmt.Sprintf(`concat('UPDATE %[1]s SET ', updateValue, ' WHERE btx_id=', quote_literal(%[1]sRow.btx_id), ' RETURNING *')`, d.Name)
	}
	if (d.IsOdataIntegration()) {
		str = fmt.Sprintf(`concat('UPDATE %[1]s SET ', updateValue, ' WHERE uuid=', quote_literal(%[1]sRow.uuid), ' RETURNING *')`, d.Name)
	}
	return str
}

// update функиця по добавлению execute insert
func (d DocType) PrintSqlFuncInsertNew() (res string) {

	//  индекс поля options для printLinkOnConflict
	optionsFldIndex := 1
	onConflictFldUpdateStr := ""

	// формирование строчки для update в случае если таблица является связью двух таблиц и эта связь уникальна
	printLinkOnConflict := func() string {
		if d.Sql.IsUniqLink {
			flds := []FldType{}
			for _, fld := range d.Flds {
				if len(fld.Sql.Ref) > 0 {
					flds = append(flds, fld)
				}
			}
			if len(flds) > 1 {
				return fmt.Sprintf(` ON CONFLICT (%s, %s) DO UPDATE SET options=$%v, deleted=false%s`, flds[0].Name, flds[1].Name, optionsFldIndex, onConflictFldUpdateStr)
			}
		}
		return ""
	}

	arr1 := []string{}
	arr2 := []string{}
	arr3 := []string{}
	cnt := 0 // счетчик для номеров полей, чтобы корретно выставить номера $1, $2... C учетом того что некоторые поля пропускаются и не вставляются
	for _, f := range d.Flds {
		if len(f.Name) == 0 {
			continue
		}
		if f.Sql.IsOptionFld || utils.CheckContainsSliceStr(f.Name, "id", "created_at", "updated_at", "deleted") {
			continue
		}
		cnt++
		arr1 = append(arr1, f.Name)
		arr2 = append(arr2, fmt.Sprintf("$%v", cnt))
		// если не ref поле то добавляем его в список обновлений при сценарии on conflict
		if len(f.Sql.Ref) == 0 {
			onConflictFldUpdateStr = fmt.Sprintf("%s, %s=$%v", onConflictFldUpdateStr, f.Name, cnt)
		}
		arrow := "->>"
		if utils.CheckContainsSliceStr(f.Type, "jsonb", FldTypeTextArray, FldTypeIntArray) {
			arrow = "->"
		}
		paramStr := fmt.Sprintf("\t\t\t(params %s '%s')::%s", arrow, f.Name, f.PgInsertType())
		// для text[] своя форма записи
		if f.Type == FldTypeTextArray {
			paramStr = fmt.Sprintf("\t\t\ttext_array_from_json(params %s '%s')", arrow, f.Name)
			//text_array_from_json(params -> 'role')
		}
		if f.Type == FldTypeIntArray {
			paramStr = fmt.Sprintf("\t\t\tint_array_from_json(params %s '%s')", arrow, f.Name)
			//int_array_from_json(params -> 'role')
		}
		// в случае наличия дефолтного значения делаем конструкцию coalesce
		if len(f.Sql.Default) > 0 {
			paramStr = fmt.Sprintf("\t\t\tcoalesce((params %s '%s')::%s, %s)::%s", arrow, f.Name, f.PgInsertType(), f.Sql.Default, f.PgInsertType())
		}
		arr3 = append(arr3, paramStr)
		// options добавляем последним, поэтому optionsFldIndex увеличиваем на единицу с каждым новым полем, которое будем добавлять
		optionsFldIndex = cnt + 1
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
		if f.Sql.IsOptionFld || f.Sql.IsNotUpdatable || utils.CheckContainsSliceStr(f.Name, "id", "created_at", "updated_at", "deleted") {
			continue
		}
		arr = append(arr, fmt.Sprintf("\t\t\t['%[1]s', '%[1]s', '%[2]s'],", f.Name, f.PgUpdateType()))
	}
	res = strings.Join(arr, "\n")
	return
}

// для BEFORE TRIGGER
// формирование строки из полей для search_txt
func (d DocType) GetSearchTextString() string {
	arr := []string{}
	for _, fld := range d.Flds {
		if fld.Sql.IsSearch {
			if (len(fld.Sql.Ref) == 0) {
				arr = append(arr, "new."+fld.Name)
			} else {
				arr = append(arr, snaker.SnakeToCamelLower(strings.TrimSuffix(fld.Name, "_id"))+"Title")
			}
		}
	}
	return strings.Join(arr, ", ' ', ")
}

// формирование json из полей для search_txt
func (d DocType) GetSearchTextJson() string {
	arr := []string{}
	for _, fld := range d.Flds {
		if fld.Sql.IsSearch {
			if (len(fld.Sql.Ref) == 0) {
				arr = append(arr, fmt.Sprintf("'%[1]s', new.%[1]s", fld.Name))
			} else {
				fldName := strings.TrimSuffix(fld.Name, "_id")
				// переменная %sTitle заполняется внутри pg функции. Это title из таблицы, на которую ссылаются
				arr = append(arr, fmt.Sprintf("'%s_title', %sTitle", fldName, snaker.SnakeToCamelLower(fldName)))
				// в случае ссылки на user еще добавляем avatar, чтобы потом использовать это на ui
				if fld.Sql.Ref == "user" {
					arr = append(arr, fmt.Sprintf("'%[1]s_avatar', %[1]sAvatar", snaker.SnakeToCamelLower(fldName)))
				}
				for _, v := range fld.Sql.RefFldsForOptions {
					arr = append(arr, fmt.Sprintf("'%s_%s', %s%s", fldName, v, snaker.SnakeToCamelLower(fldName), utils.UpperCaseFirst(v)))
				}
			}
		}
	}
	return strings.Join(arr, ", ")
}

// формирование списка переменных для before триггера
func (d DocType) GetBeforeTriggerDeclareVars() string {
	if !d.Sql.IsSearchText {
		return ""
	}
	res := ""
	for _, fld := range d.Flds {
		if fld.Sql.IsSearch && len(fld.Sql.Ref) > 0 {
			varPrefix := snaker.SnakeToCamelLower(strings.TrimSuffix(fld.Name, "_id"))
			res = fmt.Sprintf("%s\n	%sTitle TEXT;", res, varPrefix)
			if fld.Sql.Ref == "user" {
				res = fmt.Sprintf("%s\n	%sAvatar TEXT;", res, varPrefix)
			}
			for _, v := range fld.Sql.RefFldsForOptions {
				res = fmt.Sprintf("%s\n	%s%s TEXT;", res, varPrefix, utils.UpperCaseFirst(v))
			}
		}
	}
	res = res  + "\n" + d.Sql.Hooks.Print("triggerBefore", "declareVars")

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
			varPrefix := snaker.SnakeToCamelLower(strings.TrimSuffix(fld.Name, "_id"))
			if refName == "user" {
				refName = `"user"`
				// в случае user добавляем еще поле avatar
				res = fmt.Sprintf("%s\n		select title, avatar into %[2]sTitle, %[2]sAvatar from %s where id = new.%s;", res, varPrefix, refName, fld.Name)
			} else {
				fldsArr1 := []string{"title"}
				fldsArr2 := []string{varPrefix + "Title"}
				for _, v := range fld.Sql.RefFldsForOptions {
					fldsArr1 = append(fldsArr1, v)
					fldsArr2 = append(fldsArr2, varPrefix+utils.UpperCaseFirst(v))
				}
				str1 := strings.Join(fldsArr1, ", ")
				str2 := strings.Join(fldsArr2, ", ")
				res = fmt.Sprintf("%s\n		select %s into %s from %s where id = new.%s;", res, str1, str2, refName, fld.Name)
			}
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
func (ds *DocSql) FillBaseMethods(docName string, roles ...string) {
	if ds.Methods == nil {
		ds.Methods = map[string]*DocSqlMethod{}
	}
	if roles == nil {
		roles = []string{}
	}
	for _, name := range []string{"list", "update", "get_by_id"} {
		name := docName + "_" + name
		ds.Methods[name] = &DocSqlMethod{Name: name, Roles: roles}
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
		res1 := "IF (TG_OP = 'UPDATE') THEN\n-- при смене названия обновляем все ссылающиеся записи, чтобы там переписалось новое название\nif new.title != old.title then\n"
		for _, arr := range linkedDocs {
			res1 = fmt.Sprintf("%s for r in select * from %s where %s = new.id loop\n update %s set updated_at=now() where id = r.id;\n end loop;\n", res1, arr[0], arr[1], arr[0])
		}
		res1 = fmt.Sprintf("%s\n end if;\n end if;", res1)
		res = fmt.Sprintf("%s\n%s", res, res1)
	}
	return res
}

func PrintUserAfterTriggerUpdateLinkedRecords() string {
	res := ""
	// ищем таблицы, которые ссылаются на user и если такие есть, то прописываем триггер, чтобы при обновлении записи, обновляем связанные записи чтобы обновились ссылки
	linkedDocs := [][]string{}
	for _, doc := range project.Docs {
		for _, f := range doc.Flds {
			if f.Sql.Ref == "user" {
				linkedDocs = append(linkedDocs, []string{doc.Name, f.Name})
			}
		}
	}
	if len(linkedDocs) > 0 {
		res1 := "IF (TG_OP = 'UPDATE') THEN\n-- при смене имени и аватарки обновляем все ссылающиеся записи, чтобы там переписалось новое название\nif new.fullname != old.fullname OR new.avatar != old.avatar then\n"
		for _, arr := range linkedDocs {
			res1 = fmt.Sprintf("%s for r in select * from %s where %s = new.id loop\n update %s set updated_at=now() where id = r.id;\n end loop;\n", res1, arr[0], arr[1], arr[0])
		}
		res1 = fmt.Sprintf("%s\n end if;\n end if;", res1)
		res = fmt.Sprintf("%s\n%s", res, res1)
	}
	return res
}

// печать sql hook'ов
func (d DocSqlHooks) Print(tmplName, hookName string) string {
	switch hookName {
	case "declareVars":
		if d.DeclareVars != nil {
			// update, triggerBefore
			if r, ok := d.DeclareVars[tmplName]; ok {
				return "-- codogenerated from doc.Sql.Hooks.declareVars\n\t" + r
			}
		}
	case "beforeInsertUpdate":
		if d.BeforeInsertUpdate != nil {
			return strings.Join(d.BeforeInsertUpdate, "\n\n")
		}
	case "beforeInsert":
		if d.BeforeInsert != nil {
			return strings.Join(d.BeforeInsert, "\n\n")
		}
	case "afterInsert":
		if d.AfterInsert != nil {
			return strings.Join(d.AfterInsert, "\n\n")
		}
	case "afterInsertUpdate":
		if d.AfterInsertUpdate != nil {
			return strings.Join(d.AfterInsertUpdate, "\n\n")
		}
	case "BeforeTriggerBefore":
		if d.BeforeTriggerBefore != nil {
			return strings.Join(d.BeforeTriggerBefore, "\n\n")
		}
	case "AfterTriggerAfter":
		if d.AfterTriggerAfter != nil {
			return strings.Join(d.AfterTriggerAfter, "\n\n")
		}
	case "listBeforeBuildWhere":
		if d.ListBeforeBuildWhere != nil {
			return strings.Join(d.ListBeforeBuildWhere, "\n\n")
		}
	case "listAfterBuildWhere":
		if d.ListAfterBuildWhere != nil {
			return strings.Join(d.ListAfterBuildWhere, "\n\n")
		}
	case "afterCreate":
		if d.AfterCreate != nil {
			return strings.Join(d.AfterCreate, "\n\n")
		}
	default:
		return fmt.Sprintf("DocSqlHooks.Print not found code for hook '%s'", hookName)
	}
	return ""
}
