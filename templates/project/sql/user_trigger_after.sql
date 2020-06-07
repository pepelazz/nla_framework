-- функция триггер
DROP FUNCTION IF EXISTS user_trigger_after() CASCADE;
CREATE OR REPLACE FUNCTION user_trigger_after() RETURNS trigger AS
$$
DECLARE
        r record;
BEGIN
        [[PrintUserAfterTriggerUpdateLinkedRecords]]

    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

