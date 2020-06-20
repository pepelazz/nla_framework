package webServer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"[[.Config.LocalProjectPath]]/pg"
	"[[.Config.LocalProjectPath]]/types"
	"[[.Config.LocalProjectPath]]/utils"
	"io"
	"net/http"
	"strings"
)

func telegramAuth(config types.TelegramConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		tgUser := struct {
			AuthDate       int64  `json:"auth_date"`
			FirstName      string `json:"first_name"`
			LastName       string `json:"last_name"`
			PhotoUrl       string `json:"photo_url"`
			Id             int64  `json:"id"`
			Username       string `json:"username"`
			Hash           string `json:"hash"`
			UserId 		   int64  `json:"user_id"`
		}{}

		if err := utils.ExtractJsonParam(c, &tgUser); err != nil {
			utils.HttpError(c, http.StatusOK, "post json params error:"+fmt.Sprintf("%s", err))
			return
		}

		// https://core.telegram.org/widgets/login#widget-configuration
		// формируем строку для проверки
		checkStr := fmt.Sprintf("auth_date=%v", tgUser.AuthDate)
		if len(tgUser.FirstName) > 0 {
			checkStr = checkStr + fmt.Sprintf("\nfirst_name=%s", tgUser.FirstName)
		}
		checkStr = checkStr + fmt.Sprintf("\nid=%v", tgUser.Id)
		if len(tgUser.LastName) > 0 {
			checkStr = checkStr + fmt.Sprintf("\nlast_name=%s", tgUser.LastName)
		}
		if len(tgUser.PhotoUrl) > 0 {
			checkStr = checkStr + fmt.Sprintf("\nphoto_url=%s", tgUser.PhotoUrl)
		}
		if len(tgUser.Username) > 0 {
			checkStr = checkStr + fmt.Sprintf("\nusername=%s", tgUser.Username)
		}

		// формируем sha256 с ключем - токен бота
		sha256hash := sha256.New()
		io.WriteString(sha256hash, config.Token)
		hmachash := hmac.New(sha256.New, sha256hash.Sum(nil))
		io.WriteString(hmachash, checkStr)
		// проверяем что совпадает с присланным хэшом
		ss := hex.EncodeToString(hmachash.Sum(nil))

		isAuth := strings.EqualFold(ss, tgUser.Hash)

		// если hash не совпадает, то возвращаем ошибку
		if !isAuth {
			utils.HttpError(c, 400, "not auth: hash invalid")
			return
		}
		// добавляем user_id текущего пользователя
		userId, _ := utils.ExtractUserIdInt64(c)
		tgUser.UserId = userId
		jsonStr, _ := json.Marshal(tgUser)
		err := pg.CallPgFunc("user_telegram_auth", jsonStr, nil, nil)
		if err != nil {
			utils.HttpError(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.HttpSuccess(c, "ok")
	}
}

