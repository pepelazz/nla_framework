-- update в случае если документ реализует поведение state machine
-- создание [[.NameRu]]
-- параметры:

DROP FUNCTION IF EXISTS [[.Name]]_update(params JSONB);
CREATE OR REPLACE FUNCTION [[.Name]]_update(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    r           [[.Name]]%ROWTYPE;
    rNew        [[.Name]]%ROWTYPE;
    checkMsg    TEXT;
    result      JSONB;
    updateValue TEXT;
    queryStr    TEXT;
    updateFlds  text[];
    arrFlds     VARCHAR[] := '{{options, options, jsonb}}'::VARCHAR[];
    m           VARCHAR[];
BEGIN

    -- проверика наличия id
    checkMsg = check_required_params(params, ARRAY ['id', 'user_id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    select * into r from [[.Name]] where id = (params ->> 'id')::int;
    IF r.id ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
    END IF;

    case r.state
[[tmplSqlUpdatePrintCaseBlock .]]
        else
            RETURN json_build_object('ok', FALSE, 'message', 'wrong stateName in [[.Name]]_update');
        end case;

    -- оставляем только поля, которые указаны в updateFlds, котрые отфильтрованы в зависимости от текущего стейта
    FOREACH m SLICE 1 IN ARRAY ARRAY [
[[.PrintSqlFuncUpdateFlds]]
        ['options', 'options', 'jsonb'],
        ['deleted', 'deleted', 'bool']
        ]
        loop
            IF m[1] = ANY (updateFlds) then
                arrFlds = arrFlds || m;
            end if;
        end loop;

    EXECUTE (concat('UPDATE [[.Name]] SET ', '' || update_str_from_json(params, arrFlds), ' WHERE id=', params ->> 'id', ' RETURNING *;'))
        INTO rNew;

    [[.Sql.Hooks.Print "update" "afterInsertUpdate"]]

    RETURN [[.Name]]_get_by_id(params);

END

$function$;