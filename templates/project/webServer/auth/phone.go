package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pepelazz/indivigo/src/pg"
	"github.com/pepelazz/indivigo/src/types"
	"github.com/pepelazz/indivigo/src/utils"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type (
	PasswordRecoverPhoneToken struct {
		Phone       string `json:"phone"`
		Token       string `json:"token"`
		ExpiredTime time.Time
	}
)

// map для хранения токенов, для отправленных писем со сбросом пароля. Данная структура сбрасывается при перезапуске приложения
var passwordRecoverPhoneTokenMap = map[string]PasswordRecoverPhoneToken{}

func PhoneAuth(c *gin.Context) {

	reqParams := struct {
		Params struct {
			Login      string                 `json:"login"`
			Password   string                 `json:"password"`
			LastName   string                 `json:"last_name"`
			FirstName  string                 `json:"first_name"`
			Options    map[string]interface{} `json:"options"`
			IsRegister bool                   `json:"is_register"` // флаг, которым различаем регистрацию нового пользователя и авторизацию существующего
		} `json:"params"`
	}{}

	type UserRegisterData struct {
		Login     string                 `json:"login"`
		Password  string                 `json:"password"`
		LastName  string                 `json:"last_name"`
		FirstName string                 `json:"first_name"`
		AuthToken string                 `json:"auth_token"`
		Token     string                 `json:"token"`
		Phone     string                 `json:"phone"`
		Options   map[string]interface{} `json:"options"`
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
	userRegData.Phone = reqParams.Params.Login

	if reqParams.Params.IsRegister {
		// вариант регистрации нового пользователя
		userRegData.Password = string(hashedPassword)
		userRegData.Token = strconv.Itoa(100000 + rand.Intn(999999-100000))
		userRegData.LastName = reqParams.Params.LastName
		userRegData.FirstName = reqParams.Params.FirstName
		userRegData.Options = reqParams.Params.Options

		jsonStr, err := json.Marshal(userRegData)
		err = pg.CallPgFunc("user_temp_phone_auth_create", jsonStr, nil, nil)
		if err != nil {
			utils.HttpError(c, http.StatusOK, "pg call user_temp_phone_auth_create err:"+fmt.Sprintf("%s", err))
			return
		}
		// отправляем sms
		//TODO: обращение к сервису по отправке смс
		//href := fmt.Sprintf("%s/check_user_email?t=%v", webServerConfig.Url, userRegData.Token)
		//err = utils.EmailSendRegistrationConfirm(userRegData.Email, href)
		//if err != nil {
		//	fmt.Printf("utils.EmailSendRegistrationConfirm: %s", err)
		//}
		// в независимости от результатов отправки письма отправляем, что данный этап регистрации успешно пройден
		utils.HttpSuccess(c, nil)
	} else {

		// вариант авторизации существующего пользователя
		user := types.User{}
		jsonStr, err := json.Marshal(userRegData)
		err = pg.CallPgFunc("user_get_by_phone_with_password", jsonStr, &user, nil)
		if err != nil {
			utils.HttpError(c, http.StatusOK, fmt.Sprintf("pg call user_get_by_phone_with_password err %s", err))
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

// подтерждение телефона при регистрации нового пользователя
func CheckSmsCode(c *gin.Context) {
	type QueryData struct {
		Params struct {
			Phone string `json:"phone"`
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
	jsonStr, err := json.Marshal(queryData.Params)
	err = pg.CallPgFunc("user_temp_phone_auth_check_sms_code", jsonStr, &user, nil)
	if err != nil {
		if len(err.Error()) > 0 {
			utils.HttpError(c, http.StatusOK, "pg call user_temp_phone_auth_check_sms_code err:"+fmt.Sprintf("%s", err))
		} else {
			utils.HttpError(c, http.StatusOK, "")
		}
		return
	}

	utils.HttpSuccess(c, user)
}

// функция начала сброса пароля. Создаем пару phone-токен и отправляем пользователю sms с токеном
func PhoneAuthStartRecoverPassword(c *gin.Context) {
	type QueryData struct {
		Params struct {
			Phone string `json:"phone"`
			Token string `json:"token"`
		} `json:"params"`
	}

	var queryData QueryData
	// извлекаем json-параметры запроса
	if err := c.BindJSON(&queryData); err != nil {
		utils.HttpError(c, http.StatusOK, fmt.Sprintf("post json params error: %s", err))
		return
	}
	// проверяем что такой пользователь с таким phone есть в базе
	user := types.User{}
	jsonStr, err := json.Marshal(queryData.Params)
	err = pg.CallPgFunc("user_get_by_phone_with_password", jsonStr, &user, nil)
	if err != nil {
		utils.HttpError(c, http.StatusOK, "pg call user_get_by_phone_with_password err:"+fmt.Sprintf("%s", err))
		return
	}
	if user.Id == 0 {
		utils.HttpError(c, http.StatusOK, "user not found")
		return
	}

	token := strconv.Itoa(100000 + rand.Intn(999999-100000))
	phone := queryData.Params.Phone
	fmt.Printf("token: %s\n", token)
	// добавляем пару токен-email а коллекцию
	passwordRecoverPhoneTokenMap[phone] = PasswordRecoverPhoneToken{phone, token, time.Now().Add(1 * time.Minute)}
	//удаляем просроченные токены
	for k, v := range passwordRecoverPhoneTokenMap {
		if time.Now().After(v.ExpiredTime) {
			delete(passwordRecoverPhoneTokenMap, k)
		}
	}

	//TODO: отправляем sms с кодом для восставноления пароля

	utils.HttpSuccess(c, nil)
}

// функция замены пароля.
func PhoneAuthRecoverPassword(c *gin.Context) {
	type QueryData struct {
		Params struct {
			Phone    string `json:"phone"`
			Token    string `json:"token"`
			Password string `json:"password"`
		} `json:"params"`
	}

	var queryData QueryData
	// извлекаем json-параметры запроса
	if err := c.BindJSON(&queryData); err != nil {
		utils.HttpError(c, http.StatusOK, "post json params error:"+fmt.Sprintf("%s", err))
		return
	}

	// проверяем что новый пароль не пустой
	if len(queryData.Params.Password) == 0 {
		utils.HttpError(c, http.StatusOK, "password is empty")
		return
	}
	// находим запись по номеру телефона
	if v, ok := passwordRecoverPhoneTokenMap[queryData.Params.Phone]; ok {
		if time.Now().After(v.ExpiredTime) {
			// токен просрочен
			utils.HttpError(c, http.StatusOK, "invalid token")
			return
		}
		// сравниваем код полученный от пользоваателя, с тем, который был отправлен ему по sms
		if v.Token != queryData.Params.Token {
			// токен не совпадает с тем, который был отправлен по sms
			utils.HttpError(c, http.StatusOK, "invalid token")
			return
		}
		// вариант валидного токена
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(queryData.Params.Password), 8)
		if err != nil {
			utils.HttpError(c, 400, fmt.Sprintf("bcrypt.GenerateFromPassword error:%s", err))
			return
		}
		// по phone находим пользователя
		user := types.User{}
		jsonStr, err := json.Marshal(v)
		err = pg.CallPgFunc("user_get_by_phone_with_password", jsonStr, &user, nil)
		if err != nil {
			utils.HttpError(c, http.StatusOK, "pg call user_get_by_email_with_password err:"+fmt.Sprintf("%s", err))
			return
		}
		if user.Id == 0 {
			utils.HttpError(c, http.StatusOK, "user not found")
			return
		}
		// сохраняем новый пароль в базу
		err = updateUserPassword(&user, string(hashedPassword), "phone")
		if err != nil {
			utils.HttpError(c, http.StatusOK, err.Error())
			return
		}
		// стираем токен из коллекции
		delete(passwordRecoverPhoneTokenMap, queryData.Params.Phone)
		utils.HttpSuccess(c, nil)
		return
	} else {
		// токен не найден в коллекции
		utils.HttpError(c, http.StatusOK, "invalid token")
		return
	}

	utils.HttpSuccess(c, nil)
}
