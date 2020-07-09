-- получение списка {{.NameRu}}
-- параметры:
-- deleted         type: bool - удаленные / существующие. Дефолт: false
-- order_by        type: string - поле для сортировки и направление сортировки. Например, orderBy: "id desc"
-- page            type: int - номер страницы. Дефолт: 1
-- per_page        type: int - количество записей на странице. Дефолт: 1000
-- search_text     type: string - текстовый поиск

DROP FUNCTION IF EXISTS {{.PgName}}_list(params JSONB);
CREATE OR REPLACE FUNCTION {{.PgName}}_list(params JSONB)
    RETURNS JSON
    LANGUAGE plpgsql
AS
$function$

DECLARE
    result       JSON;
    condQueryStr TEXT;
    whereStr     TEXT;
    checkMsg     TEXT;
BEGIN

    checkMsg = check_required_params(params, ARRAY ['user_id']);
    IF checkMsg IS NOT NULL
    THEN
        RETURN checkMsg;
    END IF;

    {{.Sql.Hooks.Print "list" "listBeforeBuildWhere"}}

    -- сборка условия WHERE (where_str_build - функция из папки base)
    whereStr = where_str_build(params, 'doc', ARRAY [
        {{.PrintSqlFuncListWhereCond}}
    ]);

    {{.Sql.Hooks.Print "list" "listAfterBuildWhere"}}

    -- финальная сборка строки с условиями выборки (build_query_part_for_list - функция из папки base)
    condQueryStr = '' || whereStr || build_query_part_for_list(params);

    {{.PrintSqlFuncList}}

    RETURN json_build_object('ok', TRUE, 'result', coalesce(result, '[]'));

END
$function$;




