-- обновление пользователя
-- параметры:
-- first_name  type: string
-- last_name   type: string
-- role        type: string   - роль пользователя
-- avatar      type: string
-- deleted     type: bool

DROP FUNCTION IF EXISTS user_update(params JSONB);
CREATE OR REPLACE FUNCTION user_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE

    temp_var    "user"%ROWTYPE;
    result      JSONB;
    updateValue TEXT;
    queryStr    TEXT;
    checkMsg    TEXT;

BEGIN

    -- проверика наличия id
    checkMsg = check_required_params(params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    if params->>'phone' notnull then
        params = params || jsonb_build_object('phone', phone_change_8_to_7((params->>'phone')::text));
    end if;

    updateValue = '' || update_str_from_json(params, ARRAY [
        ['last_name', 'last_name', 'text'],
        ['first_name', 'first_name', 'text'],
        ['role', 'role', 'jsonArrayText'],
        ['avatar', 'avatar', 'text'],
        ['phone', 'phone', 'text'],
        ['grade', 'grade', 'text'],
        ['options', 'options', 'jsonb'],
        ['deleted', 'deleted', 'bool']
        ]);

    queryStr = concat('UPDATE "user" SET ', updateValue, ' WHERE id=', params ->> 'id', ' RETURNING *');

    raise notice 'queryStr %', queryStr;

    EXECUTE (queryStr)
        INTO temp_var;

    -- случай когда записи с таким id не найдено
    IF row_to_json(temp_var) ->> 'id' ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
    END IF;

    result = row_to_json(temp_var) :: JSONB;

    RETURN json_build_object('ok', TRUE, 'result', result - 'password');

END

$function$;
