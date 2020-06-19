-- создание Задачи по имени task_type
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

DROP FUNCTION IF EXISTS task_create_by_task_type_title(params JSONB);
CREATE OR REPLACE FUNCTION task_create_by_task_type_title(params JSONB)
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
    taskTypeId int;
BEGIN

    -- проверика наличия id
    checkMsg = check_required_params(params, ARRAY ['task_type_title']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    select id into taskTypeId from task_type where title_en = (params->>'task_type_title')::text;
    if taskTypeId isnull then
        return jsonb_build_object('ok', false, 'message', format('not found task_type with title_en "%s"', (params->>'task_type_title')));
    end if;

    -- если задача прикрепляется к документу, то проверка что есть table_id
    select * into TaskType from task_type where id = taskTypeId;
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
        taskTypeId,
         (params ->> 'table_id')::int,
         (params ->> 'executor_id')::int,
         (params ->> 'manager_id')::int,
         coalesce((params ->> 'state'), 'in_process')::text,
         (params ->> 'deadline')::timestamp,
         (params ->> 'date_completed')::timestamp,
         (params ->> 'result')::text,
         (params ->> 'success_rate')::int,
         (params -> 'options')::jsonb;

    result = row_to_json(TaskRow) :: JSONB;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;