-- получение списка Задача
-- параметры:
-- user_id         int - текущий пользователь, длдя которого получаем список задач
-- deleted         type: bool - удаленные / существующие. Дефолт: false
-- order_by        type: string - поле для сортировки и направление сортировки. Например, orderBy: "id desc"
-- page            type: int - номер страницы. Дефолт: 1
-- per_page        type: int - количество записей на странице. Дефолт: 1000
-- search_text     type: string - текстовый поиск

DROP FUNCTION IF EXISTS task_list_for_user(params JSONB);
CREATE OR REPLACE FUNCTION task_list_for_user(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    result       JSON;
    condQueryStr TEXT;
    whereStr     TEXT;
BEGIN

    -- сборка условия WHERE (where_str_build - функция из папки base)
    whereStr = where_str_build(params, 'doc', ARRAY [
        ['ilike', 'search_text', 'concat(task_type_title)'],
        ['text', 'state', 'doc.state']
    ]);

    -- финальная сборка строки с условиями выборки (build_query_part_for_list - функция из папки base)
    condQueryStr = '' || whereStr || concat(' AND executor_id = ', params->>'user_id') || build_query_part_for_list(params);

    EXECUTE (
       'with t1 as (select * from task as doc ' || condQueryStr || '),
        t2 as (select t1.*, u.fullname as executor_fullname from t1 left join "user" u on u.id = t1.executor_id),
        t3 as (select t2.*, tt.options as task_type_options from t2 left join task_type tt on t2.task_type_id = tt.id)
        select array_to_json(array_agg(t3)) from t3')
        INTO result;

    RETURN json_build_object('ok', TRUE, 'result', coalesce(result, '[]'));

END
$function$;




