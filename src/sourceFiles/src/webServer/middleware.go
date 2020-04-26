package webServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pepelazz/projectGenerator/src/cacheUtil"
	"github.com/pepelazz/projectGenerator/src/pg"
	"github.com/pepelazz/projectGenerator/src/types"
	"github.com/pepelazz/projectGenerator/src/utils"
	"net/http"
	"strings"
	"time"
)

// проверка токена авторизации
func authRequired(c *gin.Context) {
	//println("in AuthRequired", c.Request.Method)
	// ищем auth_token
	// вначале ищем в header'е запроса
	authToken := c.Request.Header.Get("Auth-token")

	contentType := c.Request.Header.Get("Content-Type")

	// вариант обработки POST запроса
	if c.Request.Method == "POST" {
		// если это не загрузка файлов (тогда в запросе нет json параметров)
		if !strings.Contains(contentType, "multipart/form-data") {
			var json JsonParamType
			//println("AuthRequired Auth-token:", c.Request.Header.Get("Auth-token"))
			// извлекаем json-параметры запроса
			if err := c.BindJSON(&json); err != nil {
				//println("AuthRequired step1.1", fmt.Sprintf("%s", err))
				utils.HttpError(c, http.StatusOK, "post json params error:"+fmt.Sprintf("%s", err))
				return
			}
			// если не находим, то смотрим в json-параметрах запроса
			if len(authToken) == 0 {
				authToken = json.AuthToken
			}
			c.Set(utils.ContextJsonParam, json)
			c.Set(utils.ContextJsonParamFldParam, json.Params)
		}
	}

	// в случае GET запроса смотрим токен в query параметрах
	if c.Request.Method == "GET" {
		if len(authToken) == 0 {
			authToken = c.Query("authToken")
		}
		// точно используется при вебсокет соединении, когда отправляется первоначальный запрос
		c.Set(utils.GinContextGetRequestQueryId, c.Query("id"))
	}
	// если токен не найден, то возвращаем ошибку
	if len(authToken) == 0 {
		utils.HttpError(c, http.StatusOK, "missed auth_token")
		return
	}

	user, err := userFindByAuthToken(authToken)
	if err != nil {
		utils.HttpError(c, http.StatusOK, fmt.Sprintf("%s", err))
		return
	}
	c.Set(utils.GinContextUser, user)
	c.Set(utils.GinContextUserId, user.Id)
}

// LiberalCORS is a very allowing CORS middleware.
func LiberalCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	if c.Request.Method == "OPTIONS" {
		if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
			c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
		}
		c.AbortWithStatus(http.StatusOK)
	}
}

func userFindByAuthToken(token string) (user *types.User, err error) {
	// ищем пользователя в кэше
	userIntreface, _ := cacheUtil.GoCacheGet(cacheUtil.GetCacheKeyUserToken(token))
	if userIntreface != nil {
		var ok bool
		user, ok = userIntreface.(*types.User)
		if ok {
			return
		}
	}
	jsonStr, _ := json.Marshal(map[string]interface{}{"auth_token": token})
	err = pg.CallPgFunc("user_get_by_auth_token", jsonStr, &user, nil)
	// записываем в кэш
	if err == nil {
		cacheUtil.GoCacheSet(cacheUtil.GetCacheKeyUserToken(token), user, time.Minute*1)
	}
	return
}