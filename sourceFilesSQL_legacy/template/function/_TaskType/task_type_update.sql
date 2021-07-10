-- создание Тип_Задачи
-- параметры:
-- title    char
-- table_name    char
-- options    jsonb

DROP FUNCTION IF EXISTS task_type_update(params JSONB);
CREATE OR REPLACE FUNCTION task_type_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    TaskTypeRow     task_type%ROWTYPE;
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
        checkMsg = check_required_params(params, ARRAY ['title']);
        IF checkMsg IS NOT NULL
        THEN
            RETURN checkMsg;
        END IF;

        EXECUTE ('INSERT INTO task_type (title, table_name, options ) VALUES ($1, $2, $3) RETURNING *;')
            INTO  TaskTypeRow
            USING
             (params ->> 'title')::text,
             (params ->> 'table_name')::text,
             (params -> 'options')::jsonb;
    else
        updateValue = '' || update_str_from_json(params, ARRAY [
            ['title', 'title', 'text'],
            ['table_name', 'table_name', 'text'],
            ['options', 'options', 'jsonb'],
            ['deleted', 'deleted', 'bool']
            ]);

        queryStr = concat('UPDATE task_type SET ', updateValue, ' WHERE id=', params ->> 'id', ' RETURNING *;');

        EXECUTE (queryStr)
            INTO TaskTypeRow;

        -- случай когда записи с таким id не найдено
        IF row_to_json(TaskTypeRow) ->> 'id' ISNULL
        THEN
            RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
        END IF;

    end if;

    result = row_to_json(TaskTypeRow) :: JSONB;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;