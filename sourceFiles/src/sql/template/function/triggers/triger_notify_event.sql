-- функция создания сообщения об изменении

CREATE OR REPLACE FUNCTION notify_event()
  RETURNS TRIGGER AS $$

DECLARE
  hString HSTORE;
  result  JSONB;
  r       RECORD;
  authToken text;
  userFullname text;
  userOptions jsonb;
  taskExecutorFullname text;
  taskManagerFullname text;
  taskTypeOptions jsonb;
BEGIN

  IF (TG_OP = 'DELETE')
  THEN
    r = OLD;
    -- в случае удаления отправляем всю запись
    hString = (hstore(OLD) || 'delete=>true' :: HSTORE) - ARRAY ['id', 'updated_at', 'created_at'];
  ELSIF (TG_OP = 'INSERT')
    THEN
      r = NEW;
      hString = hstore(NEW) - ARRAY ['id', 'updated_at', 'created_at', 'password'];
  ELSIF (TG_OP = 'UPDATE')
    THEN
      r = NEW;
      -- считаем дельту между старой и новой версией
      -- из полученной дельты убираем поле updated_at
      hString = hstore(NEW) - hstore(OLD) - ARRAY ['updated_at', 'password'];

  END IF;

  result = jsonb_build_object('table', TG_TABLE_NAME, 'id', r.id, 'flds',
                              hstore_to_json_loose(hString));

  -- в случае изменения user добавляем поля auth_token
  IF TG_TABLE_NAME = 'user'
  THEN
      select auth_token into authToken from user_auth where user_id = r.id;
      result = result || jsonb_build_object('auth_token', authToken);
  END IF;

  -- в случае изменения message добавляем поля id и TG_OP
  IF TG_TABLE_NAME = 'message'
  THEN
      select fullname, options into userFullname, userOptions from "user" where id = r.user_id;
      result = jsonb_set(result, '{flds}', result->'flds' || jsonb_build_object('id', r.id, 'tg_op', TG_OP, 'sse_type', 'message', 'user_fullname', userFullname, 'user_options', userOptions));
  END IF;

  -- в случае изменения task добавляем fullname по исполнителю и менеджеру
  IF TG_TABLE_NAME = 'task'
  THEN
      select fullname into taskExecutorFullname from "user" where id = r.executor_id;
      select fullname into taskManagerFullname from "user" where id = r.manager_id;
      select options into taskTypeOptions from task_type where id = r.task_type_id;
      result = jsonb_set(result, '{flds}', result->'flds' || row_to_json(r)::jsonb || jsonb_build_object('id', r.id, 'tg_op', TG_OP, 'sse_type', 'task', 'executor_fullname', taskExecutorFullname, 'manager_fullname', taskManagerFullname, 'task_type_options', taskTypeOptions));
  END IF;

  IF char_length(hString :: TEXT) > 0 -- отправляем notification только если есть изменения
  THEN
    PERFORM pg_notify('events', result :: TEXT);
  END IF;

  -- Result is ignored since this is an AFTER trigger
  RETURN NULL;
END;

$$ LANGUAGE plpgsql;