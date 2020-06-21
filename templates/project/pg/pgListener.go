package pg

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"[[.Config.LocalProjectPath]]/cacheUtil"
	"[[.Config.LocalProjectPath]]/types"
	"[[.Config.LocalProjectPath]]/utils"
	"[[.Config.LocalProjectPath]]/sse"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
	[[if .IsTelegramIntegration -]]
	"[[.Config.LocalProjectPath]]/tgBot"
	[[- end]]
)

type (
	PgEventListener func(event string)
)

var (
	pgListeners = []PgEventListener{}
)

func waitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			processPgEvent(n.Extra)
			for _, f := range pgListeners {
				f(n.Extra)
			}
			//printEventJson(n)
			return
		case <-time.After(90 * time.Second):
			//fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}

func pgListen(config types.Postgres) {

	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.DbName)
	db, err := sql.Open("postgres", dbinfo)
	err = db.Ping()
	utils.CheckErr(err, "Can't connect to postgres. Maybe wrong port.")
	defer db.Close()

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(dbinfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	for {
		waitForNotification(listener)
	}
}

func AddPgEventListener(f PgEventListener) {
	pgListeners = append(pgListeners, f)
}

func processPgEvent(event string) {
	fmt.Printf("event %s\n", event)
	// извлекаем тип документа для которого произошли изменения в базе
	tableName := gjson.Get(event, "table").Str
	//обрабатываем изменения
	switch tableName {
	case "user":
		// стираем пользователя из кэша
		token := gjson.Get(event, "auth_token").Str
		if len(token) > 0 {
			cacheUtil.UserRemoveByToken(token)
		}
	case "message":
		if (gjson.Get(event, "flds.tg_op").Str == "INSERT") {
			userIdInt := gjson.Get(event, "flds.user_id").Int()
			sse.SendJson(strconv.FormatInt(userIdInt, 10), gjson.Get(event, "flds").Value())
			[[if .IsTelegramIntegration -]]
			tgBot.SendMsg(gjson.Get(event, "flds.user_options.telegram_id").String(), gjson.Get(event, "flds.title").String())
			[[- end]]
		}
	case "task":
		userIdInt := gjson.Get(event, "flds.executor_id").Int()
		sse.SendJson(strconv.FormatInt(userIdInt, 10), gjson.Get(event, "flds").Value())
	case "process_error":
		fmt.Printf("postgres event %s\n", event)
	}
}
