{{$doc := . -}}
-- функция триггер
DROP FUNCTION IF EXISTS {{.PgName}}_trigger_before() CASCADE;
CREATE OR REPLACE FUNCTION {{.PgName}}_trigger_before() RETURNS trigger AS
$$
DECLARE
       {{- .GetBeforeTriggerDeclareVars}}
BEGIN
        {{.Sql.Hooks.Print "triggerBefore" "BeforeTriggerBefore"}}

        {{if .Sql.IsSearchText}}
        {{- /* заполнение ref полей */ -}}
        {{.GetBeforeTriggerFillRefVars}}
        -- заполняем options.title
        NEW.options = coalesce(OLD.options, '{}'::jsonb) || NEW.options || jsonb_build_object('title', jsonb_build_object({{.GetSearchTextJson}}));
        -- заполняем search_text
        {{if .GetSearchTextString}}
        NEW.search_text = concat({{.GetSearchTextString}});
        {{- else}}
        NEW.search_text = '';
        {{- end}}
        {{- end }}

        {{- range .Flds}}
        {{- if .Sql.FillValueInBeforeTrigger }}
        NEW.{{.Name}} = {{.Sql.FillValueInBeforeTrigger}};
        {{- end -}}
        {{- end -}}

        {{if .IsRecursion}}
        if NEW.parent_id notnull then
            if exists(select true from {{.PgName}} where parent_id=NEW.parent_id and deleted=false) then
                update {{.PgName}} set is_folder=true where id = NEW.parent_id;
            else
                update {{.PgName}} set is_folder=false where id = NEW.parent_id;
            end if;
        end if;
        {{- end }}



    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

