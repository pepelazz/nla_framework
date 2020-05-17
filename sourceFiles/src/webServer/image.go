package webServer

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	"github.com/oklog/ulid"
	"github.com/pepelazz/projectGenerator/utils"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const IMAGE_DIR = "../image"
const STAT_IMAGE_PATH = "/stat-img"

func uploadImage(c *gin.Context) {
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
	// извлекаем минимальную ширину фото для сжатия. Если 0, то не сжимаем
	width := 0
	if widthStr, ok := c.GetPostForm("width"); ok {
		w, err := strconv.Atoi(widthStr)
		if err == nil {
			width = w
		}
	}
	// извлекаем параметры для crop. Например, 300x400
	crop := []int{}
	if cropStr, ok := c.GetPostForm("crop"); ok {
		widthStr := strings.Split(cropStr, "x")
		// перекладываем []string -> []int
		if len(widthStr) == 2 {
			for _, v := range widthStr {
				w, err := strconv.Atoi(v)
				if err == nil {
					crop = append(crop, w)
				}
			}
		}
	}
	path := fmt.Sprintf("%s/%s/%s", IMAGE_DIR, tableName, tableId)
	saveImage(c, path, "", width, crop)
}

func saveImage(c *gin.Context, path, filePrefix string, width int, crop []int) {
	// извлекаем файл из парамeтров post запроса
	form, _ := c.MultipartForm()
	var fileName string

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
	}
	// извлекаем содержание присланного файла по названию файла
	file, _, err := c.Request.FormFile(fileName)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadFile c.Request.FormFile error: %s", err.Error()))
		return
	}
	defer file.Close()

	// извлекаем расширение файла
	imgExt := "jpeg"
	contentType, _ := c.GetPostForm("Content-Type")
	if len(contentType) > 0 {
		arr := strings.Split(contentType, "/")
		if len(arr) > 1 {
			imgExt = arr[1]
		}
	}

	// перекодируем файл в картинку
	var img image.Image
	switch imgExt {
	case "jpeg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	case "gif":
		img, err = gif.Decode(file)
	default:
		err = errors.New("Unsupported file type")
	}
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage jpeg.Decode error: %s", err.Error()))
		return
	}
	// создаем директорию, если еще не создана
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage os.MkdirAll error: %s", err))
		return
	}

	// открываем файл для сохранения картинки
	fullFileName := fmt.Sprintf("%s%s.%s", filePrefix, randomFilename(), imgExt)
	fileOnDisk, err := os.Create(fmt.Sprintf("%s/%s", path, fullFileName))
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage os.Create err: %s", err))
		return
	}
	defer fileOnDisk.Close()

	// если необходимо обрезать
	if crop != nil && len(crop) == 2 {
		analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
		topCrop, _ := analyzer.FindBestCrop(img, crop[0], crop[1])
		type SubImager interface {
			SubImage(r image.Rectangle) image.Image
		}
		img = img.(SubImager).SubImage(topCrop)
	}

	// сжатие размеров картинки до минимума - 500 или фактический размер
	imgWidth := uint(utils.MinInt(width, img.Bounds().Max.X))
	resizedImg := resize.Resize(imgWidth, 0, img, resize.Lanczos3)

	// сохранение файла
	err = jpeg.Encode(fileOnDisk, resizedImg, nil)
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, fmt.Sprintf("uploadImage jpeg.Encode err: %s", err))
		return
	}

	// возвращаем ссылку на файл
	utils.HttpSuccess(c, map[string]string{"file": fmt.Sprintf("%s/%s", strings.Replace(path, IMAGE_DIR, STAT_IMAGE_PATH, 1), fullFileName)})
}

// загрузка аватарки
func uploadProfileImage(c *gin.Context)  {
	if userId, ok := utils.ExtractUserIdString(c); ok {
		path := fmt.Sprintf("%s/profile", IMAGE_DIR)
		prefix := fmt.Sprintf("id_%s_", userId)
		saveImage(c, path, prefix, 200, []int{200, 200})
	}
}

// генерим случаный uid для названия файла
func randomFilename() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return strings.ToLower(fmt.Sprintf("%v", ulid.MustNew(ulid.Timestamp(t), entropy)))
}
