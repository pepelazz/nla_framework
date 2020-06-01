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

    result = row_to_json({{.Name}}Row) :: JSONB;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;