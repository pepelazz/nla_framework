package types

type (
	// State machine
	DocSm struct {
		States []*DocSmState
	}

	DocSmState struct {
		Title      string
		TitleRu    string
		Actions    []DocSmAction
		UpdateFlds []FldType // поля, которые можно редактировать в этом стейте
	}

	DocSmAction struct {
		From              string
		To                string
		UpdateFlds        []FldType // поля, которые заполняются при смене стейта
		CopyToHistoryFlds []FldType // поля, значения которых копируются в историю изменений при смене статуса
	}
)
