package webServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"[[.Config.LocalProjectPath]]/pg"
	"[[.Config.LocalProjectPath]]/types"
	"[[.Config.LocalProjectPath]]/utils"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
	"time"
)

type (
	PgMethod struct {
		Title      string
		Roles      []string
		Cache      PgMethodCache
		BeforeHook func(*gin.Context, interface{}) error
	}
	pgFuncCacheType struct {
		Data        []byte
		ExpiredTime time.Time
	}
	PgMethodCache interface {
		Duration() time.Duration                                        // время кэширования в секундах
		GetKey(userId, role string, params interface{}) (string, error) // формирование ключа
	}
)

var (
	pgFuncCache = map[string]pgFuncCacheType{}
	pgFuncList  = []PgMethod{
		PgMethod{"user_update", []string{"admin",[[ArrayStringJoin .Config.User.Roles.UserUpdate ]]}, nil, nil},
		PgMethod{"user_list", []string{[[ArrayStringJoin .Config.User.Roles.UserList ]]}, nil, nil},
		PgMethod{"user_get_by_id", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"user_get_by_id_for_ui", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"current_user_update", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"current_user_get_auth_providers", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"message_list", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"message_update", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"message_mark_as_read", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"task_list", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"task_update", []string{"admin"}, nil, BeforeHookAddUserId},
		PgMethod{"task_get_by_id", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"task_type_list", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"task_list_for_user", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"task_action_to_finished", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"task_type_update", []string{"admin"}, nil, BeforeHookAddUserId},
		PgMethod{"task_type_get_by_id", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"chat_update", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"chat_get_by_id", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"chat_for_table_id", []string{}, nil, BeforeHookAddUserId},
		PgMethod{"chat_message_update", []string{}, nil, BeforeHookAddUserId},
		[[.PrintApiCallPgFuncMethods]]
	}
)

func apiCallPgFunc(c *gin.Context) {

	var jsonParam JsonParamType
	if v, ok := c.Get(utils.ContextJsonParam); !ok {
		utils.HttpError(c, http.StatusMethodNotAllowed, "missed params")
		return
	} else {
		jsonParam = v.(JsonParamType)
	}
	if len(jsonParam.Method) == 0 {
		utils.HttpError(c, http.StatusMethodNotAllowed, "missed method name")
		return
	}

	// проверяем что метод из списка разрешенных для вызова через api
	isCorrectMethod := false
	isAllowedMethod := false
	// создаем переменную для хранении информации о кэше
	cacheResult := struct {
		Data           []byte
		NewExpiredTime time.Time
		Key            string
	}{}

	u, _ := c.Get(utils.GinContextUser)
	user := u.(*types.User)

	for _, v := range pgFuncList {
		if v.Title == jsonParam.Method {
			isCorrectMethod = true
			// если роли для данного метода не указаны, то метод автоматически разрешен
			if len(v.Roles) == 0 {
				isAllowedMethod = true
			}
			// проверяем роль
			for _, role := range v.Roles {
				for _, userRole := range user.Role {
					if userRole == role {
						isAllowedMethod = true
						break
					}
				}
			}
			// проверяем есть ли beforeHook
			if v.BeforeHook != nil {
				err := v.BeforeHook(c, jsonParam)
				if err != nil {
					utils.HttpError(c, http.StatusMethodNotAllowed, err.Error())
					return
				}
			}
			if v.Cache != nil {
				key, err := v.Cache.GetKey(user.IdString(), user.GetRoleAsString(), jsonParam.Params)
				if err == nil {
					// в случае ошибки не формируем ключ для кэширования. Детали об ошибки выводим внутри метода.
					// заполняем информацию, необходимую для кэширования
					cacheResult.Key = key
					cacheResult.NewExpiredTime = time.Now().Add(v.Cache.Duration() * time.Second)
					if res, ok := pgFuncCache[key]; ok {
						// проверяем не истекло ли время кэша
						if res.ExpiredTime.After(time.Now()) {
							cacheResult.Data = res.Data
						}
					}
				}
			}
			break
		}
	}
	// если метода нет в списке то выходим
	if !isCorrectMethod {
		utils.HttpError(c, http.StatusMethodNotAllowed, "not allowed method: "+jsonParam.Method)
		return
	}

	// если метода нет в списке то выходим
	if !isAllowedMethod {
		utils.HttpError(c, http.StatusMethodNotAllowed, "for this role not allowed method: "+jsonParam.Method)
		return
	}

	// запрос к postgres, если нет данных из кэша
	queryRes := []byte("")
	if len(cacheResult.Data) == 0 {
		var err error
		queryRes, err = callPgFuncToJson(jsonParam.Method, jsonParam.Params)
		if err != nil {
			utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("%s", err))
			return
		}
		// в случае если указан ключ для кэширования, сохраняем полученные данные из базы в кэш
		if len(cacheResult.Key) > 0 {
			pgFuncCache[cacheResult.Key] = pgFuncCacheType{queryRes, cacheResult.NewExpiredTime}
		}
	} else {
		queryRes = cacheResult.Data
	}

	//// логгирование обращений к сайту
	//{
	//	u, _ := c.Get(common.GinContextUser)
	//	user := u.(*model.User)
	//	common.GraylogConfig.L().Infom(map[string]interface{}{"web_api_call": jsonParam.Method, "userId": user.Id, "username": user.Fullname}, fmt.Sprintf("%s. Is from cache: %v", jsonParam.Method, len(cacheResult.Data)))
	//
	//}

	// перекладываем ответ из jsonParam в jsonParam
	c.JSON(http.StatusOK, gin.H{
		"ok":        gjson.Get(fmt.Sprintf("%s", queryRes), "ok").Bool(),
		"result":    gjson.Get(fmt.Sprintf("%s", queryRes), "result").Value(),
		"message":   gjson.Get(fmt.Sprintf("%s", queryRes), "message").Value(),
		"meta_info": gjson.Get(fmt.Sprintf("%s", queryRes), "meta_info").Value(),
	})
}

