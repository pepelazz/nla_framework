-- создание Задача
-- параметры:
-- content      text
-- type         text
-- table_name   text - заполняются триггером
-- task_type_title     text - заполняются триггером
-- table_id     text
-- executor_id  int
-- manager_id   int
-- state        text
-- deadline    timestamp
-- date_completed    timestamp
-- result       text
-- success_rate int
-- options      jsonb

DROP FUNCTION IF EXISTS task_update(params JSONB);
CREATE OR REPLACE FUNCTION task_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    TaskRow     task%ROWTYPE;
    TaskType    task_type%ROWTYPE;
    checkMsg    TEXT;
    result      JSONB;
    updateValue TEXT;
    queryStr    TEXT;
BEGIN

    -- проверика наличия id
    checkMsg = check_required_params(params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    if (params ->> 'id')::int = -1 then
        -- проверика наличия обязательных параметров
        checkMsg = check_required_params(params, ARRAY ['task_type_id']);
        IF checkMsg IS NOT NULL
        THEN
            RETURN checkMsg;
        END IF;

        -- если задача прикрепляется к документу, то проверка что есть table_id
        select * into TaskType from task_type where id = (params->>'task_type_id')::int;
        if TaskType.table_name notnull AND (params->>'table_id'):: int isnull then
            return jsonb_build_object('ok', false, 'message', 'missed table_id');
        end if;

        if params->>'manager_id' isnull then
            params = params || jsonb_build_object('manager_id', params->>'user_id');
        end if;
        if params->>'executor_id' isnull then
            params = params || jsonb_build_object('executor_id', params->>'user_id');
        end if;

        EXECUTE ('INSERT INTO task (content, task_type_id, table_id, executor_id, manager_id, state, deadline, date_completed, result, success_rate, options ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;')
            INTO  TaskRow
            USING
             (params ->> 'content')::text,
             (params ->> 'task_type_id')::int,
             (params ->> 'table_id')::int,
             (params ->> 'executor_id')::int,
             (params ->> 'manager_id')::int,
             coalesce((params ->> 'state'), 'in_process')::text,
             (params ->> 'deadline')::timestamp,
             (params ->> 'date_completed')::timestamp,
             (params ->> 'result')::text,
             (params ->> 'success_rate')::int,
             (params -> 'options')::jsonb;
    else
        updateValue = '' || update_str_from_json(params, ARRAY [
            ['content', 'content', 'text'],
--             ['task_type_id', 'task_type_id', 'number'], - запрет на изменения
--             ['table_id', 'table_id', 'number'],
            ['executor_id', 'executor_id', 'number'],
            ['manager_id', 'manager_id', 'number'],
            ['state', 'state', 'text'],
            ['deadline', 'deadline', 'timestamp'],
            ['date_completed', 'date_completed', 'timestamp'],
            ['result', 'result', 'text'],
            ['success_rate', 'success_rate', 'number'],
            ['options', 'options', 'jsonb'],
            ['deleted', 'deleted', 'bool']
            ]);

        queryStr = concat('UPDATE task SET ', updateValue, ' WHERE id=', params ->> 'id', ' RETURNING *;');

        EXECUTE (queryStr)
            INTO TaskRow;

        -- случай когда записи с таким id не найдено
        IF row_to_json(TaskRow) ->> 'id' ISNULL
        THEN
            RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
        END IF;

    end if;

    result = row_to_json(TaskRow) :: JSONB;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;