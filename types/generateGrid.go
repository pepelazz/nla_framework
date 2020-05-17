package types

// генерация сетки для Vue
func makeGrid(doc DocType) []VueGridDiv {
	res := []VueGridDiv{}
	for _, fld := range doc.Flds {
		var cell *VueGridDiv
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
func getCell(grid *[]VueGridDiv, rowNum, colNum int, class string) *VueGridDiv {
	// создание строки если необходимо
	if len(*grid) < rowNum {
		for {
			*grid = append(*grid, VueGridDiv{Class: "row q-col-gutter-md q-mb-sm", Grid: []VueGridDiv{}})
			if len(*grid) == rowNum {
				break
			}
		}
	}
	// создание ячейки если необходимо
	if len((*grid)[rowNum-1].Grid) < colNum {
		for {
			(*grid)[rowNum-1].Grid = append((*grid)[rowNum-1].Grid, VueGridDiv{Class: class, Grid: []VueGridDiv{}})
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