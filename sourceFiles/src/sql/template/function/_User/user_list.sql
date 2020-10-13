-- получение списка пользователей
-- параметры:
-- state           type: user_state - статус пользователя
-- deleted         type: bool - удаленные / существующие. Дефолт: false
-- order_by        type: string - поле для сортировки и направление сортировки. Например, orderBy: "id desc"
-- page            type: int - номер страницы. Дефолт: 1
-- per_page        type: int - количество записей на странице. Дефолт: 1000
-- search_fullname type: string - текстовый поиск по fullname
-- roles           type: bool - ожидающие авторизации

DROP FUNCTION IF EXISTS user_list(params JSONB );
CREATE OR REPLACE FUNCTION user_list(params JSONB)
  RETURNS JSON
LANGUAGE plpgsql
AS $function$

DECLARE

  result       JSON;
  condQueryStr TEXT;
  whereStr     TEXT;

BEGIN

  -- сборка условия WHERE (where_str_build - функция из папки base)
  whereStr = where_str_build(params, 'doc', ARRAY [
  ['enum', 'state', 'doc.state'],
  ['jsonArrayText', 'role', 'doc.role'],
  ['ilike', 'search_fullname', 'doc.fullname'],
  ['ilike', 'search_text', 'doc.fullname']
  ]);

  -- финальная сборка строки с условиями выборки (build_query_part_for_list - функция из папки base)
  condQueryStr = '' || whereStr || build_query_part_for_list(params);

  EXECUTE (
    ' SELECT array_to_json(array_agg(t)) FROM (SELECT id, avatar, first_name, last_name, fullname, title, role, email, options, deleted, created_at  FROM "user" as doc ' ||  condQueryStr || ') AS t')
  INTO result;

  RETURN json_build_object('ok', TRUE, 'result', coalesce(result, '[]'));

END

$function$;




