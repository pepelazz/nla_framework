-- поиск Chat по id таблицы
-- параметры:
-- id       type: int

DROP FUNCTION IF EXISTS chat_for_table_id(params JSONB);
CREATE OR REPLACE FUNCTION chat_for_table_id(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    ChatRow         chat%Rowtype;
    checkMsg               TEXT;
    result                 jsonb;
    tableIdTitle text;
    tableId int;
BEGIN

    -- проверика наличия id
    checkMsg = check_required_params_with_func_name('chat_for_table_id', params, ARRAY ['table_name', 'table_id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    with t1 as (select * from chat where table_name=(params->>'table_name') AND table_id=(params ->> 'table_id')::int),
         t2 as (select t1.*,
                       (select array_to_json(array_agg(cm)) from (
                                                                     with tt1 as (select * from chat_message where chat_id = t1.id)
                                                                     select tt1.*, u.fullname, u.avatar  from tt1 left join "user" u on u.id = tt1.user_id) cm) as message_list
                from t1)
    select array_to_json(array_agg(t2)) into result from t2;

    -- случай когда записи с таким id не найдено
    IF result ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'not found');
    END IF;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;