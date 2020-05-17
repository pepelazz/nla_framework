-- параметры:
-- user_id			type: int
-- title            type: string
-- data			    type: json
-- table_name		type: string
-- table_id			type: int
-- state			type: string
-- type			    type: string
-- is_read			type: bool
-- options			type: json
-- is_read_conditions type: json

DROP FUNCTION IF EXISTS message_update(params JSONB);
CREATE OR REPLACE FUNCTION message_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE

    MessageRow   message%ROWTYPE;
    result      JSONB;
    updateValue TEXT;
    queryStr    TEXT;
    checkMsg    TEXT;
    tokenStr    TEXT;

BEGIN

    -- проверика наличия обязательных параметров
    checkMsg = check_required_params(params, ARRAY ['id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    IF (params->>'id')::int = -1
    THEN
        checkMsg = check_required_params(params, ARRAY ['user_id', 'title']);
        IF checkMsg IS NOT NULL
        THEN
            RETURN checkMsg;
        END IF;
        -- вариант создания
        EXECUTE ('INSERT INTO message (user_id, title, data, table_name, table_id, state, type, is_read, is_read_conditions, options) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)  RETURNING *;')
            INTO MessageRow
            USING
                (params ->> 'user_id')::int,
                params ->> 'title',
                coalesce((params ->> 'data')::jsonb, '{}'),
                params ->> 'table_name',
                (params ->> 'table_id')::int,
                params ->> 'state',
                params ->> 'type',
                coalesce((params ->> 'is_read')::bool, false),
                coalesce((params ->> 'is_read_conditions')::jsonb, '[]'),
                coalesce((params ->> 'options')::jsonb, '{}');
    ELSE
        -- вариант обновления существующей записи
        updateValue = '' || update_str_from_json(params, ARRAY [
            ['user_id', 'user_id', 'number'],
            ['title', 'title', 'text'],
            ['data', 'data', 'jsonb'],
            ['table_name', 'table_name', 'text'],
            ['table_id', 'table_id', 'number'],
            ['state', 'state', 'text'],
            ['type', 'type', 'text'],
            ['is_read', 'is_read', 'bool'],
            ['is_read_conditions', 'is_read_conditions', 'jsonb'],
            ['options', 'options', 'jsonb'],
            ['deleted', 'deleted', 'bool']]);

        queryStr = concat('UPDATE message SET ', updateValue, ' WHERE id=', params->>'id',
                          ' RETURNING *');
        EXECUTE (queryStr)
            INTO MessageRow;
    END IF;

    RETURN json_build_object('ok', TRUE, 'result',row_to_json(MessageRow));

END

$function$;

		