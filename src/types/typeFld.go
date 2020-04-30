package types

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/src/utils"
)

type (
	FldType struct {
		Name   string
		NameRu string
		Type   string
		Vue    FldVue
		Sql    FldSql
	}

	FldVue struct {
		Name      string
		NameRu    string
		Type      string
		RowCol    [][]int
		Class     []string
		IsRequred bool
		Ext       map[string]string
	}

	FldSql struct {
		IsSearch   bool
		IsRequired bool
		Ref        string
		IsUniq     bool
		Size       int
		IsOptionFld bool
	}
)

func (fld *FldType) PrintPgModel() string {
	typeStr := fmt.Sprintf(`type="%s"`, fld.Type)
	extStr := ""
	if fld.Type == "string" {
		if fld.Sql.Size > 0 {
			typeStr = fmt.Sprintf("type=\"char\",\tsize=%v", fld.Sql.Size)
		} else {
			typeStr = `type="text"`
		}
	}
	if utils.CheckContainsSliceStr(fld.Type, "date", "datetime") {
		typeStr = `type="timestamp"`
	}
	if fld.Sql.IsRequired {
		extStr = "not null"
	}
	// ext может быть пустой
	ext := ""
	if len(extStr) > 0 {
		ext = fmt.Sprintf(" \text=\"%s\",", extStr)
	}
	res := fmt.Sprintf("\t{name=\"%s\",\t\t\t\t\t%s,%s\t comment=\"%s\"}", fld.Name, typeStr, ext, fld.NameRu)

	return res
}

func (fld *FldType) PgInsertType() string {
	switch fld.Type {
	case "double":
		return "double precision"
	case "string":
		return "text"
	case "date", "datetime":
		return "timestamp"
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
	case "date", "datetime":
		return "timestamp"
	default:
		return fld.Type
	}
}
