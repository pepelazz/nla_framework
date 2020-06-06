package types

import (
	"fmt"
	"strings"
)

type (
	BitrixFld struct {
		Name string
		Type string
		CastToGoType string // кастомная функция по приведению типа из битрикса к типу для записи в постгрес
	}
)

// печатаем условие проверки, что результат запроса не пустой. Есть разница в проверке в зависимости от получаемых данных
// дефолтно это массив, поэтому просто проверяем длину массива
// но может быть и массив внутри другого массива, тогда нам нужно сперва проверить что первый массив не пустой и потом уже проверять длину первого элемента
func (db DocIntegrationsBitrix) PrintCheckResultIsEmpty() string {
	if len(db.Result.PathStr) > 0 {
		//res.Result[0].List
		if strings.HasPrefix(db.Result.PathStr, "Result[0]") {
			return fmt.Sprintf("len(res.Result) == 0 || len(res.%s) == 0", db.Result.PathStr)
		}
		return fmt.Sprintf("len(res.%s) == 0", db.Result.PathStr)
	}
	return "len(res.Result) == 0"
}
