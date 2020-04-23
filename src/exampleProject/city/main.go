package city

import (
	"github.com/pepelazz/projectGenerator/src/types"
)

func GetDoc() types.DocType {
	doc := types.DocType{
		Name:   "city",
		NameRu: "город",
		Flds: []types.FldType{
			{Name: "title", NameRu: "название", Type: "string", Ext: map[string]interface{}{"size": 50}, Vue: types.FldVue{RowCol: [][]int{{1,1}}, Class: []string{"col-12"}}},
			{Name: "region", NameRu: "регион", Type: "string", Ext: map[string]interface{}{"size": 50}, Vue: types.FldVue{RowCol: [][]int{{2,1}}, Class: []string{"col-6"}}},
			{Name: "index", NameRu: "индекс", Type: "string", Ext: map[string]interface{}{"size": 50}, Vue: types.FldVue{RowCol: [][]int{{2,2}}, Class: []string{"col-6"}}},
			//{Name: "fld1", Type: "string", Ext: map[string]interface{}{"size": 50}, Vue: types.FldVue{RowCol: [][]int{{1,2},{1,1}}}},
			//{Name: "fld2", Type: "string", Ext: map[string]interface{}{"size": 50}, Vue: types.FldVue{RowCol: [][]int{{1,2},{1,2}}}},
			//{Name: "fld3", Type: "string", Ext: map[string]interface{}{"size": 50}, Vue: types.FldVue{RowCol: [][]int{{2,1}}}},
			//{Name: "fld4", Type: "string", Ext: map[string]interface{}{"size": 50}, Vue: types.FldVue{RowCol: [][]int{{3,2},{1,1}}, Class: []string{"col-3", "col-4"}}},
		},
		Vue: types.DocVue{
			Route: "city",
		},
	}
	doc.Vue.Grid = makeGrid(doc)
	return doc
}

func makeGrid(doc types.DocType) []types.VueGridDiv {
	res := []types.VueGridDiv{}
	for _, fld := range doc.Flds {
		var cell *types.VueGridDiv
		for i, arr := range fld.Vue.RowCol {
			// дефолтная ширина колонки, которую перезаписываем классом в зависимости от уровня вложенности
			class := "col-6"
			if len(fld.Vue.Class) > i {
				class = fld.Vue.Class[i]
			}
			if cell == nil {
				cell = getCell(&res, arr[0], arr[1], class)
			} else {
				cell = getCell(&cell.Grid, arr[0], arr[1], class)
			}
			if i == len(fld.Vue.RowCol) - 1 {
				// если последний элемент в массиве, то создаем финальный div для field
				cell.Fld = fld
			}
		}
	}
	return res
}

// функция по созданию или получению существующей ячейки из сетки
func getCell(grid *[]types.VueGridDiv, rowNum, colNum int, class string) *types.VueGridDiv {
	// создание строки если необходимо
	if (len(*grid) < rowNum) {
		for {
			*grid = append(*grid, types.VueGridDiv{Class: "row q-col-gutter-md q-mb-sm", Grid: []types.VueGridDiv{}})
			if len(*grid) == rowNum {
				break
			}
		}
	}
	// создание ячейки если необходимо
	if (len((*grid)[rowNum-1].Grid) < colNum) {
		for {
			(*grid)[rowNum-1].Grid = append((*grid)[rowNum-1].Grid, types.VueGridDiv{Class: class, Grid: []types.VueGridDiv{}})
			if len((*grid)[rowNum-1].Grid) == colNum {
				break
			}
		}
	}
	return &(*grid)[rowNum-1].Grid[colNum-1]
}

// пример структуры для темплейта сетки
//func makeGridOld() []types.VueGridDiv  {
//	res := []types.VueGridDiv{}
//	res = append(res, types.VueGridDiv{
//		Class: "row q-col-gutter-md q-mb-sm",
//		Grid: []types.VueGridDiv{
//			{Class: "col-6", Fld: types.FldType{Name: "fld1"}},
//			{
//				Class: "col-6",
//				Grid: []types.VueGridDiv {
//					{Class: "row q-col-gutter-md q-mb-sm", Fld: types.FldType{Name: "fld2"}},
//					{Class: "row q-col-gutter-md q-mb-sm", Fld: types.FldType{Name: "fld3"}},
//					{Class: "row q-col-gutter-md q-mb-sm", Fld: types.FldType{Name: "fld4"}},
//				},
//			},
//		},
//	})
//	res = append(res, types.VueGridDiv{
//		Class: "row",
//		Grid: []types.VueGridDiv{
//			{Class: "col-6", Fld: types.FldType{Name: "fld5"}},
//			{Class: "col-6", Fld: types.FldType{Name: "fld6"}},
//		},
//	})
//	return res
//}