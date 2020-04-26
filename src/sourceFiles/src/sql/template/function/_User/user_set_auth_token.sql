-- смена токена пользователя
-- параметры:
-- user_id     type: int
-- auth_token  type: string

DROP FUNCTION IF EXISTS user_set_auth_token(params JSONB );
CREATE OR REPLACE FUNCTION user_set_auth_token(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  checkMsg    TEXT;
  roleStr     TEXT;

  temp_var    "user"%ROWTYPE;
  result      JSONB;
  updateValue TEXT;
  queryStr    TEXT;

BEGIN

  -- проверка наличия id
  checkMsg = check_required_params_with_func_name('user_set_auth_token', params, ARRAY ['user_id', 'auth_token']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  queryStr = concat('UPDATE "user" SET auth_token=', quote_literal(params ->> 'auth_token'), ' WHERE id=', params ->> 'user_id', ' RETURNING *');

  EXECUTE (queryStr)
  INTO temp_var;

  -- случай когда записи с таким id не найдено
  IF row_to_json(temp_var) ->> 'id' ISNULL
  THEN
    RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
  END IF;

  result = row_to_json(temp_var) :: JSONB;

  RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;
