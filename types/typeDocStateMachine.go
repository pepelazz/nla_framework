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
		IconSrc    string
	}

	DocSmAction struct {
		From              string
		To                string
		Label             string
		UpdateFlds        []FldType              // поля, которые заполняются при смене стейта
		CopyToHistoryFlds []FldType              // поля, значения которых копируются в историю изменений при смене статуса
		Conditions        []DocSmActionCondition // условия выполнения экшена
	}

	DocSmActionCondition struct {
		SqlText string
	}
)
