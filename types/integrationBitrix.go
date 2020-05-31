package types

type (
	BitrixFld struct {
		Name string
		Type string
		CastToGoType string // кастомная функция по приведению типа из битрикса к типу для записи в постгрес
	}
)

