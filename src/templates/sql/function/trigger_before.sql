{{$doc := . -}}
-- функция триггер
DROP FUNCTION IF EXISTS {{.PgName}}_trigger_before() CASCADE;
CREATE OR REPLACE FUNCTION {{.PgName}}_trigger_before() RETURNS trigger AS
$$
DECLARE
       {{- .GetBeforeTriggerDeclareVars}}
BEGIN
        {{if .Sql.IsSearchText}}
        {{- /* заполнение ref полей */ -}}
        {{.GetBeforeTriggerFillRefVars}}
        -- заполняем options.title
        NEW.options = coalesce(OLD.options, '{}'::jsonb) || NEW.options || jsonb_build_object('title', jsonb_build_object({{.GetSearchTextJson}}));
        -- заполняем search_text
        NEW.search_text = concat({{.GetSearchTextString}});
        {{- end }}

        {{if .Sql.ComputedTitle}}
        NEW.title = {{.Sql.ComputedTitle}}
        {{- end }}

    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

