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
        {{- if eq .Vue.Type "tags"}}
         searchTxtVar = searchTxtVar || (NEW.{{.Name}})::text;
        {{- end -}}
        {{- if .Sql.FillValueInBeforeTrigger }}
        NEW.{{.Name}} = {{.Sql.FillValueInBeforeTrigger}};
        {{- end -}}
        {{- if eq .Vue.Type "phone" }}
        NEW.{{.Name}} = phone_change_8_to_7(NEW.{{.Name}});
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

        {{.Sql.Hooks.Print "triggerBefore" "AfterTriggerBefore"}}


    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

