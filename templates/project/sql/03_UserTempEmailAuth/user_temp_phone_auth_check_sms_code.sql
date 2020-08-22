-- проверка токена пользователя, с помощью которого подтверждаем номер телефона
-- параметры:
-- token            type: string

DROP FUNCTION IF EXISTS user_temp_phone_auth_check_sms_code(params JSONB );
CREATE OR REPLACE FUNCTION user_temp_phone_auth_check_sms_code(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  temp_var user_temp_email_auth%ROWTYPE;
  checkMsg TEXT;
  result   JSONB;
BEGIN

  -- проверка наличия обязательных параметров
  checkMsg = check_required_params(params, ARRAY ['token', 'phone']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  -- вначале находим все истекшие токены и стираем их
  UPDATE user_temp_email_auth
  SET token = NULL, auth_token = NULL
  WHERE phone notnull and updated_at < (now() - INTERVAL '45 second');

  -- находим запись с токеном
  EXECUTE 'SELECT * FROM user_temp_email_auth WHERE token=$1 AND phone=$2'
  INTO temp_var
  USING params ->> 'token', params ->> 'phone';

  IF temp_var ISNULL
  THEN
    RETURN jsonb_build_object('ok', FALSE);
  END IF;

  SELECT *
  FROM user_auth_create(
               jsonb_build_object('auth_provider', 'phone', 'auth_provider_id', temp_var.phone, 'auth_token',
                                  temp_var.auth_token, 'last_name', temp_var.last_name, 'first_name', temp_var.first_name,
                                  'options', coalesce(temp_var.options, '{}'::jsonb) || jsonb_build_object('state', [[if .Config.Auth.IsPassStepWaitingAuth]]'working'[[else]]'waiting_auth'[[end]]),
                                  'password', temp_var.password))
  INTO result;

  -- стираем запись из временной таблицы
  DELETE FROM user_temp_email_auth
  WHERE id = temp_var.id;

  result = result - 'password';

  RETURN result;

END

$function$;

