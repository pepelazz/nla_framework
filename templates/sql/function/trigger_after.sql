{{$doc := . -}}
-- функция триггер
DROP FUNCTION IF EXISTS {{.PgName}}_trigger_after() CASCADE;
CREATE OR REPLACE FUNCTION {{.PgName}}_trigger_after() RETURNS trigger AS
$$
DECLARE
        r record;
        jsonbEl      jsonb;
BEGIN
        {{.PrintAfterTriggerUpdateLinkedRecords}}

        {{if .IsRecursion}}
        if NEW.parent_id notnull then
            if exists(select true from {{.PgName}} where parent_id=NEW.parent_id and deleted=false) then
                update {{.PgName}} set is_folder=true where id = NEW.parent_id;
            else
                update {{.PgName}} set is_folder=false where id = NEW.parent_id;
            end if;
        end if;
        {{- end }}

        {{.Sql.Hooks.Print "triggerAfter" "AfterTriggerAfter"}}

    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

