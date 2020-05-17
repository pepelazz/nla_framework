{{$doc := . -}}
-- функция триггер
DROP FUNCTION IF EXISTS {{.PgName}}_trigger_after() CASCADE;
CREATE OR REPLACE FUNCTION {{.PgName}}_trigger_after() RETURNS trigger AS
$$
DECLARE
        r record;
BEGIN
        {{.PrintAfterTriggerUpdateLinkedRecords}}

    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