func callPgFuncToJson(funcName string, jsonMap interface{}) (queryRes []byte, err error) {
	var queryStr string
	var jsonStr []byte

	// по разному обрабатываем параметры запроса в зависмости от типа: может быть строка, а может быть map
	switch v := jsonMap.(type) {
	case string:
		v = strings.Replace(v, "'", "''", -1)
		jsonStr = []byte(v)
	case map[string]interface{}:
		if jsonStr, err = json.Marshal(v); err != nil {
			return
		}
		jsonStr = []byte(strings.Replace(string(jsonStr), "'", "''", -1))
	default:
		jsonStr = []byte{}
	}

	if len(jsonStr) > 0 {
		queryStr = fmt.Sprintf("select * from %s('%s')", funcName, jsonStr)
	} else {
		queryStr = fmt.Sprintf("select * from %s()", funcName)
	}

	//println("!!! callPgFuncToJson queryStr:", queryStr)

	err = pg.Pg.QueryRow(queryStr).Scan(&queryRes)
	if err != nil {
		return
	}

	return
}

func BeforeHookAddUserId(c *gin.Context, p interface{}) error {
	jsonParam, ok := p.(JsonParamType)
	if !ok {
		return errors.New(fmt.Sprintf("BeforeHookAddUserId wrong type assertion %s not JsonParamType", p))
	}
	if jsonParam.Params == nil {
		return errors.New(fmt.Sprintf("missed params in method %s", jsonParam.Method))
	}
	if u, ok := c.Get(utils.GinContextUser); ok {
		user := u.(*types.User)
		m, ok := jsonParam.Params.(map[string]interface{})
		if !ok {
			return errors.New(fmt.Sprintf("BeforeHookAddUserId wrong type assertion %s not map[string]interface{}", jsonParam.Params))
		}
		m["user_id"] = user.Id
	}
	return nil
}
