-- поиск пользователя по email и паролю
-- параметры:
-- email     type: string

DROP FUNCTION IF EXISTS user_get_by_email_with_password(params JSONB);
CREATE OR REPLACE FUNCTION user_get_by_email_with_password(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    checkMsg  TEXT;
    userRow  "user"%ROWTYPE;
    temp_var  user_auth%ROWTYPE;
    result    JSONB;
    tmpResult JSON;
    queryStr  TEXT;

BEGIN

    -- проверика наличия id
    checkMsg = check_required_params_with_func_name('user_get_by_email_with_password', params, ARRAY ['email']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    -- проверяем что пользователь с таким email есть
    EXECUTE ('SELECT * FROM user_auth WHERE auth_provider=$1 AND auth_provider_id=$2')
        INTO temp_var
        USING 'email', params ->> 'email';

    -- случай когда записи с таким email не найдено
    IF row_to_json(temp_var) ->> 'id' ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'user not found');
    END IF;

    select * into userRow from "user" where id = temp_var.user_id;

    result = row_to_json(userRow) :: JSONB || jsonb_build_object('password', temp_var.password, 'auth_token', temp_var.auth_token);

    -- из итогового результата убираем пароль
    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;
