-- получение списка Задача
-- параметры:
-- deleted         type: bool - удаленные / существующие. Дефолт: false
-- order_by        type: string - поле для сортировки и направление сортировки. Например, orderBy: "id desc"
-- page            type: int - номер страницы. Дефолт: 1
-- per_page        type: int - количество записей на странице. Дефолт: 1000
-- search_text     type: string - текстовый поиск

DROP FUNCTION IF EXISTS task_list(params JSONB);
CREATE OR REPLACE FUNCTION task_list(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    result       JSON;
    condQueryStr TEXT;
    whereStr     TEXT;
    totalRows    INT;
    dateFrom     timestamp;
    dateTo       timestamp;
    dateBetween  text := '';
    tableTitlePath  text := format('t1.table_options ->>%s', quote_literal('title'));
BEGIN

    -- сборка условия WHERE (where_str_build - функция из папки base)
    whereStr = where_str_build(params, 'doc', ARRAY [
        ['ilike', 'search_text', 'concat(task_type_title)'],
        ['text', 'state', 'doc.state']
    ]);

    dateFrom = to_timestamp((params ->> 'date_from'), 'YYYY-MM-DD"T"HH24:MI:SS');
    dateTo = to_timestamp((params ->> 'date_to'), 'YYYY-MM-DD"T"HH24:MI:SS');
    if dateFrom notnull then
        dateBetween = format(' WHERE deadline >= %s', quote_literal(dateFrom));
    end if;
    if dateTo notnull then
        if dateFrom notnull then
            dateBetween = format(' %s AND deadline <= %s', dateBetween, quote_literal(dateTo));
        else
            dateBetween = format(' WHERE deadline <= %s', quote_literal(dateTo));
        end if;
    end if;

    -- для вывода в таблицу нужно считать общее количество записей. Поэтому делаем операцию в два шага.
    -- Сперва сччитаем вссе результаты с учетом фильтрации в отдельную таблицу
    -- потом из этой таблицы забираем нужное количество записей и считаем полное количество записей
    EXECUTE ('CREATE TEMP TABLE IF NOT EXISTS temp_table AS select t.* from (
             with t1 as (select * from task ' || dateBetween || '),
            t2 as (select t1.*, '|| tableTitlePath ||' as table_title, u.fullname as executor_fullname from t1 left join "user" u on u.id = t1.executor_id)
            select * from t2 as doc ' || whereStr || '
             ) t;');

    params = params || jsonb_build_object('prefix', 'doc.');

    -- финальная сборка строки с условиями выборки (build_query_part_for_list - функция из папки base)
    condQueryStr = build_query_part_for_list(params);

    select count(*) into totalRows from temp_table;

    EXECUTE ('with t1 as (select * from temp_table as doc ' || condQueryStr || ')
        select array_to_json(array_agg(t1)) from t1')
        INTO result;

    DROP TABLE temp_table;

    RETURN json_build_object('ok', TRUE, 'result', coalesce(result, '[]'), 'meta_info',
                             jsonb_build_object('total_rows', totalRows));

--
--     -- финальная сборка строки с условиями выборки (build_query_part_for_list - функция из папки base)
--     condQueryStr = '' || whereStr || build_query_part_for_list(params);
--
--     EXECUTE (
--        'with t1 as (select * from task as doc ' || condQueryStr || '),
--         t2 as (select t1.*, u.fullname as executor_fullname from t1 left join "user" u on u.id = t1.executor_id)
--         select array_to_json(array_agg(t2)) from t2')
--         INTO result;
--
--     RETURN json_build_object('ok', TRUE, 'result', coalesce(result, '[]'));

END
$function$;




