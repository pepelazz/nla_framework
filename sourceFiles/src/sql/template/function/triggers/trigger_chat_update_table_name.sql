
-- проставляем занчение полей table_name и task_type_title
CREATE OR REPLACE FUNCTION trigger_chat_update_table_name() RETURNS trigger AS $$
DECLARE
BEGIN

    IF (TG_OP = 'INSERT') THEN
        -- заполняем table_options
        NEW.table_options = '{}'::jsonb;
        -- for codeGenerate #trigger_task_update_table_name_slot
    end if;

    if (TG_OP = 'UPDATE') then
        NEW.table_name = old.table_name;
        NEW.table_id = old.table_id;
        NEW.table_options = old.table_options;
    end if;

  RETURN NEW;
END;

$$ LANGUAGE plpgsql;
