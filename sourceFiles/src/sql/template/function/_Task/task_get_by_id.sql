-- поиск Задача по id
-- параметры:
-- id       type: int

DROP FUNCTION IF EXISTS task_get_by_id(params JSONB);
CREATE OR REPLACE FUNCTION task_get_by_id(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    TaskRow         task%Rowtype;
    checkMsg               TEXT;
    result                 jsonb;
    tableIdTitle text;
    tableId int;
BEGIN

    -- проверика наличия id
    checkMsg = check_required_params_with_func_name('task_get_by_id', params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    with t1 as (select * from task where id = (params ->> 'id')::int),
    t2 as (select t1.*, tt.title as task_type_title from t1 left join task_type tt on tt.id = t1.task_type_id),
    t3 as (select t2.*, u.fullname as executor_fullname from t2 left join "user" u on u.id = t2.executor_id)
    select row_to_json(t3.*)::jsonb
    into result
    from t3;

    -- случай когда записи с таким id не найдено
    IF result ->> 'id' ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'not found');
    END IF;

    -- добавляем поле table_id_title - для селектора на интерфейсе
    if result->>'table_id' notnull then
        tableId = (result->>'table_id')::int;
        if (result->>'table_name')::text = 'client' then
            select title into tableIdTitle from client where id = tableId;
        end if;
    end if;

    result = result || jsonb_build_object('table_id_title', tableIdTitle);

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;