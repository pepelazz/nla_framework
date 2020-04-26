-- функция создания сообщения об изменении

CREATE OR REPLACE FUNCTION notify_event()
  RETURNS TRIGGER AS $$

DECLARE
  hString HSTORE;
  result  JSONB;
  r       RECORD;
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

  -- в случае изменения user добавляем поля auth_provider и auth_provider_id
--   IF TG_TABLE_NAME = 'user'
--   THEN
--     result = result || jsonb_build_object('auth_provider', r.auth_provider, 'auth_provider_id', r.auth_provider_id);
--   END IF;


  IF char_length(hString :: TEXT) > 0 -- отправляем notification только если есть изменения
  THEN
    PERFORM pg_notify('events', result :: TEXT);
  END IF;

  -- Result is ignored since this is an AFTER trigger
  RETURN NULL;
END;




$$ LANGUAGE plpgsql;