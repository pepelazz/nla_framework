package types

import (
	"fmt"
)

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
		IsSearch bool
		IsRequired bool
		Ref      string
		IsUniq   bool
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
	res := fmt.Sprintf("\t{name=\"%s\",\t\t\t\t\t%s,\t comment=\"%s\"}", fld.Name, typeStr, fld.NameRu)

	return res
}

func (fld *FldType) PgInsertType() string {
	switch fld.Type {
	case "double":
		return "double precision"
	case "string":
		return "text"
	default:
		return fld.Type
	}
}

func (fld *FldType) PgUpdateType() string {
	switch fld.Type {
	case "int", "double":
		return "number"
	case "string":
		return "text"
	default:
		return fld.Type
	}
}
