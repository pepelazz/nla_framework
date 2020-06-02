CREATE OR REPLACE FUNCTION random(NUMERIC, NUMERIC)
  RETURNS NUMERIC AS
$$
SELECT ($1 + ($2 - $1) * random()) :: NUMERIC;
$$ LANGUAGE 'sql' VOLATILE;

-- функция конвертации json массива в текстовый массив
DROP FUNCTION IF EXISTS text_array_from_json(jsonArr JSONB );
CREATE OR REPLACE FUNCTION text_array_from_json(jsonArr JSONB)
  RETURNS TEXT []
LANGUAGE plpgsql
AS $function$
BEGIN

  IF jsonArr ISNULL OR jsonArr = 'null'
  THEN RETURN NULL;
  END IF;

  RETURN COALESCE((SELECT array_agg(e)
                   FROM jsonb_array_elements_text(jsonArr) e), '{}' :: TEXT []);
END
$function$;

-- функция конвертации json массива в текстовый массив
DROP FUNCTION IF EXISTS int_array_from_json(jsonArr JSONB );
CREATE OR REPLACE FUNCTION int_array_from_json(jsonArr JSONB)
  RETURNS INT []
LANGUAGE plpgsql
AS $function$
BEGIN

  IF jsonArr ISNULL OR jsonArr = 'null'
  THEN RETURN NULL;
  END IF;

  RETURN COALESCE((SELECT array_agg(e) :: INT []
                   FROM jsonb_array_elements_text(jsonArr) e), '{}' :: INT []);
END
$function$;

-- функция для модификации options - используется в функции first_raw_transition_order_update
DROP FUNCTION IF EXISTS options_add_fld(userId int, options JSONB, fldName text, jsonObj jsonb);
CREATE OR REPLACE FUNCTION options_add_fld(userId int, options JSONB, fldName text, jsonObj jsonb)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
BEGIN
    return jsonb_set(options, string_to_array(fldName, ''), coalesce(options -> fldName, '[]'::jsonb) ||
                                                            (jsonObj || jsonb_build_object('user_id', userId, 'date', now() at time zone '[[Config.Postgres.TimeZone]]')));
END
$function$;

-- количество дней, которое надо прибавить чтобы получить следующий рабочий день
DROP FUNCTION IF EXISTS next_business_day(timestamp);
CREATE OR REPLACE FUNCTION next_business_day(timestamp)
    RETURNS interval
    LANGUAGE plpgsql
AS
$function$
DECLARE
    weekday integer;
BEGIN
    weekday := extract(dow from $1);
    IF weekday = 0 THEN
        return format('%s days', 2);
    ELSIF weekday = 6 THEN
        return format('%s days', 3);
    ELSE
        return format('%s days', 1);
    END IF;
END;
$function$;




