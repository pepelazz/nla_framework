package types

type (
	FldType struct {
		Name string
		NameRu string
		Type string
		Ext map[string]interface{}
		Vue FldVue
	}

	FldVue struct {
		Name string
		NameRu string
		Type string
		RowCol [][]int
		Class []string
	}
)
