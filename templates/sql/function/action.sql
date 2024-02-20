-- вызов action in state machine: [[.NameRu]]
-- параметры:

DROP FUNCTION IF EXISTS [[.PgName]]_action(params JSONB);
CREATE OR REPLACE FUNCTION [[.PgName]]_action(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    r           [[.PgName]]%ROWTYPE;
    rJson       jsonb;
    result      json;
    checkMsg    TEXT;
    updateValue TEXT;
    newStateName TEXT;
    allowedStates TEXT[];
    updateFlds  text[];
    copyToParamsFlds  text[];
    copyFldName  text;
    arrFlds     VARCHAR[] := '{{options, options, jsonb}}'::VARCHAR[];
    m           VARCHAR[];
    [[tmplSqlActionPrintRefUpdateVarDeclare .]]
BEGIN

    -- проверка наличия id
    checkMsg = check_required_params(params, ARRAY ['id', 'action_name', 'user_id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    select * into r from [[.PgName]] where id = (params ->> 'id')::int;
    IF r.id ISNULL
    THEN
        RETURN json_build_object('ok', FALSE, 'message', 'wrong id');
    END IF;
    -- создаем json объект из записи, чтобы можно было обращаться к значениям колонок через название переменных
    rJson = row_to_json(r)::jsonb;

    if params->'options'->'states' isnull then
        params = jsonb_set(params, '{options, states}'::text[], '[]'::jsonb);
    end if;

    case params->>'action_name'
[[tmplSqlActionPrintCaseBlock .]]
        else
            RETURN json_build_object('ok', FALSE, 'message', 'wrong action name');
        end case;


    -- проверка что экшен из того стейта, в котором сейчас находится документ
    if array_length(allowedStates, 1) > 0 AND r.state != ALL(allowedStates) then
        RETURN json_build_object('ok', FALSE, 'message', 'wrong action for current state');
    end if;

    --записываем название нового стейта
    params = params || jsonb_build_object('state', newStateName);

    -- если это смена стейта.
    if r.state != newStateName then
        -- сразу сохраняем коммент, потому что это поле есть во всех стейтах
        params = params || jsonb_set(params, '{options, states, 0}'::text[] || '{comment}'::text[], rJson -> 'comment');

        -- копируем в options значения необходимых полей из предыдущего стейта
        if array_length(copyToParamsFlds, 1) > 0 then
            FOREACH copyFldName IN ARRAY copyToParamsFlds
                LOOP
                    params = params || jsonb_set(params, '{options, states, 0}'::text[] || copyFldName, rJson -> copyFldName);
                    [[tmplSqlActionPrintRefUpdateBlock .]]
                END LOOP;
        end if;
    -- прописываем кто изменил статус и когда
--     params = params || jsonb_build_object('options', options_add_fld((params->>'user_id')::int, params->'options', 'states', jsonb_build_object('state', newStateName)));
    params = params || jsonb_build_object('options', jsonb_insert( params->'options', '{states, 0}'::text[], jsonb_build_object('state', newStateName, 'user_id', (params->>'user_id')::int, 'date', now() at time zone '[[GetPgTimeZone]]')));

    end if;

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

    EXECUTE (concat('UPDATE [[.PgName]] SET ', '' || update_str_from_json(params, arrFlds), ' WHERE id=', params ->> 'id', ' RETURNING *;'))
        INTO r;

    -- действия в случае успешного выполнения action
    [[tmplSqlActionPrintAfterHook .]]

    RETURN json_build_object('ok', TRUE, 'result', row_to_json(r) :: JSONB);
END

$function$;