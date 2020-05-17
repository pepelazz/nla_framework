-- смена email по id auth_provider_id
-- параметры:
-- auth_provider     type: string
-- auth_provider_id     type: string
-- email  type: string

DROP FUNCTION IF EXISTS user_auth_set_email_by_auth_provider_id(params JSONB );
CREATE OR REPLACE FUNCTION user_auth_set_email_by_auth_provider_id(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  checkMsg     TEXT;
  existEmail   TEXT;
  userAuthRow  user_auth%ROWTYPE;
  userAuthRow1 user_auth%ROWTYPE;
  userRow      "user"%ROWTYPE;
  result       JSONB;

BEGIN

  -- проверика наличия id
  checkMsg = check_required_params(params, ARRAY ['auth_provider', 'auth_provider_id', 'email']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  IF char_length(params->>'email') < 5 THEN
    RETURN json_build_object('ok', FALSE, 'message', 'wrong email');
  END IF;

  UPDATE user_auth
  SET email = params ->> 'email'
  WHERE auth_provider = params ->> 'auth_provider' AND
        auth_provider_id = params ->> 'auth_provider_id'
  RETURNING *
    INTO userAuthRow;

  -- случай когда записи с таким id не найдено
  IF row_to_json(userAuthRow) ->> 'id' ISNULL
  THEN
    RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
  END IF;
  -- проверяем что если есть user_auth с таким email, то проставляем связь с пользователем, который связан с ним. Если нет и у основного пользователя не проставлен email, то заполняем поле
  SELECT *
  INTO userAuthRow1
  FROM user_auth
  WHERE id != userAuthRow.id AND email = params ->> 'email';
  IF userAuthRow1.id NOTNULL
  THEN
    UPDATE user_auth
    SET user_id = userAuthRow1.user_id
    WHERE id = userAuthRow.id;
  END IF;

  -- проверяем что если у пользователя не заполнено поле email, то заполняем его
  SELECT email
  INTO existEmail
  FROM "user"
  WHERE id = userAuthRow.user_id;
  IF existEmail ISNULL OR length(existEmail) = 0
  THEN
    UPDATE "user"
    SET email = userAuthRow.email
    WHERE id = userAuthRow.user_id;
  END IF;

  result = row_to_json(userAuthRow) :: JSONB;

  RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;
