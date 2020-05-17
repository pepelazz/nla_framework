-- поиск пользователя по токену
-- параметры:
-- auth_token  type: string

DROP FUNCTION IF EXISTS user_get_by_auth_token(params JSONB);
CREATE OR REPLACE FUNCTION user_get_by_auth_token(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    checkMsg  TEXT;
    userRow   "user"%ROWTYPE;
    result    JSONB;
    queryStr  TEXT;
    userId    BIGINT;
    authToken TEXT;

BEGIN

    -- проверика наличия id
    checkMsg = check_required_params_with_func_name('user_get_by_auth_token', params, ARRAY ['auth_token']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    SELECT user_id,
           auth_token
           INTO userId, authToken
    FROM user_auth
    WHERE auth_token = params ->> 'auth_token';

    IF userId ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'invalid token');
    END IF;

    SELECT * INTO userRow
    FROM "user"
    WHERE id = userId;

    -- случай когда записи с таким id не найдено
    IF userRow.id ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'invalid token');
    END IF;

    result = row_to_json(userRow) :: JSONB;
    -- добавляем auth_token
    result = result || jsonb_build_object('auth_token', authToken);

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;
