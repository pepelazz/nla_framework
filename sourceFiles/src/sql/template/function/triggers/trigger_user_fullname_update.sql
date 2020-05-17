
-- функция обновления поля fullname
CREATE OR REPLACE FUNCTION trigger_user_fullname_update() RETURNS trigger AS $$

BEGIN

    NEW.fullname  := btrim(COALESCE(NEW.last_name, NEW.last_name, '') || ' ' || COALESCE(NEW.first_name, NEW.first_name, ''));
    NEW.title  := btrim(COALESCE(NEW.last_name, NEW.last_name, '') || ' ' || COALESCE(NEW.first_name, NEW.first_name, ''));

  RETURN NEW;
END;

$$ LANGUAGE plpgsql;
