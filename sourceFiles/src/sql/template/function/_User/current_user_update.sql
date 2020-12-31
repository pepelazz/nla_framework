-- обновление данных текущего пользователя
-- параметры:
-- user_id     type: int
-- first_name  type: string
-- last_name   type: string
-- avatar      type: string
-- options     type: json

DROP FUNCTION IF EXISTS current_user_update(params JSONB );
CREATE OR REPLACE FUNCTION current_user_update(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE

  userRow     "user"%ROWTYPE;
  result       JSONB;
  existOptions JSONB;
  updateValue  TEXT;
  queryStr     TEXT;
  checkMsg     TEXT;

BEGIN

  -- проверика наличия id
  checkMsg = check_required_params(params, ARRAY ['user_id']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  SELECT options
  INTO existOptions
  FROM "user"
  WHERE id = (params ->> 'user_id') :: BIGINT;

  -- если существующией options не null, то объединяем их с пришедшими параметрами, а уже потом сохраняем в базу. Иначе новые параметры просто перезапишут старые
  IF existOptions NOTNULL AND
     (params -> 'options') NOTNULL -- проверка что options есть в переданных параметрах на обновление
  THEN
    IF (params ->> 'options') NOTNULL -- если не options:null, то объекдиням с существующими options
    THEN
      params = params || jsonb_build_object('options', existOptions || (params -> 'options'));
    END IF;
  END IF;

  if params->>'phone' notnull then
      params = params || jsonb_build_object('phone', phone_change_8_to_7((params->>'phone')::text));
  end if;

  updateValue = '' || update_str_from_json(params, ARRAY [
  ['last_name', 'last_name', 'text'],
  ['first_name', 'first_name', 'text'],
  ['phone', 'phone', 'text'],
  ['avatar', 'avatar', 'text'],
  ['options', 'options', 'jsonb']
  ]);

  queryStr = concat('UPDATE "user" SET ', updateValue, ' WHERE id=', params ->> 'user_id', ' RETURNING *');

  EXECUTE (queryStr)
  INTO userRow;

  -- случай когда записи с таким id не найдено
  IF row_to_json(userRow) ->> 'id' ISNULL
  THEN
    RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
  END IF;

  result = row_to_json(userRow) :: JSONB;

  RETURN json_build_object('ok', TRUE, 'result', result - 'password');

END

$function$;
