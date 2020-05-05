-- параметры:
-- filename			      type: string
-- ext			          type: string
-- table_name			  type: string
-- table_id			      type: int
-- size			          type: int
-- options			      type: json

DROP FUNCTION IF EXISTS file_update(params JSONB);
CREATE OR REPLACE FUNCTION file_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE

    FileRow   file%ROWTYPE;
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
        checkMsg = check_required_params(params, ARRAY ['filename', 'ext', 'table_name', 'table_id']);
        IF checkMsg IS NOT NULL
        THEN
            RETURN checkMsg;
        END IF;
        -- проверка, что если тако файл уже загружен, то возвращаем уже существующий токен
        select token into tokenStr from file where filename = (params->>'filename') AND ext = (params->>'ext')
                                               AND table_name=(params->>'table_name') AND table_id = (params->>'table_id')::int and deleted= false;

        if tokenStr notnull then
            RETURN json_build_object('ok', TRUE, 'result', jsonb_build_object('token', tokenStr));
        end if;
        -- генерация токена
        SELECT md5(random() :: TEXT)
        INTO tokenStr;

        -- вариант создания
        EXECUTE ('INSERT INTO file (filename, ext, table_name, table_id, size, options, token) VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING *;')
            INTO FileRow
            USING
                params ->> 'filename',
                params ->> 'ext',
                params ->> 'table_name',
                (params ->> 'table_id')::int,
                (params ->> 'size')::int,
                coalesce((params ->> 'options')::jsonb, '{}'),
                tokenStr;
    ELSE
        -- вариант обновления существующей записи
        updateValue = '' || update_str_from_json(params, ARRAY [
            ['filename', 'filename', 'text'],
            ['ext', 'ext', 'text'],
            ['table_name', 'table_name', 'text'],
            ['table_id', 'table_id', 'number'],
            ['size', 'size', 'number'],
            ['options', 'options', 'jsonb'],
            ['deleted', 'deleted', 'bool']]);

        queryStr = concat('UPDATE file SET ', updateValue, ' WHERE id=', params->>'id',
                          ' RETURNING *');
        EXECUTE (queryStr)
            INTO FileRow;
    END IF;

    RETURN json_build_object('ok', TRUE, 'result', jsonb_build_object('token', FileRow.token));

END

$function$;

		