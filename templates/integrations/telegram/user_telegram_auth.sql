-- линковка аккаунта в телеграм
-- параметры:
-- username string
-- first_name string
-- last_name string
-- photo_url string
-- id int64
-- user_id int64

DROP FUNCTION IF EXISTS user_telegram_auth(params JSONB);
CREATE OR REPLACE FUNCTION user_telegram_auth(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    checkMsg     TEXT;
    existEmail   TEXT;
    userAuthRow  user_auth%ROWTYPE;
    userAuthRow1 user_auth%ROWTYPE;
    userRow      "user"%ROWTYPE;
    result       JSONB;

BEGIN

    -- проверика наличия id
    checkMsg = check_required_params(params, ARRAY ['id', 'user_id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    select * into userAuthRow from user_auth where auth_provider = 'telegram' AND auth_provider_id = (params ->> 'id');

    if userAuthRow.id notnull then
        return json_build_object('ok', true, 'result', 'user already exist');
    end if;

    insert into user_auth (user_id, auth_provider, auth_provider_id, username, first_name, last_name, avatar)
    values ((params ->> 'user_id')::int,
            'telegram',
            params ->> 'id',
            params ->> 'username',
            params ->> 'first_name',
            params ->> 'last_name',
            params ->> 'photo_url');


    update "user" set options = jsonb_set(options, '{telegram_id}', params->'id') where id = (params->>'user_id')::int;

    RETURN json_build_object('ok', TRUE, 'result', 'ok');

END

$function$;
