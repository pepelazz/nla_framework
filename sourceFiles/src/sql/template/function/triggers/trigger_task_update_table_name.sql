
-- проставляем занчение полей table_name и task_type_title
CREATE OR REPLACE FUNCTION trigger_task_update_table_name() RETURNS trigger AS $$
DECLARE
    taskTypeRow task_type%ROWTYPE;
BEGIN

    IF (TG_OP = 'INSERT') THEN
        select * into taskTypeRow from task_type where id = NEW.task_type_id;
        NEW.table_name = taskTypeRow.table_name;
        NEW.task_type_title = taskTypeRow.title;
        -- заполняем table_options
        NEW.table_options = '{}'::jsonb;
        -- for codeGenerate #trigger_task_update_table_name_slot
    end if;

    if (TG_OP = 'UPDATE') then
        NEW.table_name = old.table_name;
        NEW.table_id = old.table_id;
        NEW.table_options = old.table_options;
        new.task_type_id = old.task_type_id;
    end if;

  RETURN NEW;
END;

$$ LANGUAGE plpgsql;
