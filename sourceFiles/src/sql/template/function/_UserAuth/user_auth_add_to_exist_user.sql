-- добавление авторизационного профиля к уже существующему пользователю
-- параметры:
-- user_id          type: int
-- last_name        type: string
-- first_name       type: string
-- avatar           type: string
-- username         type: string
-- auth_provider    type:string
-- auth_provider_id type:string
-- auth_token       type:string
-- email            type:string

DROP FUNCTION IF EXISTS user_auth_add_to_exist_user(params JSONB );
CREATE OR REPLACE FUNCTION user_auth_add_to_exist_user(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  userAuthRow user_auth%ROWTYPE;
  userRow     "user"%ROWTYPE;
  checkMsg    TEXT;
  roleArr     TEXT [] := '{student}';
  result      JSONB;
BEGIN

  -- проверка наличия обязательных параметров
  checkMsg = check_required_params_with_func_name('user_auth_add_to_exist_user', params,
                                                  ARRAY ['user_id', 'auth_provider', 'auth_provider_id', 'auth_token']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  -- проверяем что user с таким id существует
  SELECT *
  INTO userRow
  FROM "user"
  WHERE id = (params ->> 'user_id') :: BIGINT;
  IF userRow.id ISNULL
  THEN
    RETURN jsonb_build_object('ok', FALSE, 'message', concat('not found user with id: ', params ->> 'user_id'));
  END IF;

  -- проверяем есть ли уже такая авторизация
  EXECUTE 'SELECT * FROM user_auth WHERE auth_provider=$1 AND auth_provider_id=$2'
  INTO userAuthRow
  USING params ->> 'auth_provider', params ->> 'auth_provider_id';

  -- если такая авторизация уже есть
  IF userAuthRow.id NOTNULL
  THEN
    -- если она связана с этим же пользователем то ничего не меняем и просто возвращаем данного пользователя
    -- иначе меняем user_id на нового пользователя. Старого пользователя не удаляем
    IF userAuthRow.user_id != userRow.id
    THEN
      UPDATE user_auth
      SET user_id = userRow.id
      WHERE id = userAuthRow.id;
    END IF;

  ELSE
    -- если такой авторизации нет, то создаем
    EXECUTE ('INSERT INTO user_auth (user_id, auth_provider, auth_provider_id, last_name, first_name, username, avatar, auth_token, email, options) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;')
    INTO userAuthRow
    USING
      userRow.id,
      params ->> 'auth_provider',
      params ->> 'auth_provider_id',
      params ->> 'last_name',
      params ->> 'first_name',
      params ->> 'username',
      params ->> 'avatar',
      COALESCE((params ->> 'auth_token') :: TEXT, NULL),
      COALESCE((params ->> 'email') :: TEXT, NULL),
      COALESCE((params -> 'options') :: JSONB, NULL);
  END IF;

  result = (row_to_json(userRow) :: JSONB) || jsonb_build_object('auth_token', userAuthRow.auth_token);
  RETURN jsonb_build_object('ok', TRUE, 'result', result);

END

$function$;

