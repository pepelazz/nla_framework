package webServer

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/pepelazz/projectGenerator/graylog"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"net/http"
)

func logToGraylog(c *gin.Context) {

	var user *types.User

	if u, ok := c.Get(utils.GinContextUser); ok {
		user = u.(*types.User)
	}

	logMsg := struct {
		Type string                 `json:"type"`
		Code string                 `json:"code"`
		Msg  map[string]interface{} `json:"msg"`
	}{}

	// извлекаем параметр page из json
	// вариант когда пользователь авторизован и параметры уже извлечены и находятся в ContextJsonParamFldParam
	if user != nil {
		if ok := utils.ExtractPostReqParams(c, &logMsg); !ok {
			return
		}
	} else {
		// вариант когда пользователь неавторизован и передан только параметры
		if err := c.BindJSON(&logMsg); err != nil {
			utils.HttpError(c, http.StatusOK, "post json params error:"+fmt.Sprintf("%s", err))
			return
		}
	}

	msgAttrs := map[string]interface{}{"type": logMsg.Code}
	if user != nil {
		msgAttrs["user_id"] = user.Id
		msgAttrs["fullname"] = user.Fullname
	}
	for k, v := range logMsg.Msg {
		msgAttrs[k] = v
	}

	if logMsg.Type == "error" {
		err := graylog.Graylog.L().Errm(msgAttrs, fmt.Sprintf("%s", logMsg.Msg))
		if err != nil {
			fmt.Printf("logToGraylog send to graylog error %s\n", err)
		}
	} else {
		err := graylog.Graylog.L().Infom(msgAttrs, fmt.Sprintf("%s", logMsg.Msg))
		if err != nil {
			fmt.Printf("logToGraylog send to graylog error %s\n", err)
		}
	}
	utils.HttpSuccess(c, "ok")
}
