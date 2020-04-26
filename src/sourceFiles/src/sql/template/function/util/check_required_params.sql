-- проверка наличия обязательных параметров
-- передается входящий json с параметрами и массив полей, для которых осуществляется проверка

DROP FUNCTION IF EXISTS check_required_params(params JSONB, fldNames TEXT [] );
CREATE OR REPLACE FUNCTION check_required_params(params JSONB, fldNames TEXT [])
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  fld TEXT;
BEGIN

  FOREACH fld IN ARRAY fldNames
  LOOP
    IF (params ->> fld) ISNULL
    THEN
      RETURN json_build_object('ok', FALSE, 'message', concat('missing prop: ', fld));
    END IF;
  END LOOP;

  RETURN NULL;

END

$function$;


DROP FUNCTION IF EXISTS check_required_params_with_func_name(funcName TEXT, params JSONB, fldNames TEXT [] );
CREATE OR REPLACE FUNCTION check_required_params_with_func_name(funcName TEXT, params JSONB, fldNames TEXT [])
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE
  fld TEXT;
BEGIN

  FOREACH fld IN ARRAY fldNames
  LOOP
    IF (params ->> fld) ISNULL
    THEN
      RETURN json_build_object('ok', FALSE, 'message', concat(funcName, ' missing prop: ', fld));
    END IF;
  END LOOP;

  RETURN NULL;

END

$function$;
