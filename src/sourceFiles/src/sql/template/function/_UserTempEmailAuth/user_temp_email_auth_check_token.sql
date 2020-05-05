-- проверка токена пользователя, с помощью которого подтверждаем email
-- параметры:
-- token            type: string

DROP FUNCTION IF EXISTS user_temp_email_auth_check_token(params JSONB );
CREATE OR REPLACE FUNCTION user_temp_email_auth_check_token(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  temp_var user_temp_email_auth%ROWTYPE;
  checkMsg TEXT;
  result   JSONB;
BEGIN

  -- проверка наличия обязательных параметров
  checkMsg = check_required_params(params, ARRAY ['token']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  -- находим запись с токеном
  EXECUTE 'SELECT * FROM user_temp_email_auth WHERE token=$1'
  INTO temp_var
  USING params ->> 'token';

  IF temp_var ISNULL
  THEN
    RETURN jsonb_build_object('ok', FALSE);
  END IF;

  SELECT *
  FROM user_auth_create(
               jsonb_build_object('auth_provider', 'email', 'auth_provider_id', temp_var.email, 'auth_token',
                                  temp_var.auth_token, 'last_name', temp_var.last_name, 'first_name', temp_var.first_name,
                                  'username', temp_var.email, 'email', temp_var.email, 'options', jsonb_build_object('state', 'waiting_auth'),
                                  'password', temp_var.password))
  INTO result;

  -- стираем запись из временной таблицы
  DELETE FROM user_temp_email_auth
  WHERE id = temp_var.id;

  result = result - 'password';

  RETURN result;

END

$function$;

