-- функция триггер
DROP FUNCTION IF EXISTS user_trigger_before() CASCADE;
CREATE OR REPLACE FUNCTION user_trigger_before() RETURNS trigger AS
$$
DECLARE
        r record;
	senderTitle TEXT;
	recipientTitle TEXT;

       searchTxtVar TEXT := '';
BEGIN
        

    -- при удалении пользователя меняем статус на 'уволен'
    IF new.deleted = true and old.deleted != new.deleted then
        new.options = new.options || jsonb_build_object('state', 'fired');
    end if;


    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

