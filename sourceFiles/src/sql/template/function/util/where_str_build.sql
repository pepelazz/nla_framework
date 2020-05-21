-- Пример
-- whereStr = where_str_build(params, ARRAY[
--     ['enum', 'state', 'q.state'],
--     ['notQuoted', 'surveyId', 'q.survey_id']
--   ])
-- tableAlias - буква для названия таблицы для которой определяем свойство delete

DROP FUNCTION IF EXISTS where_str_build(params JSONB, tableAlias VARCHAR, arr VARCHAR[]);
CREATE OR REPLACE FUNCTION where_str_build(params JSONB, tableAlias VARCHAR, arr VARCHAR[])
    RETURNS TEXT
    LANGUAGE plpgsql
AS
$function$
DECLARE
    m        VARCHAR[];
    whereStr TEXT := concat(' where ', tableAlias, '.deleted=', COALESCE((params ->> 'deleted'), 'false'));
BEGIN

    FOREACH m SLICE 1 IN ARRAY arr
        LOOP

            -- ENUM
            IF m[1] = 'enum'
            THEN
                IF (params ->> m[2]) IS NOT NULL AND (params ->> m[2]) != 'all'
                THEN
                    whereStr = concat(whereStr, concat(' AND ', m[3], '='), quote_nullable(params ->> m[2]));
                END IF;
            END IF;

            -- ЗНАЧЕНИЕ В КОВЫЧКАХ
            IF m[1] = 'text'
            THEN
                IF (params ->> m[2]) IS NOT NULL
                THEN
                    whereStr = concat(whereStr, concat(' AND ', m[3], '='), quote_literal(params ->> m[2]));
                END IF;
            END IF;

            -- ЗНАЧЕНИЕ БЕЗ КОВЫЧЕК
            IF m[1] = 'notQuoted'
            THEN
                IF (params ->> m[2]) IS NOT NULL
                THEN
                    if (params ->> m[2]) = 'null' then
                        whereStr = concat(whereStr, concat(' AND ', m[3], ' is null '));
                    else
                        whereStr = concat(whereStr, concat(' AND ', m[3], '='), params ->> m[2]);
                    end if;
                END IF;
            END IF;

            -- ПОИСК ПО ТЕКСТУ
            IF m[1] = 'ilike'
            THEN
                IF (params ->> m[2]) IS NOT NULL
                THEN
                    whereStr = concat(whereStr, concat(' AND ', m[3], ' ilike '),
                                      quote_literal(concat('%', (params ->> m[2]), '%')));
                END IF;
            END IF;

            -- ПОИСК ПО json МАССИВУ
            IF m[1] = 'jsonArrayText'
            THEN
                IF (params ->> m[2]) IS NOT NULL
                THEN
                    -- проверка что параметр является массивом
                    BEGIN
                        PERFORM text_array_from_json((params -> m[2]) :: JSONB);
                    EXCEPTION
                        WHEN OTHERS
                            THEN
                                RAISE EXCEPTION 'params "%" must be array', m[2];
                    END;
                    whereStr = concat(whereStr, concat(' AND ', m[3], ' @> ',
                                                       quote_literal(text_array_from_json((params -> m[2]) :: JSONB))));
                END IF;
            END IF;

            -- FULL TEXT SEARCH
            IF m[1] = 'fts'
            THEN
                IF (params ->> m[2]) IS NOT NULL
                THEN
                    whereStr = concat(whereStr, concat(' AND ', m[3], ' @@ '),
                                      quote_literal(replace(trim((params ->> m[2])), ' ', '&')), ':: tsquery');
                END IF;
            END IF;

        END LOOP;

    RETURN whereStr;
END;
$function$;

