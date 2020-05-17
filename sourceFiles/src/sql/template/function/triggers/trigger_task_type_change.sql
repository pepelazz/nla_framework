
-- обновляем task_title в таблице task, а также не даем менять table_name
CREATE OR REPLACE FUNCTION trigger_task_type_change() RETURNS trigger AS $$
DECLARE
    taskTypeRow task_type%ROWTYPE;
BEGIN

    NEW.table_name = old.table_name;
    IF new.title != old.title THEN
        update task set task_type_title = new.title where task_type_id = new.id;
    end if;

  RETURN NEW;
END;

$$ LANGUAGE plpgsql;
