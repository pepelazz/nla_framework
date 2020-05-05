-- поиск по токену
-- параметры:
-- token       type: string
-- user_id     type: string

DROP FUNCTION IF EXISTS file_get_by_token(params JSONB);
CREATE OR REPLACE FUNCTION file_get_by_token(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    FileRow       file%Rowtype;
    checkMsg               TEXT;
    result                 jsonb;
    productParts           jsonb;
BEGIN

    -- проверика наличия id
    checkMsg = check_required_params_with_func_name('file_get_by_token', params, ARRAY ['token']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    with t1 as (select * from file where token = params ->> 'token')
    select row_to_json(t1.*)::jsonb
    into result
    from t1;

    -- случай когда записи с таким id не найдено
    IF result ->> 'id' ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'not found');
    END IF;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;