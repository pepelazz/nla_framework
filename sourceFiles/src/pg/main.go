package pg

import (
	"database/sql"
	"fmt"
	"github.com/pepelazz/projectGenerator/types"
	"github.com/pepelazz/pg_generate"
)

var Pg *sql.DB

func StartPostgres(config types.Postgres) error {
	var err error
	// создаем базу
	pgGenerate.Start(false)
	// создаем подключение к базе
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.DbName)
	Pg, err = sql.Open("postgres", dbinfo)
	err = Pg.Ping()
	if err != nil {
		return err
	}
	// подписываемся на канал обновлений
	go pgListen(config)
	return nil
}
