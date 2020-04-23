package client

import "github.com/pepelazz/projectGenerator/src/types"

func GetDoc() types.DocType  {
	doc := types.DocType{
		Name: "client",
		NameRu: "клиент",
		Flds: []types.FldType{
			{Name: "title", Type: "string", Ext: map[string]interface{}{"size": 50}},
		},
		Vue: types.DocVue{Route: "client"},
	}
	return doc
}
