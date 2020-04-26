-- смена роли пользователя
-- параметры:
-- id     type: int
-- role  type: string

DROP FUNCTION IF EXISTS user_change_role(params JSONB );
CREATE OR REPLACE FUNCTION user_change_role(params JSONB)
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

  -- проверика наличия id
  checkMsg = check_required_params(params, ARRAY ['id']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  roleStr = params ->> 'role';

  IF roleStr ISNULL
  THEN RETURN json_build_object('ok', FALSE, 'message', 'missed "role" value');
  END IF;

  IF roleStr != 'admin' AND roleStr != 'student'
  THEN RETURN json_build_object('ok', FALSE, 'message', 'role must be "admin" or "student"');
  END IF;

  queryStr = concat('UPDATE "user" SET role=', quote_literal(roleStr), ' WHERE id=', params ->> 'id', ' RETURNING *');

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
