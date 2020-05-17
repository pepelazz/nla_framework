package webServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/projectGenerator/utils"
	"github.com/pepelazz/projectGenerator/webServer/auth"
	"github.com/pepelazz/projectGenerator/sse"
	"net/http"
)

func StartWebServer(config types.Config) {
	r := gin.New()

	// передаем конфиги для модуля авторизации
	auth.SetWebServerConfig(config.WebServer)

	// вырубаем CORS
	r.Use(LiberalCORS)
	r.Static("/stat-img", "../image")
	r.Static("/static", "./webClient/dist")
	r.Static("/statics", "./webClient/dist/statics")
	r.StaticFile("/", "./webClient/dist/index.html")

	// АВТОРИЗАЦИЯ
	authRoute := r.Group("/auth")
	{
		// авторизация через email
		authRoute.POST("/email", auth.EmailAuth)
		authRoute.POST("/check_user_email", auth.EmailAuthCheckUserEmail)
		authRoute.POST("/email_auth_start_recover_password", auth.EmailAuthStartRecoverPassword)
		authRoute.POST("/email_auth_recover_password", auth.EmailAuthRecoverPassword)
	}

	apiRoute := r.Group("/api", authRequired)
	{
		apiRoute.POST("/current_user", apiCurrentUser)
		apiRoute.POST("/call_pg_func", apiCallPgFunc)
		// отправка логов в graylog
		apiRoute.POST("/log", logToGraylog)
		// подключение по SSE
		apiRoute.GET("/sse", sse.AddConn)
		// операции с файлами
		apiRoute.POST("/upload_file", uploadFile)
		apiRoute.GET("/file/:fileToken", downloadFile)
		apiRoute.POST("/remove_file/:fileToken", deleteFile)
		// загрузка фото
		apiRoute.POST("/upload_image", uploadImage)
		apiRoute.POST("/upload_profile_image", uploadProfileImage)
	}

	// на ненайденный url отправляем статический файл для запуска vuejs приложения
	r.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./webClient/dist/index.html")
	})

	err := r.Run(fmt.Sprintf(":%v", config.WebServer.Port))
	utils.CheckErr(err, "webserver run")
}

func apiCurrentUser(c *gin.Context) {
	if u, ok := c.Get(utils.GinContextUser); ok {
		if user, ok := u.(*types.User); ok {
			utils.HttpSuccess(c, user)
		} else {
			utils.HttpError(c, http.StatusNonAuthoritativeInfo, "user problem cast to model.User")
		}
	} else {
		utils.HttpError(c, http.StatusNonAuthoritativeInfo, "user not found in gin.context")
	}
}