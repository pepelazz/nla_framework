-- получение списка email'ов админов
-- параметры:

DROP FUNCTION IF EXISTS user_get_admin_emails();
CREATE OR REPLACE FUNCTION user_get_admin_emails()
  RETURNS JSONB
LANGUAGE plpgsql
AS $function$

DECLARE
  result JSONB;
BEGIN

  SELECT array_to_json(array_agg(email))
  INTO result
  FROM "user"
  WHERE 'admin' = ANY (role);

  RETURN jsonb_build_object('ok', TRUE, 'result', result);

END

$function$;
