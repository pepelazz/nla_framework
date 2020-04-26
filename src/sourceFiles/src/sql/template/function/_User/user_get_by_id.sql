-- поиск пользователя по id
-- параметры:
-- id  type: int

DROP FUNCTION IF EXISTS user_get_by_id(params JSONB );
CREATE OR REPLACE FUNCTION user_get_by_id(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  checkMsg    TEXT;

  userRow    "user"%ROWTYPE;
  result      JSONB;
  queryStr    TEXT;

BEGIN

  -- проверика наличия id
  checkMsg = check_required_params_with_func_name('user_get_by_id', params, ARRAY ['id']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  queryStr = concat('SELECT * FROM "user" WHERE id= ', params ->> 'id');

  EXECUTE (queryStr)
  INTO userRow;

  -- случай когда записи с таким id не найдено
  IF row_to_json(userRow) ->> 'id' ISNULL
  THEN
    RETURN json_build_object('ok', FALSE, 'message', 'not found');
  END IF;

  result = row_to_json(userRow) :: JSONB;

  RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;
