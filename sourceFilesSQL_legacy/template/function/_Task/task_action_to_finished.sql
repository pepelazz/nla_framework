-- создание Задача
-- параметры:
-- id            int
-- result        text
-- success_rate  int

DROP FUNCTION IF EXISTS task_action_to_finished(params JSONB);
CREATE OR REPLACE FUNCTION task_action_to_finished(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    TaskRow     task%ROWTYPE;
    checkMsg    TEXT;
    result      JSONB;
    updateValue TEXT;
    queryStr    TEXT;
BEGIN

    -- проверка наличия id
    checkMsg = check_required_params(params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    update task
    set (result, success_rate, state, date_completed) = (params ->> 'result', (params ->> 'success_rate')::int,
                                                         'finished', now() at time zone '[[Config.Postgres.TimeZone]]')
    where id = (params ->> 'id')::int returning * into TaskRow;

    -- случай когда записи с таким id не найдено
    IF TaskRow.id ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
    END IF;

    result = row_to_json(TaskRow) :: JSONB;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;