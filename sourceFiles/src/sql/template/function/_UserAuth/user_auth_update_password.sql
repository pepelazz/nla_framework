-- обновление пароля пользователя
-- параметры:
-- id             type: int
-- password       type: string
-- auth_provider  type: string - default: email

DROP FUNCTION IF EXISTS user_auth_update_password(params JSONB );
CREATE OR REPLACE FUNCTION user_auth_update_password(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS $function$

DECLARE
    userAuthRow user_auth%ROWTYPE;
    userRow     "user"%ROWTYPE;
    checkMsg    TEXT;
    roleArr     TEXT [] := '{student}';
    result      JSONB;
    authProvider  text := 'email';
BEGIN

    -- проверка наличия обязательных параметров
    -- id в данном случае id user_auth
    checkMsg = check_required_params_with_func_name('user_auth_update_password', params,
                                                    ARRAY ['id', 'password']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    if params->>'auth_provider' notnull then
        authProvider = (params->>'auth_provider')::text;
    end if;

    EXECUTE 'SELECT * FROM user_auth WHERE auth_provider=$1 AND user_id=$2'
        INTO userAuthRow
        USING authProvider, (params ->> 'id') :: INT;

    IF userAuthRow.id ISNULL
    THEN
        RETURN jsonb_build_object('ok', FALSE, 'message', 'wrong user_auth_id');
    END IF;

    UPDATE user_auth
    SET password = (params ->> 'password') :: TEXT
    WHERE id = userAuthRow.id;

    RETURN jsonb_build_object('ok', TRUE, 'result', NULL);

END

$function$;

