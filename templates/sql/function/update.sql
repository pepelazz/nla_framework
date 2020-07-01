{{$doc := . -}}
-- создание {{.NameRu}}

DROP FUNCTION IF EXISTS {{.PgName}}_update(params JSONB);
CREATE OR REPLACE FUNCTION {{.PgName}}_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    {{.Name}}Row     {{.PgName}}%ROWTYPE;
    checkMsg    TEXT;
    result      JSONB;
    updateValue TEXT;
    queryStr    TEXT;
    {{.Sql.Hooks.Print "update" "declareVars"}}
BEGIN

    {{.PrintSqlFuncUpdateCheckParams}}

    {{.Sql.Hooks.Print "update" "beforeInsertUpdate"}}

    {{- range .Flds}}
    {{if eq .Type "uuid" -}}
    if char_length((params->>'{{.Name}}')::text) = 0 then
        params = params || jsonb_build_object('{{.Name}}', '00000000-0000-0000-0000-000000000000');
    end if;
    {{- end}}
    {{- end}}

    {{.PrintSqlFuncUpdateCheckIsNew}}
        {{if .RequiredFldsString -}}
        -- проверика наличия обязательных параметров
        checkMsg = check_required_params(params, ARRAY [{{.RequiredFldsString}}]);
        IF checkMsg IS NOT NULL
        THEN
            RETURN checkMsg;
        END IF;
        {{end -}}

        {{.Sql.Hooks.Print "update" "beforeInsert"}}

        {{.PrintSqlFuncInsertNew}}

    else
        updateValue = '' || update_str_from_json(params, ARRAY [
{{.PrintSqlFuncUpdateFlds}}
            ['options', 'options', 'jsonb'],
            ['deleted', 'deleted', 'bool']
            ]);

        queryStr = {{.PrintSqlFuncUpdateQueryStr}};

        EXECUTE (queryStr)
            INTO {{.Name}}Row;

        -- случай когда записи с таким id не найдено
        IF row_to_json({{.Name}}Row) ->> 'id' ISNULL
        THEN
            RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
        END IF;

    end if;

    {{.Sql.Hooks.Print "update" "afterInsertUpdate"}}

    RETURN {{.PgName}}_get_by_id(jsonb_build_object('id', {{.Name}}Row.id));

END

$function$;