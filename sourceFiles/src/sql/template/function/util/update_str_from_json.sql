-- Пример
--  updateValue = '' || update_str_from_json(params, ARRAY [
-- ['infoMsg', 'info_msg', 'text'],
-- ['state', 'state', 'enum']
-- ]);
-- первое значение - поле в json
-- второе значение - поле в postgres
-- третье значение - тип

DROP FUNCTION IF EXISTS update_str_from_json(params JSONB, arr VARCHAR[]);
CREATE OR REPLACE FUNCTION update_str_from_json(params JSONB, arr VARCHAR[])
    RETURNS TEXT
    LANGUAGE plpgsql
AS
$function$
DECLARE
    i             RECORD;
    m             VARCHAR[];
    columnNameStr TEXT := '(';
    valueStr      TEXT := '(';
    cnt           int  := 0;
BEGIN

    FOR i IN SELECT *
             FROM jsonb_each_text(params)

        LOOP
            FOREACH m SLICE 1 IN ARRAY arr
                LOOP
                    IF m[1] = i.key
                    THEN
                        columnNameStr = concat(columnNameStr, concat(m[2], ','));
                        cnt = cnt + 1;
                        CASE m[3]
                            WHEN 'text'
                                THEN valueStr = concat(valueStr, COALESCE(quote_literal(i.value), 'NULL'), ',');
                            WHEN 'enum'
                                THEN valueStr = concat(valueStr, COALESCE(quote_literal(i.value), 'NULL'), ',');
                            WHEN 'jsonb'
                                THEN
                                    IF length(i.value) > 0
                                    THEN
                                        valueStr = concat(valueStr, COALESCE(quote_literal(i.value :: JSONB), 'NULL'),
                                                          ',');
                                    ELSE
                                        valueStr = concat(valueStr, 'NULL', ',');
                                    END IF;
                            WHEN 'number'
                                THEN valueStr = concat(valueStr, COALESCE(NULLIF(trim(i.value), ''), 'NULL'), ',');
                            WHEN 'bool'
                                THEN valueStr = concat(valueStr, COALESCE(i.value, 'NULL'), ',');
                            WHEN 'arrayText'
                                THEN valueStr = concat(valueStr,
                                                       COALESCE(quote_literal(string_to_array(trim(i.value), '|')),
                                                                'NULL'), ',');
                            WHEN 'jsonArrayText'
                                THEN
                                    valueStr = concat(valueStr,
                                                      COALESCE(quote_literal(text_array_from_json(i.value :: JSONB)),
                                                               'NULL'), ',');
                            WHEN 'jsonArrayInt'
                                THEN
                                    valueStr = concat(valueStr,
                                                      COALESCE(quote_literal(int_array_from_json(i.value :: JSONB)),
                                                               'NULL'), ',');
                            WHEN 'timestamp'
                                THEN
                                    valueStr = concat(valueStr, COALESCE(
                                            quote_literal(to_timestamp(i.value, 'YYYY-MM-DD"T"HH24:MI:SS')), 'NULL'),
                                                      ',');
                            WHEN 'time'
                                THEN
                                    valueStr = concat(valueStr, COALESCE(quote_literal(i.value), 'NULL'), ',');
                            ELSE
                                RAISE NOTICE 'update_str_from_json uknown type: %', m[3];
                            END CASE;
                    END IF;
                END LOOP;
        END LOOP;

    columnNameStr = rtrim(columnNameStr, ',');
    columnNameStr = concat(columnNameStr, ')');

    valueStr = rtrim(valueStr, ',');
    valueStr = concat(valueStr, ')');

    -- если обновление только одного значения то убираем скобки
    if cnt = 1 then
        valueStr = replace(valueStr, ')', '');
        valueStr = replace(valueStr, '(', '');
        columnNameStr = replace(columnNameStr, ')', '');
        columnNameStr = replace(columnNameStr, '(', '');
    end if;

    RETURN concat(columnNameStr, ' = ', valueStr);
END ;
$function$;

