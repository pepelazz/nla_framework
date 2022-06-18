package types

import (
	"fmt"
	"github.com/pepelazz/go-toml"
	"os"
	"strconv"
)

type Config struct {
	Postgres Postgres

	WebServer WebServer

	Graylog GraylogConfig

	Email EmailConfig
[[if.IsBitrixIntegration -]]
	Bitrix BitrixConfig
[[- end]]
[[if.IsTelegramIntegration -]]
	Telegram TelegramConfig
[[- end]]
[[if.IsOdataIntegration -]]
	Odata OdataConfig
[[- end]]
}

func ReadConfigFile(path string) (c *Config, err error) {

	tree, err := toml.LoadFile(path)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Printf("current directory (pwd): %s\n", pwd)
		return nil, err
	}

	c = &Config{}

	if tree.Has("postgres") {
		c.Postgres.User = tree.Get("postgres.user").(string)
		c.Postgres.Password = tree.Get("postgres.password").(string)
		if len(os.Getenv("PG_PASSWORD")) > 0 {
			// перезаписываем пароль, если есть глобальная переменная
			c.Postgres.Password = os.Getenv("PG_PASSWORD")
		}
		c.Postgres.DbName = tree.Get("postgres.dbName").(string)
		if len(os.Getenv("PG_DBNAME")) > 0 {
			// перезаписываем пароль, если есть глобальная переменная
			c.Postgres.DbName = os.Getenv("PG_DBNAME")
		}
		c.Postgres.Host = tree.Get("postgres.host").(string)
		if len(os.Getenv("PG_HOST")) > 0 {
			// перезаписываем имя хоста, если есть глобальная переменная (для docker-compose)
			c.Postgres.Host = os.Getenv("PG_HOST")
		}
		c.Postgres.Port = tree.Get("postgres.port").(int64)
		if len(os.Getenv("PG_PORT")) > 0 {
			// перезаписываем порт, если есть глобальная переменная (для docker-compose)
			var port int64
			port, err = strconv.ParseInt(os.Getenv("PG_PORT"), 10, 64)
			if err != nil {
				return
			}
			c.Postgres.Port = port
		}
	}

	if tree.Has("webServer") {
		if tree.Has("webServer.enable") {
			c.WebServer.Enable = true
		}
		if tree.Has("webServer.port") {
			c.WebServer.Port = tree.Get("webServer.port").(int64)
		} else {
			c.WebServer.Port = 8085
		}
		if tree.Has("webServer.url") {
			c.WebServer.Url = tree.Get("webServer.url").(string)
			if os.Getenv("IS_DEVELOPMENT") == "true" {
				c.WebServer.Url = "http://localhost:8080"
			}
		} else {
			c.WebServer.Url = "localhost"
		}
	}

	if tree.Has("graylog") {
		if tree.Has("graylog.host") {
			c.Graylog.Host = tree.Get("graylog.host").(string)
		}
		if tree.Has("graylog.port") {
			c.Graylog.Port = int(tree.Get("graylog.port").(int64))
		}
		if tree.Has("graylog.appName") {
			c.Graylog.AppName = tree.Get("graylog.appName").(string)
		}
	}

	if tree.Has("email") {
		c.Email.Sender = tree.Get("email.sender").(string)
		if len(os.Getenv("EMAIL_SENDER")) > 0 {
			c.Email.Sender = os.Getenv("EMAIL_SENDER")
		}
		c.Email.Password = tree.Get("email.password").(string)
		if len(os.Getenv("EMAIL_PASSWORD")) > 0 {
			c.Email.Password = os.Getenv("EMAIL_PASSWORD")
		}
		c.Email.Host = tree.Get("email.host").(string)
		if len(os.Getenv("EMAIL_HOST")) > 0 {
			c.Email.Host = os.Getenv("EMAIL_HOST")
		}
		if tree.Has("email.port") {
			c.Email.Port = tree.Get("email.port").(int64)
		} else {
			c.Email.Port = 25
		}
		if len(os.Getenv("EMAIL_PORT")) > 0 {
			c.Email.Port, err = strconv.ParseInt(os.Getenv("EMAIL_PORT"), 10, 64)
			if err != nil {
				return nil, err
			}
		}
		if tree.Has("email.senderName") {
			c.Email.SenderName = tree.Get("email.senderName").(string)
		}
		if tree.Has("email.senderLogo") {
			c.Email.SenderLogo = tree.Get("email.senderLogo").(string)
		}
		if tree.Has("email.isSendWithEmptySender") {
			c.Email.IsSendWithEmptySender = tree.Get("email.isSendWithEmptySender").(bool)
		}
	}
[[if.IsBitrixIntegration -]]
if tree.Has("bitrix") {
if tree.Has("bitrix.apiUrl") {
c.Bitrix.ApiUrl = tree.Get("bitrix.apiUrl").(string)
}
if tree.Has("bitrix.userId") {
c.Bitrix.UserId = tree.Get("bitrix.userId").(string)
}
if tree.Has("bitrix.webhookToken") {
c.Bitrix.WebhookToken = tree.Get("bitrix.webhookToken").(string)
}
}
[[- end]]

[[if.IsTelegramIntegration -]]
if tree.Has("telegram") {
if tree.Has("telegram.botName") {
c.Telegram.BotName = tree.Get("telegram.botName").(string)
if len(os.Getenv("TG_BOT_NAME")) > 0 {
// перезаписываем, если есть глобальная переменная
c.Telegram.BotName = os.Getenv("TG_BOT_NAME")
}
}
if tree.Has("telegram.token") {
c.Telegram.Token = tree.Get("telegram.token").(string)
if len(os.Getenv("TELEGRAM_BOT_TOKEN")) > 0 {
// перезаписываем, если есть глобальная переменная
c.Telegram.Token = os.Getenv("TELEGRAM_BOT_TOKEN")
}
}
}
[[- end]]

[[if.IsOdataIntegration -]]
if tree.Has("odata") {
if tree.Has("odata.url") {
c.Odata.Url = tree.Get("odata.url").(string)
}
if tree.Has("odata.login") {
c.Odata.Login = tree.Get("odata.login").(string)
}
if tree.Has("odata.password") {
c.Odata.Password = tree.Get("odata.password").(string)
}
if tree.Has("odata.exchangePlanName") {
c.Odata.ExchangePlanName = tree.Get("odata.exchangePlanName").(string)
}
if tree.Has("odata.exchangePlanGuid") {
c.Odata.ExchangePlanGuid = tree.Get("odata.exchangePlanGuid").(string)
}
}
[[- end]]

return
}
