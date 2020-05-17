-- параметры:
-- roles			type: []string
-- title            type: string
-- data			    type: json
-- table_name		type: string
-- table_id			type: int
-- is_read_conditions type: json
-- state			type: string
-- type			    type: string
-- is_read			type: bool
-- options			type: json

DROP FUNCTION IF EXISTS message_for_role(params JSONB);
CREATE OR REPLACE FUNCTION message_for_role(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE

    MessageRow  message%ROWTYPE;
    result      JSONB;
    updateValue TEXT;
    queryStr    TEXT;
    checkMsg    TEXT;
    roleStr    TEXT;
    roleStr1    TEXT :='admin';
    r           record;
    r1           record;

BEGIN

    -- проверика наличия обязательных параметров
    checkMsg = check_required_params(params, ARRAY ['roles']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    for r in (select jsonb_array_elements_text(params->'roles') roleStr)
        LOOP
            for r1 in select * from "user" where r.roleStr = ANY(role) and deleted = false
            loop
                Perform message_update(params || jsonb_build_object('id', -1, 'user_id', r1.id));
            end loop;
        END LOOP;

    RETURN json_build_object('ok', TRUE, 'result', 'ok');

END

$function$;

		