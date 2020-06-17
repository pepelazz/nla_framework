
DROP FUNCTION IF EXISTS chat_message_update(params JSONB);
CREATE OR REPLACE FUNCTION chat_message_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    ChatMessageRow     chat_message%ROWTYPE;
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
        checkMsg = check_required_params(params, ARRAY ['title', 'chat_id']);
        IF checkMsg IS NOT NULL
        THEN
            RETURN checkMsg;
        END IF;

        EXECUTE ('INSERT INTO chat_message (chat_id, user_id, title, options ) VALUES ($1, $2, $3, $4) RETURNING *;')
            INTO  ChatMessageRow
            USING
             (params ->> 'chat_id')::int,
             (params ->> 'user_id')::int,
             (params ->> 'title')::text,
             (params -> 'options')::jsonb;
    else
        updateValue = '' || update_str_from_json(params, ARRAY [
            ['title', 'title', 'text'],
            ['options', 'options', 'jsonb'],
            ['deleted', 'deleted', 'bool']
            ]);

        queryStr = concat('UPDATE chat_message SET ', updateValue, ' WHERE id=', params ->> 'id', ' RETURNING *;');

        EXECUTE (queryStr)
            INTO ChatMessageRow;

        -- случай когда записи с таким id не найдено
        IF row_to_json(ChatMessageRow) ->> 'id' ISNULL
        THEN
            RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
        END IF;

    end if;

    result = row_to_json(ChatMessageRow) :: JSONB;

    RETURN json_build_object('ok', TRUE, 'result', result);

END

$function$;