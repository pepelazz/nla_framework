-- создание новой записи пользователе, который должен подтвердить свой phone
-- параметры:
-- email            type: string
-- phone            type: string
-- last_name        type: string
-- first_name       type: string
-- password         type: string
-- token            type: string
-- options          type: json

DROP FUNCTION IF EXISTS user_temp_phone_auth_create(params JSONB );
CREATE OR REPLACE FUNCTION user_temp_phone_auth_create(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  existUserAuthId INT;
  checkMsg    TEXT;
  authToken    TEXT;
BEGIN

  -- проверка наличия обязательных параметров
  checkMsg = check_required_params(params, ARRAY ['phone', 'password', 'token']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  -- вначале находим все истекшие токены и стираем их
  UPDATE user_temp_email_auth
  SET token = NULL, auth_token = NULL
  WHERE phone notnull and updated_at < (now() - INTERVAL '45 second');

  -- проверяем что пользователя с таким телефоном в системе нет
  SELECT id INTO existUserAuthId
  FROM user_auth
  WHERE auth_provider = 'phone' AND auth_provider_id = (params ->> 'phone');

  IF existUserAuthId NOTNULL
  THEN
    RETURN jsonb_build_object('ok', FALSE, 'message', 'phone already exist');
  END IF;

  -- SqlHooks.CheckIsUserExist
  [[range .Config.Auth.SqlHooks.CheckIsUserExist -]]
  [[.]]
  [[- end]]

  -- генерим токен
  SELECT md5(random() :: TEXT)
  INTO authToken;

  EXECUTE (
    'INSERT INTO user_temp_email_auth (phone, last_name, first_name, password, token, auth_token, options) VALUES ($1, $2, $3, $4, $5, $6, $7) '
    ||
    'ON CONFLICT (phone) DO UPDATE SET phone=$1, last_name=$2, first_name=$3, password=$4, token=$5, auth_token=$6, options=$7')
  USING
    params ->> 'phone',
    params ->> 'last_name',
    params ->> 'first_name',
    params ->> 'password',
    params ->> 'token',
    authToken,
    coalesce(params -> 'options', '{}'::jsonb);

  RETURN jsonb_build_object('ok', TRUE, 'result', NULL);

END

$function$;

