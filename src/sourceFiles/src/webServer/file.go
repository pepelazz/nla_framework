package webServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pepelazz/projectGenerator/src/pg"
	"github.com/pepelazz/projectGenerator/src/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const FILE_DIR = "../uploaded_files"

func uploadFile(c *gin.Context) {
	// извлекаем название таблицы и id записи, к которой крепится файл
	tableName, _ := c.GetPostForm("tableName")
	{
		if len(tableName) == 0 {
			utils.HttpError(c, http.StatusBadRequest, "missed tableName")
			return
		}
	}
	tableId, _ := c.GetPostForm("tableId")
	{
		if len(tableId) == 0 {
			utils.HttpError(c, http.StatusBadRequest, "missed tableId")
			return
		}
	}
	// TODO: проверка, может ли пользователь с этой ролью прикреплять файлы к данным таблицам

	path := fmt.Sprintf("%s/%s/%s", FILE_DIR, tableName, tableId)
	err := os.MkdirAll(path, os.ModePerm) // создаем директорию, если еще не создана
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadFile os.MkdirAll error: %s", err))
		return
	}

	// извлекаем файл из парамeтров post запроса
	form, _ := c.MultipartForm()
	var fileName, fileExt string
	var fileSize int

	if len(form.File) == 0 {
		utils.HttpError(c, http.StatusBadRequest, "list of files is empty")
		return
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
		utils.HttpError(c, http.StatusBadRequest, "wrong file extansion")
		return
	}
	// извлекаем содержание присланного файла по названию файла
	file, _, err := c.Request.FormFile(fileName)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadFile c.Request.FormFile error: %s", err.Error()))
		return
	}
	defer file.Close()

	// читаем содержание присланного файл в []byte
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(fileBytes) > 10000000 {
		utils.HttpError(c, http.StatusBadRequest, "FILE_SIZE_LIMIT: 10Mb")
		return
	}
	fileSize = len(fileBytes) / 1000

	// создаем запись в БД с информацией о файле и получаем token для имени файла
	jsonStr, _ := json.Marshal(map[string]interface{}{"id": -1, "filename": fileName, "ext": fileExt, "table_name": tableName, "table_id": tableId, "size": fileSize})
	fileUpdateRes := struct {
		Token string `json:"token"`
	}{}
	err = pg.CallPgFunc("file_update", jsonStr, &fileUpdateRes, nil)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	// открываем файл для сохранения
	fileOnDisk, err := os.Create(fmt.Sprintf("%s/%s", path, fileName))
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage os.Create err: %s", err))
		return
	}
	defer fileOnDisk.Close()

	_, err = fileOnDisk.Write(fileBytes)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	// возвращаем ссылку на файл
	utils.HttpSuccess(c, map[string]string{"filename": fileName, "ext": fileExt, "url": fmt.Sprintf("/api/file/%s", fileUpdateRes.Token)})
}

func downloadFile(c *gin.Context)  {
	// извлекаем токен файла из параметро запроса
	fileToken := c.Param("fileToken")
	if len(fileToken) == 0 {
		utils.HttpError(c, http.StatusMethodNotAllowed, "missed file token")
		return
	}
	userId, _ := c.Get("user_id")
	path, err := getFilePathByToken(userId, fileToken)
	if err != nil {
		utils.HttpError(c, http.StatusMethodNotAllowed, err.Error())
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=test")
	c.Header("Content-Type", "octet-stream")
	c.File(path)
}

func deleteFile(c *gin.Context)  {
	// извлекаем токен файла из параметро запроса
	fileToken := c.Param("fileToken")
	if len(fileToken) == 0 {
		utils.HttpError(c, http.StatusMethodNotAllowed, "missed file token")
		return
	}
	userId, _ := c.Get("user_id")
	path, err := getFilePathByToken(userId, fileToken)
	if err != nil {
		utils.HttpError(c, http.StatusMethodNotAllowed, err.Error())
		return
	}
	err = os.Remove(path)
	if err != nil {
		utils.HttpError(c, http.StatusMethodNotAllowed, err.Error())
		return
	}
	utils.HttpSuccess(c, "ok")
}

func getFilePathByToken(userId interface{}, fileToken string) (string, error) {
	jsonStr, _ := json.Marshal(map[string]interface{}{"user_id": userId, "token": fileToken})
	res := struct {
		Id int64 `json:"id"`
		Filename string `json:"filename"`
		TableName string `json:"table_name"`
		TableId int64 `json:"table_id"`
	}{}
	err := pg.CallPgFunc("file_get_by_token",jsonStr, &res, nil)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/%s/%v/%s", FILE_DIR, res.TableName, res.TableId, res.Filename)
	if !fileExists(path) {
		return "", errors.New("file not found on disk")
	}
	return path, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}