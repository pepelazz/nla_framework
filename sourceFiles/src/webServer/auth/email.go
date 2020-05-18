package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pepelazz/projectGenerator/pg"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type (
	PasswordRecoverEmailToken struct {
		Email       string `json:"email"`
		ExpiredTime time.Time
	}
)

// map для хранения токенов, для отправленных писем со сбросом пароля. Данная структура сбрасывается при перезапуске приложения
var passwordRecoverEmailTokenMap = map[string]PasswordRecoverEmailToken{}

func EmailAuth(c *gin.Context) {

	reqParams := struct {
		Params struct {
			Login      string `json:"login"`
			Password   string `json:"password"`
			LastName   string `json:"last_name"`
			FirstName  string `json:"first_name"`
			IsRegister bool   `json:"is_register"` // флаг, которым различаем регистрацию нового пользователя и авторизацию существующего
		} `json:"params"`
	}{}

	type UserRegisterData struct {
		Login      string `json:"login"`
		Password   string `json:"password"`
		LastName   string `json:"last_name"`
		FirstName  string `json:"first_name"`
		AuthToken  string `json:"auth_token"`
		Token      string `json:"token"`
		Email      string `json:"email"`
	}

	var userRegData UserRegisterData
	// извлекаем json-параметры запроса
	if err := c.BindJSON(&reqParams); err != nil {
		utils.HttpError(c, http.StatusOK, "post json params error:"+fmt.Sprintf("%s", err))
		return
	}
	// трансформируем пароль в hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqParams.Params.Password), 8)
	if err != nil {
		utils.HttpError(c, 400, fmt.Sprintf("bcrypt.GenerateFromPassword error:%s", err))
		return
	}
	userRegData.Login = reqParams.Params.Login
	userRegData.Email = reqParams.Params.Login

	if reqParams.Params.IsRegister {
		// вариант регистрации нового пользователя
		userRegData.Password = string(hashedPassword)
		userRegData.Token = xid.New().String()
		userRegData.LastName = reqParams.Params.LastName
		userRegData.FirstName = reqParams.Params.FirstName

		jsonStr, err := json.Marshal(userRegData);
		err = pg.CallPgFunc("user_temp_email_auth_create", jsonStr, nil, nil)
		if err != nil {
			utils.HttpError(c, http.StatusOK, "pg call user_temp_email_auth_create err:"+fmt.Sprintf("%s", err))
			return
		}
		// отправляем письмо
		href := fmt.Sprintf("%s/check_user_email?t=%v", webServerConfig.Url, userRegData.Token)
		err = utils.EmailSendRegistrationConfirm(userRegData.Email, href)
		if err != nil {
			fmt.Printf("utils.EmailSendRegistrationConfirm: %s", err)
		}
		// в независимости от результатов отправки письма отправляем, что данный этап регистрации успешно пройден
		utils.HttpSuccess(c, nil)
	} else {

		// вариант авторизации существующего пользователя
		user := types.User{}
		jsonStr, err := json.Marshal(userRegData)
		err = pg.CallPgFunc("user_get_by_email_with_password", jsonStr, &user, nil)
		if err != nil {
			utils.HttpError(c, http.StatusOK, fmt.Sprintf("pg call user_get_by_email_with_password err %s", err))
			return
		}
		// проверяем пароль
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqParams.Params.Password))
		if err != nil {
			utils.HttpError(c, http.StatusOK, "wrong password")
			return
		}
		// перед отправкой в данных пользователя стираем пароль
		user.Password = ""

		utils.HttpSuccess(c, user)
	}
}

// подтерждение email при регистрации нового пользователя (вариант POST метода)
func EmailAuthCheckUserEmail(c *gin.Context) {
	type QueryData struct {
		Params struct {
			Token string `json:"token"`
		} `json:"params"`
	}

	var queryData QueryData
	// извлекаем json-параметры запроса
	if err := c.BindJSON(&queryData); err != nil {
		utils.HttpError(c, http.StatusOK, fmt.Sprintf("post json params error: %s", err))
		return
	}

	user := types.User{}
	jsonStr, err := json.Marshal(queryData.Params);
	err = pg.CallPgFunc("user_temp_email_auth_check_token", jsonStr, &user, nil)
	if err != nil {
		if len(err.Error()) > 0 {
			utils.HttpError(c, http.StatusOK, "pg call user_temp_email_auth_check_token err:"+fmt.Sprintf("%s", err))
		} else {
			utils.HttpError(c, http.StatusOK, "")
		}
		return
	}

	utils.HttpSuccess(c, user)
}

