package types

import (
	"github.com/pelletier/go-toml"
	"os"
	"strconv"
)

type Config struct {
	Postgres Postgres

	WebServer WebServer

	Graylog GraylogConfig

	Email EmailConfig
}

func ReadConfigFile(path string) (c *Config, err error) {

	tree, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	c = &Config{}

	if tree.Has("postgres") {
		c.Postgres.User = tree.Get("postgres.user").(string)
		c.Postgres.Password = tree.Get("postgres.password").(string)
		c.Postgres.DbName = tree.Get("postgres.dbName").(string)
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
	}

	if tree.Has("email") {
		c.Email.Sender = tree.Get("email.sender").(string)
		c.Email.Password = tree.Get("email.password").(string)
		c.Email.Host = tree.Get("email.host").(string)
		if tree.Has("email.port") {
			c.Email.Port = tree.Get("email.port").(int64)
		} else {
			c.Email.Port = 25
		}
		if tree.Has("email.senderName") {
			c.Email.SenderName = tree.Get("email.senderName").(string)
		}
		if tree.Has("email.senderLogo") {
			c.Email.SenderLogo = tree.Get("email.senderLogo").(string)
		}
	}

	return
}
