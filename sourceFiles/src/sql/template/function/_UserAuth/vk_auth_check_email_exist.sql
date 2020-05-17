-- проверка что у пользователя, авторизованного чере вк, заполнено поле email
-- параметры:
-- auth_provider_id     type: string

DROP FUNCTION IF EXISTS vk_auth_check_email_exist(params JSONB );
CREATE OR REPLACE FUNCTION vk_auth_check_email_exist(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  checkMsg     TEXT;
  emailExist   TEXT;
  userAuthRow  user_auth%ROWTYPE;
  userAuthRow1 user_auth%ROWTYPE;
  userRow      "user"%ROWTYPE;
  result       JSONB;

BEGIN

  -- проверика наличия id
  checkMsg = check_required_params(params, ARRAY ['auth_provider_id']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  SELECT email
  INTO emailExist
  FROM user_auth
  WHERE auth_provider_id = params ->> 'auth_provider_id' AND auth_provider = 'vk';

  IF emailExist NOTNULL AND length(emailExist) > 0
  THEN
    RETURN json_build_object('ok', TRUE, 'result', jsonb_build_object('email', emailExist));
  ELSE RETURN json_build_object('ok', FALSE, 'message', 'email is empty');
  END IF;

END

$function$;
