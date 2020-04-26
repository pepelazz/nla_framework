-- сервисы авторизации, которые привязаны к данному пользователю
-- параметры:
-- user_id  type: int

DROP FUNCTION IF EXISTS current_user_get_auth_providers(params JSONB );
CREATE OR REPLACE FUNCTION current_user_get_auth_providers(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  checkMsg TEXT;
  result   JSONB;

BEGIN

  -- проверика наличия id
  checkMsg = check_required_params_with_func_name('user_get_by_auth_token', params, ARRAY ['user_id']);
  IF checkMsg IS NOT NULL
  THEN
    RETURN checkMsg;
  END IF;

  SELECT array_to_json(array_agg(t))
  INTO result
  FROM (SELECT auth_provider
        FROM user_auth
        WHERE user_id = (params ->> 'user_id') :: BIGINT) t;

  RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;
