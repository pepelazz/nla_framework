package webServer

type JsonParamType struct {
	AuthToken string      `json:"auth_token"`
	Method    string      `json:"method"`
	Params    interface{} `json:"params"`
}


