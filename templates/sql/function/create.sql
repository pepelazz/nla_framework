-- создание in state machine: [[.NameRu]]
-- параметры:

DROP FUNCTION IF EXISTS [[.PgName]]_create(params JSONB);
CREATE OR REPLACE FUNCTION [[.PgName]]_create(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    [[.PgName]]Row     [[.PgName]]%ROWTYPE;
    checkMsg    TEXT;
 [[.Sql.Hooks.Print "update" "declareVars"]]
BEGIN

    [[.PrintSqlFuncUpdateCheckParams]]

    [[.Sql.Hooks.Print "update" "beforeInsertUpdate"]]

    [[if .RequiredFldsString -]]
    -- проверка наличия обязательных параметров
    checkMsg = check_required_params(params, ARRAY [ [[.RequiredFldsString]] ]);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;
    [[end -]]
    -- прописываем кто изменил статус и когда
    params = params || jsonb_build_object('options', options_add_fld((params->>'user_id')::int, coalesce(params->'options', '{}'::jsonb), 'states',
                                                                     jsonb_build_object('state', '[[with .StateMachine.GetFirstState]][[.Title]][[end]]')));

    [[.Sql.Hooks.Print "update" "beforeInsert"]]

    [[.PrintSqlFuncInsertNew]]

    [[.Sql.Hooks.Print "create" "afterCreate"]]

    RETURN json_build_object('ok', TRUE, 'result', row_to_json([[.PgName]]Row) :: JSONB);

END

$function$;