[[/* Шаблон для триггера */]]
[[$tableName := .TmplMain.TableName]]
[[- range $e := .TmplMain.Triggers]]
DROP TRIGGER IF EXISTS [[.Name]] ON [[$tableName]];
CREATE TRIGGER [[.Name]] [[uppercase .When]] ON [[$tableName]] [[uppercase .Ref]] EXECUTE PROCEDURE [[.FuncName]]();
[[- end -]]