-- создание пользователя
-- параметры:
-- last_name        type: string
-- first_name       type: string
-- avatar           type: string
-- role             type: []string
-- auth_provider    type:string
-- auth_provider_id type:string

DROP FUNCTION IF EXISTS user_create(params JSONB );
CREATE OR REPLACE FUNCTION user_create(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  temp_var "user"%ROWTYPE;
  checkMsg TEXT;
BEGIN

  -- проверка наличия обязательных параметров
  checkMsg = check_required_params(params, ARRAY ['auth_provider', 'auth_provider_id']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  EXECUTE 'SELECT * FROM "user" WHERE auth_provider=$1 AND auth_provider_id=$2'
  INTO temp_var
  USING params ->> 'auth_provider', params ->> 'auth_provider_id';

  IF temp_var ISNULL
  THEN -- case создания нового пользователя
    EXECUTE ('INSERT INTO "user" (last_name, first_name, avatar, role, auth_provider, auth_provider_id, auth_token, options) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;')
    INTO temp_var
    USING
      params ->> 'last_name',
      params ->> 'first_name',
      params ->> 'avatar',
      text_array_from_json(params -> 'role'),
--       COALESCE(params ->> 'role', 'student'),
      params ->> 'auth_provider',
      params ->> 'auth_provider_id',
      COALESCE((params ->> 'auth_token') :: TEXT, NULL),
      COALESCE((params -> 'options') :: JSONB, NULL);
  END IF;

  RETURN jsonb_build_object('ok', TRUE, 'result', (row_to_json(temp_var) :: JSONB - 'created_at' - 'updated_at'));

END

$function$;

