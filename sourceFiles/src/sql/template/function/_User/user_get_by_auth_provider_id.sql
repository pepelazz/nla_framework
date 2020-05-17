-- поиск пользователя по auth_provider_id
-- параметры:
-- auth_provider     type: string
-- auth_provider_id  type: string

DROP FUNCTION IF EXISTS user_get_by_auth_provider_id(params JSONB);
CREATE OR REPLACE FUNCTION user_get_by_auth_provider_id(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    checkMsg  TEXT;
    userRow   "user"%ROWTYPE;
    result    JSONB;
    tmpResult JSON;
    queryStr  TEXT;
    userId    BIGINT;
    authToken TEXT;

BEGIN

    -- проверика наличия id
    checkMsg = check_required_params_with_func_name('user_get_by_auth_provider_id', params,
                                                    ARRAY ['auth_provider', 'auth_provider_id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    SELECT user_id,
           auth_token
           INTO userId, authToken
    FROM user_auth
    WHERE auth_provider = params ->> 'auth_provider'
      AND auth_provider_id = params ->> 'auth_provider_id';

    IF userId ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'not found');
    END IF;

    SELECT * INTO userRow
    FROM "user"
    WHERE id = userId;

    -- случай когда записи с таким id не найдено
    IF row_to_json(userRow) ->> 'id' ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'not found');
    END IF;

    result = row_to_json(userRow) :: JSONB;
    -- добавляем auth_token
    result = result || jsonb_build_object('auth_token', authToken);

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;
