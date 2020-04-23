package types

import "fmt"

type (
	FldType struct {
		Name string
		NameRu string
		Type string
		Ext map[string]interface{}
		Vue FldVue
		Sql FldSql
	}

	FldVue struct {
		Name string
		NameRu string
		Type string
		RowCol [][]int
		Class []string
	}

	FldSql struct {
		Ref string
		IsUniq bool
	}
)

func (fld *FldType) PrintPgModel() string {
	typeStr := ""
	if fld.Type == "string" {
		if s, ok := fld.Ext["size"]; ok {
			typeStr = fmt.Sprintf("type=\"char\",\tsize=%v", s)
		} else {
			typeStr = `type="text"`
		}
	}
	res := fmt.Sprintf("{name=\"%s\",\t\t%s,\t comment=\"%s\"},", fld.Name, typeStr, fld.NameRu)

	return res
}
