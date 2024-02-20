-- поиск пользователя по telegram id
-- параметры:
-- id  type: string

DROP FUNCTION IF EXISTS user_get_by_telegram_id(params JSONB);
CREATE OR REPLACE FUNCTION user_get_by_telegram_id(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    checkMsg  TEXT;
    userRow   "user"%ROWTYPE;
    result    JSONB;
BEGIN

    -- проверка наличия id
    checkMsg = check_required_params_with_func_name('user_get_by_telegram_id', params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    with t1 as (select * from user_auth where auth_provider_id=(params->>'id') and auth_provider='telegram'),
         t2 as (select u.* from t1 left join "user" u on u.id = t1.user_id)
    select * into userRow from t2;

    -- случай когда записи с таким id не найдено
    IF userRow.id ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'user not found');
    END IF;

    result = row_to_json(userRow) :: JSONB;
    -- добавляем auth_token
    result = result || jsonb_build_object('telegram_id', params->>'id');

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;
