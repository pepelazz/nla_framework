package main

import (
	"encoding/gob"
	"flag"
	"{{.Config.LocalProjectPath}}/graylog"
	"{{.Config.LocalProjectPath}}/jobs"
	"{{.Config.LocalProjectPath}}/pg"
	"{{.Config.LocalProjectPath}}/types"
	"{{.Config.LocalProjectPath}}/utils"
	"{{.Config.LocalProjectPath}}/webServer"
	"{{.Config.LocalProjectPath}}/sse"
	{{if .IsBitrixIntegration -}}
	"{{.Config.LocalProjectPath}}/bitrix"
	{{- end}}
	{{if .IsTelegramIntegration -}}
	"{{.Config.LocalProjectPath}}/tgBot"
	{{- end}}
	{{if .IsOdataIntegration -}}
	"{{.Config.LocalProjectPath}}/odata"
	{{- end}}
	"math/rand"
	"os"
	"time"
)

var (
	config *types.Config
	err    error
)

func main() {

	// считываем флаг dev. Если режим разработки, то меняем глобальные переменные
	isDev := flag.Bool("dev", false, "a bool")
	pgPort := flag.String("pg_port", "", "an string")
	pgPassword := flag.String("pg_pass", "", "an string")
	dbName := flag.String("dbname", "", "an string")
	{{if .IsTelegramIntegration -}}
	tgBotName := flag.String("telegram_bot_name", "", "an string")
	tgBotToken:= flag.String("telegram_bot_token", "", "an string")
	{{- end}}
	flag.Parse()

	if *isDev {
		_ = os.Setenv("PG_PORT", "5438")
		if len(*pgPort) > 0 {
			_ = os.Setenv("PG_PORT", *pgPort)
		}
		if len(*pgPassword) > 0 {
			_ = os.Setenv("PG_PASSWORD", *pgPassword)
		}
		_ = os.Setenv("PG_HOST", "localhost")
		if len(*dbName) > 0 {
			_ = os.Setenv("PG_DBNAME", *dbName)
		}
		{{if .IsTelegramIntegration -}}
		if len(*tgBotName) > 0 {
			_ = os.Setenv("TELEGRAM_BOT_NAME", *tgBotName)
		} else {
			utils.Panic("Write 'telegram_bot_name' and 'telegram_bot_token' in go parameters for developmeent mode")
		}
		if len(*tgBotToken) > 0 {
			_ = os.Setenv("TELEGRAM_BOT_TOKEN", *tgBotToken)
		}
		{{- end}}
		_ = os.Setenv("IS_DEVELOPMENT", "true")
	}

	// read config.toml
	config, err = types.ReadConfigFile("./config.toml")
	utils.CheckErr(err, "Read config")

	// postgres
	err = pg.StartPostgres(config.Postgres)
	utils.CheckErr(err, "StartPostgres")

	// подключаемся к серверу сбора логов
	err = graylog.Init(config.Graylog)
	utils.CheckErr(err, "Connect to GraylogConfig")

	// инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())
	//
	gob.Register(map[string]interface{}{})
	//
	jobs.StartJobs()

	// передаем часть конфига в utils
	utils.SetWebServerConfig(config.WebServer)
	utils.SetEmailConfig(config.Email)
	{{if .IsBitrixIntegration -}}
	bitrix.SetBitrixConfig(config.Bitrix)
	{{- end}}
	{{if .IsOdataIntegration -}}
	odata.SetOdataConfig(config.Odata)
	{{- end}}

	//go pg.GenerateFakeUsers(100)
	{{if .IsTelegramIntegration -}}
	go tgBot.Start(*config)
	{{- end}}

	// инициализируем брокера для обработки подключений по SSE
	sse.Init()

	webServer.StartWebServer(*config)
}
