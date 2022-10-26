-- создание новой записи пользователе, который должен подтвердить свой  email
-- параметры:
-- email            type: string
-- phone            type: string
-- last_name        type: string
-- first_name       type: string
-- password         type: string
-- token            type: string
-- options          type: json

DROP FUNCTION IF EXISTS user_temp_email_auth_create(params JSONB );
CREATE OR REPLACE FUNCTION user_temp_email_auth_create(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  existUserAuthId INT;
  checkMsg    TEXT;
  authToken    TEXT;
BEGIN

  -- проверка наличия обязательных параметров
  checkMsg = check_required_params(params, ARRAY ['email', 'password', 'token', 'auth_token']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  -- вначале находим все истекшие токены и стираем их
  UPDATE user_temp_email_auth
  SET token = NULL, auth_token = NULL
  WHERE updated_at < (now() - INTERVAL '1 hour');

  -- проверяем что пользователя с таким email'ом в системе нет

  SELECT id INTO existUserAuthId
  FROM user_auth
  WHERE auth_provider = 'email' AND auth_provider_id = (params ->> 'email');

  IF existUserAuthId NOTNULL
  THEN
    RETURN jsonb_build_object('ok', FALSE, 'message', 'email already exist');
  END IF;

  -- генерим токен
  SELECT md5(random() :: TEXT)
  INTO authToken;

  EXECUTE (
    'INSERT INTO user_temp_email_auth (email, last_name, first_name, password, token, auth_token) VALUES ($1, $2, $3, $4, $5, $6) '
    ||
    'ON CONFLICT (email) DO UPDATE SET email=$1, last_name=$2, first_name=$3, password=$4, token=$5, auth_token=$6')
  USING
    params ->> 'email',
    params ->> 'last_name',
    params ->> 'first_name',
    params ->> 'password',
    params ->> 'token',
    authToken;

  RETURN jsonb_build_object('ok', TRUE, 'result', NULL);

END

$function$;

