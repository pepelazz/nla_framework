package webServer

import (
	"[[.Config.LocalProjectPath]]/sse"
	"[[.Config.LocalProjectPath]]/types"
	"[[.Config.LocalProjectPath]]/utils"
	"[[.Config.LocalProjectPath]]/webServer/auth"
	"github.com/gin-gonic/gin"

[[if .IsBitrixIntegration -]]
	"[[.Config.LocalProjectPath]]/bitrix"
	[[- end]]
	[[if .IsOdataIntegration -]]
	"[[.Config.LocalProjectPath]]/odata"
	[[- end]]
	"net/http"
	[[- range .Go.Routes.Imports]]
		"[[.]]"
	[[- end]]
	"fmt"
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
		[[if .Config.Auth.ByPhone -]]
		// авторизация по номеру телефона
		authRoute.POST("/phone", auth.PhoneAuth)
		authRoute.POST("/check_sms_code", auth.CheckSmsCode)
		authRoute.POST("/phone_auth_start_recover_password", auth.PhoneAuthStartRecoverPassword)
		authRoute.POST("/phone_auth_recover_password", auth.PhoneAuthRecoverPassword)
		[[- end]]
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
		[[if .IsTelegramIntegration -]]
		apiRoute.POST("/telegram_auth", telegramAuth(config.Telegram))
		[[- end]]

		[[if .IsBitrixIntegration -]]
		// импорт данных из Битрикс
		btxRoute := apiRoute.Group("/bitrix")
		[[- range .Docs -]]
			[[- if .IsBitrixIntegration]]
		btxRoute.POST("/import_[[.Name]]", bitrix.Get[[ToCamel .Name]]History)
			[[- end ]]
		[[- end ]]
		[[- end ]]

		[[if .IsOdataIntegration -]]
		// импорт данных из 1С Odata
		odataRoute := apiRoute.Group("/odata")
		[[- range .Docs -]]
			[[- if .IsOdataIntegration]]
			odataRoute.POST("/import_[[.Name]]", odata.Start[[ToCamel .Name]]Sync)
			[[- end ]]
		[[- end ]]
		[[- end ]]
		[[- range .Go.Routes.Api]]
			[[.]]
		[[- end]]
	}

	[[- range .Go.Routes.NotAuth]]
	[[.]]
	[[- end]]

	[[if .IsBitrixIntegration -]]
	// отладочные методы для импорта данных из Битрикс
	[[- range .Docs -]]
	[[- if .IsBitrixIntegrationDebugMode]]
	r.GET("/bitrix/import_[[.Name]]", bitrix.Get[[ToCamel .Name]]HistoryDebug)
	[[- end ]]
	[[- end ]]
	[[- end ]]

	[[if .IsOdataIntegration -]]
	// отладочные методы для импорта данных из Битрикс
	[[- range .Docs -]]
	[[- if .IsOdataIntegrationDebugMode]]
	r.GET("/odata/import_[[.Name]]",  odata.Sync[[ToCamel .Name]]With1CDebug)
	[[- end ]]
	[[- end ]]
	[[- end ]]


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