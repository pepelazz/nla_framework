package utils

import (
	"bytes"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/pepelazz/projectGenerator/types"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"crypto/rand"
	"errors"
	"log"
	"strings"
)

const (
	GinContextUser              = "user"
	GinContextUserId            = "user_id"
	ContextJsonParam            = "jsonParam"         //параметры в web запросах
	ContextJsonParamFldParam    = "jsonParamFldParam" //поле params в параметры в web запросах
	GinContextGetRequestQueryId = "getRequestQueryId"
	GinContextAppAuth           = "app_auth"
	GinContextAppAuthId         = "app_auth_id"
)

var (
	webServerConfig types.WebServer
	emailConfig     types.EmailConfig
)

func SetWebServerConfig(config types.WebServer) {
	webServerConfig = config
}

func SetEmailConfig(config types.EmailConfig) {
	emailConfig = config
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func HttpError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"ok":      false,
		"message": message,
	})
	c.Abort()
}

func HttpSuccess(c *gin.Context, res interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": res,
	})
}

func CheckErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Panic(msg string) {
	log.Fatalf("%s", msg)
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// функция извлечения json параметров, переданных строкой
func ExtractPostReqParams(c *gin.Context, res interface{}) bool {

	v, ok := c.Get(ContextJsonParamFldParam)
	if !ok {
		HttpError(c, http.StatusMethodNotAllowed, "missed params")
		return false
	}
	paramStr, ok := v.(string)
	if !ok {
		HttpError(c, http.StatusMethodNotAllowed, fmt.Sprintf("extractPostReqParams wrong type assertion %s not string", v))
		return false
	}

	err := json.Unmarshal([]byte(paramStr), &res)
	if err != nil {
		HttpError(c, http.StatusMethodNotAllowed, fmt.Sprintf("extractPostReqParams json.Unmarshal %s params: %s", err.Error(), paramStr))
		return false
	}

	return true
}

// функция извлечения json параметров, переданных строкой
func ExtractPostReqParamsMap(c *gin.Context) (map[string]interface{}, bool) {

	v, ok := c.Get(ContextJsonParamFldParam)
	if !ok {
		HttpError(c, http.StatusMethodNotAllowed, "missed params")
		return nil, false
	}

	paramStr, ok := v.(string)
	if !ok {
		HttpError(c, http.StatusBadRequest, fmt.Sprintf("extractPostReqParamsMap wrong type assertion %s not string", v))
		return nil, false
	}

	mapRes := map[string]interface{}{}

	err := json.Unmarshal([]byte(paramStr), &mapRes)
	if err != nil {
		HttpError(c, http.StatusBadRequest, fmt.Sprintf("extractPostReqParamsMap json.Unmarshal %s", err))
		return nil, false
	}

	return mapRes, true
}

// функция извлечения из контекста запроса userId в виде строки
func ExtractUserIdString(c *gin.Context) (string, bool) {
	userId, ok := c.Get(GinContextUserId)
	if !ok {
		HttpError(c, http.StatusBadRequest, "not found user")
		return "", false
	}

	var userIdStr string

	switch v := userId.(type) {
	case string:
		userIdStr = v
	case int:
		userIdStr = strconv.Itoa(v)
	case int64:
		userIdStr = strconv.FormatInt(v, 10)
	}
	if len(userIdStr) > 0 {
		return userIdStr, true
	} else {
		return "", false
	}
}

// функция извлечения из контекста запроса userId в виде строки
func ExtractUserIdInt64(c *gin.Context) (int64, bool) {
	userId, ok := c.Get(GinContextUserId)
	if !ok {
		HttpError(c, http.StatusBadRequest, "not found user")
		return 0, false
	}

	var userIdInt64 int64

	switch v := userId.(type) {
	case string:
		var err error
		userIdInt64, err = strconv.ParseInt(v, 0, 64)
		if err != nil {
			return 0, false
		}
	case int:
		userIdInt64 = int64(v)
	case int64:
		userIdInt64 = v
	}
	return userIdInt64, true
}

func ExtractJsonParam(c *gin.Context, res interface{}) error {
	jsonParam, _ := c.Get(ContextJsonParamFldParam)
	paramstr, ok := jsonParam.(string)
	errMsg := ""
	if !ok {
		errMsg = "json params convert error - need string (JSON.stringify)"
		HttpError(c, http.StatusBadRequest, errMsg)
		return errors.New(errMsg)
	}
	err := json.Unmarshal([]byte(paramstr), &res)
	if err != nil {
		errMsg = fmt.Sprintf("json.Unmarshal err: %s\n", err)
		HttpError(c, http.StatusBadRequest, errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func RandToken(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func GetJsonByUrl(url string, res interface{}) error {

	httpRes, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	return nil
}

func ReadUploadedFile(c *gin.Context, exts []string) (multipart.File, error) {
	// извлекаем файл из парамeтров post запроса
	form, _ := c.MultipartForm()
	var fileName, fileExt string

	if len(form.File) == 0 {
		return nil, errors.New("list of files is empty")
	}
	// берем первое имя файла из присланного списка
	for key := range form.File {
		if len(fileName) > 0 {
			continue
		}
		fileName = key
		// извлекаем расширение файла
		arr := strings.Split(fileName, ".")
		if len(arr) > 1 {
			fileExt = arr[len(arr)-1]
		}
	}
	if len(fileExt) == 0 {
		return nil, errors.New("wrong file extansion")
	}
	isExtTrue := false
	for _, v := range exts {
		if fileExt == v {
			isExtTrue = true
		}
	}
	if !isExtTrue {
		return nil, errors.New(fmt.Sprintf("file extansion must be %s", exts))
	}
	// извлекаем содержание присланного файла по названию файла
	file, _, err := c.Request.FormFile(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("uploadFile c.Request.FormFile error: %s", err.Error()))
	}
	return file, nil
}
