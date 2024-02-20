{{$doc := . -}}
{{$PgName := .Name -}}
-- поиск {{.NameRu}} по id
-- параметры:
-- id       type: int

DROP FUNCTION IF EXISTS {{$PgName}}_get_by_id(params JSONB);
CREATE OR REPLACE FUNCTION {{$PgName}}_get_by_id(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    {{.Name}}Row         {{$PgName}}%Rowtype;
    checkMsg               TEXT;
    result                 jsonb;
BEGIN

    -- проверка наличия id
    checkMsg = check_required_params_with_func_name('{{$PgName}}_get_by_id', params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    {{.PrintSqlFuncGetById}}

    -- случай когда записи с таким id не найдено
    IF result ->> 'id' ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'not found');
    END IF;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;