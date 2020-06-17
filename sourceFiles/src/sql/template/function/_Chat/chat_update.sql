-- создание Chat
-- параметры:

DROP FUNCTION IF EXISTS chat_update(params JSONB);
CREATE OR REPLACE FUNCTION chat_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    ChatRow     chat%ROWTYPE;
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
        checkMsg = check_required_params(params, ARRAY ['table_name', 'table_id']);
        IF checkMsg IS NOT NULL
        THEN
            RETURN checkMsg;
        END IF;

        EXECUTE ('INSERT INTO chat (title, table_name, table_id, options ) VALUES ($1, $2, $3, $4) RETURNING *;')
            INTO  ChatRow
            USING
             (params ->> 'title')::text,
             (params ->> 'table_name')::text,
             (params ->> 'table_id')::int,
             (params -> 'options')::jsonb;
    else
        updateValue = '' || update_str_from_json(params, ARRAY [
            ['title', 'title', 'text'],
            ['options', 'options', 'jsonb'],
            ['deleted', 'deleted', 'bool']
            ]);

        queryStr = concat('UPDATE chat SET ', updateValue, ' WHERE id=', params ->> 'id', ' RETURNING *;');

        EXECUTE (queryStr)
            INTO ChatRow;

        -- случай когда записи с таким id не найдено
        IF row_to_json(ChatRow) ->> 'id' ISNULL
        THEN
            RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
        END IF;

    end if;

    result = row_to_json(ChatRow) :: JSONB;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;