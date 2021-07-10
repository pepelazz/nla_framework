-- параметры:
-- user_id			type: int
-- id               type: int

DROP FUNCTION IF EXISTS message_mark_as_read(params JSONB);
CREATE OR REPLACE FUNCTION message_mark_as_read(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    messageRow  message%ROWTYPE;
    checkMsg    TEXT;
    r    RECORD;
    isTrue bool;
BEGIN

    -- проверика наличия обязательных параметров
    checkMsg = check_required_params(params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    select * into messageRow from message where id = (params->>'id')::int;
    if messageRow.is_read != true then
        -- проверка возможных условий
        if messageRow.is_read_conditions notnull then
            for r in (select jsonb_array_elements(messageRow.is_read_conditions) cond)
            loop
                execute r.cond ->> 'query' into isTrue;
                if isTrue is not true then
                    return jsonb_build_object('ok', false, 'message', r.cond ->>'message');
                end if;
            end loop;
        end if;
    end if;

    return message_update(params || jsonb_build_object('is_read', true));

END

$function$;

		