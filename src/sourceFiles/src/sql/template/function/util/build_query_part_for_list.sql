-- построение части строки запроса списка документов
-- параметры:
-- order_by      type: string - поле для сортировки и направление сортировки.
-- page_num      type: int - номер страницы. Дефолт: 1
-- per_page      type: int - количество записей на странице. Дефолт: 10

DROP FUNCTION IF EXISTS build_query_part_for_list(params JSONB);
CREATE OR REPLACE FUNCTION build_query_part_for_list(params JSONB)
    RETURNS TEXT
    LANGUAGE plpgsql
AS
$function$

DECLARE

    orderBy  TEXT;
    limitNum TEXT;
    pageNum  INT;
    perPage  INT := COALESCE((params ->> 'per_page') :: INT, 1000);
    page     INT := COALESCE((params ->> 'page') :: INT, 1);

BEGIN

    -- сборка сортировки
    IF (params ->> 'order_by') IS NOT NULL
    THEN
        -- вариант когда, например, doc.order_by
        if params ->> 'prefix' is not null then
            orderBy = concat(' ORDER BY ', (params ->> 'prefix'), (params ->> 'order_by'));
        else
            orderBy = concat(' ORDER BY ', (params ->> 'order_by'));
        end if;
    END IF;

    -- сборка pagination
    limitNum = concat(' LIMIT ', perPage, ' OFFSET ', COALESCE((page - 1) * perPage, 0));

    RETURN '' || COALESCE(orderBy, '') || limitNum;

END

$function$;
