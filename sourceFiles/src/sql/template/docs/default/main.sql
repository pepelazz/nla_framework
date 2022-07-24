[[$DocTypeUnderscore := camelToSnake .DocType]]
[[$TableName := .TmplMain.TableName]]
[[- if .TmplMain.Enums]]
DO $$
BEGIN
  [[- range $k, $v := .TmplMain.Enums ]]
  [[- $typeName := printf "%s_%s" $DocTypeUnderscore (camelToSnake $k) ]]
  IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = '[[$typeName]]')
  THEN
CREATE TYPE [[$typeName]] AS ENUM ([[joinWithQuotes $v "," "'"]]);
END IF;
[[- end]]
END;
$$ LANGUAGE plpgsql;
[[- end]]

CREATE TABLE IF NOT EXISTS [[.TmplMain.TableName]] (
[[- range $i, $fld := .TmplMain.Fields]][[if $i]],
[[end]]
[[if eq .Type "serial" -]]     [[- template "docFld_serial" . -]]    [[- end -]]
[[if eq .Type "char" -]]       [[- template "docFld_char" . -]]      [[- end -]]
[[if eq .Type "text" -]]       [[- template "docFld_text" . -]]      [[- end -]]
[[if eq .Type "int" -]]        [[- template "docFld_int" . -]]       [[- end -]]
[[if eq .Type "bigint" -]]     [[- template "docFld_bigint" . -]]    [[- end -]]
[[if eq .Type "uuid" -]]       [[- template "docFld_uuid" . -]]      [[- end -]]
[[if eq .Type "double" -]]     [[- template "docFld_double" . -]]    [[- end -]]
[[if eq .Type "bool" -]]       [[- template "docFld_bool" . -]]      [[- end -]]
[[if eq .Type "json" -]]       [[- template "docFld_json" . -]]      [[- end -]]
[[if eq .Type "jsonb" -]]      [[- template "docFld_jsonb" . -]]     [[- end -]]
[[if eq .Type "text[]" -]]     [[- template "docFld_text[]" . -]]    [[- end -]]
[[if eq .Type "int[]" -]]      [[- template "docFld_int[]" . -]]     [[- end -]]
[[if eq .Type "double precision[]" -]]      [[- template "docFld_double[]" . -]]     [[- end -]]
[[if eq .Type "timestamp" -]]  [[- template "docFld_timestamp" . -]] [[- end -]]
[[if eq .Type "time" -]]       [[- template "docFld_time" . -]] [[- end -]]
[[if eq .Type "constraint" -]] [[- template "docFld_constraint" . -]][[- end -]]
[[if eq .Type "enum" -]]       [[- template "docFld_enum" dict "Fld" . "DocTypeUnderscore" $DocTypeUnderscore -]]   [[- end -]]
[[if eq .Type "tsvector" -]]   [[- template "docFld_tsvector" . -]]    [[- end -]]
[[ end -]]
);

-- скрипты по изменению таблицы AlterScripts
[[ range $e := .TmplMain.AlterScripts]]
[[.Name]]
[[- end ]]

-- комментарии к полям таблицы
[[range $i, $fld := .TmplMain.Fields]] [[ if .Comment]]  COMMENT ON COLUMN [[$TableName]].[[ .Name ]] IS '[[ .Comment ]]';  [[end]]
[[ end]]

-- комментарий к таблице
[[ if .TmplMain.TableComment]] COMMENT ON TABLE [[$TableName]] IS '[[.TmplMain.TableComment]]'; [[end]]

[[- range $e := .TmplMain.FkConstraints]]
ALTER TABLE [[ $TableName ]] DROP CONSTRAINT IF EXISTS [[.Name]];
ALTER TABLE [[ $TableName ]] ADD CONSTRAINT [[.Name]] [[ if .Ref]]FOREIGN KEY ([[.Fld]]) REFERENCES [[.Ref]] ([[.Fk]])[[end]] [[.Ext]];
[[- end ]]

[[- range $e := .TmplMain.Indexes]]
CREATE [[if .Unique]]UNIQUE[[end]] INDEX IF NOT EXISTS [[.Name]] ON [[$TableName]] [[if .Using]]USING [[.Using]][[end]] ([[joinWithQuotes .Fld "," ""]]) [[if .Where]] WHERE [[.Where]][[end]];
[[- end ]]

[[- template "tableExt"]]

[[ define "docFld_serial" ]]    [[ .Name ]] SERIAL PRIMARY KEY [[- end -]]
[[ define "docFld_char" ]]      [[ .Name ]] CHARACTER VARYING([[.Size]]) [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_text" ]]      [[ .Name ]] TEXT [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_int" ]]       [[ .Name ]] INTEGER [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_bigint" ]]    [[ .Name ]] BIGINT [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_uuid" ]]      [[ .Name ]] UUID  [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_double" ]]    [[ .Name ]] DOUBLE PRECISION [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_json" ]]      [[ .Name ]] JSON [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_jsonb" ]]     [[ .Name ]] JSONB [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_bool" ]]      [[ .Name ]] BOOL [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_text[]" ]]    [[ .Name ]] TEXT [] [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_int[]" ]]     [[ .Name ]] INT [] [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_double[]" ]]     [[ .Name ]] DOUBLE PRECISION [] [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_timestamp" ]] [[ .Name ]] TIMESTAMP [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_time" ]]      [[ .Name ]] TIME [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_tsvector" ]]  [[ .Name ]] TSVECTOR [[ uppercase .Ext -]] [[- end -]]
[[ define "docFld_constraint" ]][[ .Ext ]][[- end -]]

[[ define "docFld_enum" ]]
[[- $Fld := index . "Fld" -]]
[[- $enumName := printf "%s_%s" .DocTypeUnderscore (camelToSnake $Fld.Enum.Name) -]]
[[- $Fld.Name ]] [[$enumName]] [[if $Fld.Enum.Default]] DEFAULT '[[$Fld.Enum.Default]]' :: [[$enumName]] [[end]]
[[- end -]]