// функция начала сброса пароля. Создаем пару email-токен и отправляем пользователю письмо со ссылкой для восстановления пароля
func EmailAuthStartRecoverPassword(c *gin.Context) {
	type QueryData struct {
		Params struct {
			Email string `json:"email"`
		} `json:"params"`
	}

	var queryData QueryData
	// извлекаем json-параметры запроса
	if err := c.BindJSON(&queryData); err != nil {
		utils.HttpError(c, http.StatusOK, fmt.Sprintf("post json params error: %s", err))
		return
	}
	// проверяем что такой пользователь с таким email есть в базе
	user := types.User{}
	jsonStr, err := json.Marshal(queryData.Params);
	err = pg.CallPgFunc("user_get_by_email_with_password", jsonStr, &user, nil)
	if err != nil {
		utils.HttpError(c, http.StatusOK, "pg call user_get_by_email_with_password err:"+fmt.Sprintf("%s", err))
		return
	}
	if user.Id == 0 {
		utils.HttpError(c, http.StatusOK, "user not found")
		return
	}

	token := xid.New().String()
	// добавляем пару токен-email а коллекцию
	passwordRecoverEmailTokenMap[token] = PasswordRecoverEmailToken{queryData.Params.Email, time.Now().Add(20 * time.Minute)}
	//удаляем просроченные токены
	for k, v := range passwordRecoverEmailTokenMap {
		if time.Now().After(v.ExpiredTime) {
			delete(passwordRecoverEmailTokenMap, k)
		}
	}

	//отправляем письмо со ссылкой для восставноления пароля
	href := fmt.Sprintf("%s/email_auth_recover_password?t=%v", webServerConfig.Url, token)
	err = utils.EmailSendChangePassword(queryData.Params.Email, href)
	if err != nil {
		utils.HttpError(c, http.StatusOK, fmt.Sprintf("EmailSendMessage: %s", err))
		return
	}

	utils.HttpSuccess(c, nil)
}

// функция замены пароля.
func EmailAuthRecoverPassword(c *gin.Context) {
	type QueryData struct {
		Params struct {
			Password     string `json:"password"`
			Token        string `json:"token"`
			IsTokenCheck bool   `json:"is_token_check"` // флаг для того чтобы отличать post запрос первого шага, когда пользователь открывает ссылку с токеном, присланную ему на email
		} `json:"params"`
	}

	var queryData QueryData
	// извлекаем json-параметры запроса
	if err := c.BindJSON(&queryData); err != nil {
		utils.HttpError(c, http.StatusOK, "post json params error:"+fmt.Sprintf("%s", err))
		return
	}
	if queryData.Params.IsTokenCheck {
		// первый шаг процесса восстановления пароля. Проверяем, что такой токен есть в коллекции passwordRecoverEmailTokenMap и он не протух
		if v, ok := passwordRecoverEmailTokenMap[queryData.Params.Token]; ok {
			if time.Now().After(v.ExpiredTime) {
				// токен просрочен
				utils.HttpError(c, http.StatusOK, "invalid token")
			}
			// вариант валидного токена
			utils.HttpSuccess(c, "ok")
		} else {
			// токен не найден в коллекции
			utils.HttpError(c, http.StatusOK, "invalid token")
		}
		return
	} else {
		// проверяем что новый пароль не пустой
		if len(queryData.Params.Password) == 0 {
			utils.HttpError(c, http.StatusOK, "password is empty")
			return
		}
		// находим запись по токену
		if v, ok := passwordRecoverEmailTokenMap[queryData.Params.Token]; ok {
			if time.Now().After(v.ExpiredTime) {
				// токен просрочен
				utils.HttpError(c, http.StatusOK, "invalid token")
			}
			// вариант валидного токена
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(queryData.Params.Password), 8)
			if err != nil {
				utils.HttpError(c, 400, fmt.Sprintf("bcrypt.GenerateFromPassword error:%s", err))
				return
			}
			// по email находим пользователя
			user := types.User{}
			jsonStr, err := json.Marshal(v);
			err = pg.CallPgFunc("user_get_by_email_with_password", jsonStr, &user, nil)
			if err != nil {
				utils.HttpError(c, http.StatusOK, "pg call user_get_by_email_with_password err:"+fmt.Sprintf("%s", err))
				return
			}
			if user.Id == 0 {
				utils.HttpError(c, http.StatusOK, "user not found")
				return
			}
			// сохраняем новый пароль в базу
			err = updateUserPassword(&user, string(hashedPassword))
			if err != nil {
				utils.HttpError(c, http.StatusOK, err.Error())
				return
			}
			// стираем токен из коллекции
			delete(passwordRecoverEmailTokenMap, queryData.Params.Token)
			utils.HttpSuccess(c, nil)
			return
		} else {
			// токен не найден в коллекции
			utils.HttpError(c, http.StatusOK, "invalid token")
		}
	}

	utils.HttpSuccess(c, nil)
}

func updateUserPassword(user *types.User, pwd string) (err error) {
	if user == nil {
		return errors.New("updateUserPassword user is nil")
	}
	user.Password = pwd

	jsonStr, _ := json.Marshal(map[string]interface{}{"id": user.Id, "password": pwd})

	err = pg.CallPgFunc("user_auth_update_password", jsonStr, &user, nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("UpdateUserPassword error: %s", err))
		return
	}

	return
}