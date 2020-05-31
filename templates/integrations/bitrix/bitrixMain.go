package bitrix

import (
	"encoding/json"
	"[[.Config.LocalProjectPath]]/pg"
	"[[.Config.LocalProjectPath]]/types"
)

var (
	bitrixConfig     types.BitrixConfig
)

func SetBitrixConfig(config types.BitrixConfig) {
	bitrixConfig = config
}

func saveResultMsgToPg(userId string, msg string) error {
	jsonStr, _ := json.Marshal(map[string]interface{}{"id": -1, "user_id": userId, "title": msg})
	return pg.CallPgFunc("message_update", jsonStr, nil, nil)
}
