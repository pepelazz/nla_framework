-- проверка что пользователь админ
-- параметры:
-- user_id   type: int

DROP FUNCTION IF EXISTS user_check_is_admin(user_id BIGINT );
CREATE OR REPLACE FUNCTION user_check_is_admin(user_id BIGINT)
  RETURNS BOOLEAN
LANGUAGE plpgsql
AS $function$

DECLARE
  result BOOLEAN;
BEGIN

  SELECT 'admin' = ANY (role)
  INTO result
  FROM "user"
  WHERE id = user_id;

  RETURN result;

END

$function$;
