{{$doc := . -}}
-- функция триггер
DROP FUNCTION IF EXISTS {{.PgName}}_trigger_before() CASCADE;
CREATE OR REPLACE FUNCTION {{.PgName}}_trigger_before() RETURNS trigger AS
$$
DECLARE
       {{- .GetBeforeTriggerDeclareVars}}
       searchTxtVar TEXT := '';
BEGIN
        {{.Sql.Hooks.Print "triggerBefore" "BeforeTriggerBefore"}}

        {{if .Sql.IsSearchText}}
        {{- /* заполнение ref полей */ -}}
        {{.GetBeforeTriggerFillRefVars}}
        {{- end}}

        {{- range .Flds}}
        {{- if .Sql.FillValueInBeforeTrigger }}
        NEW.{{.Name}} = {{.Sql.FillValueInBeforeTrigger}};
        {{- end -}}
        {{- end -}}

        {{if .Sql.IsSearchText}}
        -- заполняем options.title
        NEW.options = coalesce(OLD.options, '{}'::jsonb) || NEW.options || jsonb_build_object('title', jsonb_build_object({{.GetSearchTextJson}}));
        -- заполняем search_text
        {{if .GetSearchTextString}}
        NEW.search_text = concat({{.GetSearchTextString}}, ' ', searchTxtVar);
        {{- else}}
        NEW.search_text = '';
        {{- end}}
        {{- end }}


    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

