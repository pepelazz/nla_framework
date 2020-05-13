-- получение списка тэгов
-- параметры:

DROP FUNCTION IF EXISTS {{.PgName}}_{{GetFld}}_list(params JSONB );
CREATE OR REPLACE FUNCTION {{.PgName}}_{{GetFld}}_list(params JSONB)
  RETURNS JSONB
LANGUAGE plpgsql
AS $function$

DECLARE
  result JSON;
BEGIN

  EXECUTE (
    'SELECT array_to_json(array_agg(t.unnest)) FROM (select DISTINCT unnest({{GetFld}}) from {{.PgName}}) AS t')
  INTO result;

  RETURN json_build_object('ok', TRUE, 'result', coalesce(result, '[]'));

END
$function$;