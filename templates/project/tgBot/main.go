package tgBot

import (
	"[[.Config.LocalProjectPath]]/cacheUtil"
	"[[.Config.LocalProjectPath]]/pg"
	"[[.Config.LocalProjectPath]]/types"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
)

type tgUser struct {
	Id string
}

func (u *tgUser) Recipient() string {
	return u.Id
}

var (
	bot *tb.Bot
	// Universal markup builders.
	menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	selector = &tb.ReplyMarkup{}

	// Reply buttons.
	btnHelp     = menu.Text("ℹ Help")
	btnSettings = menu.Text("⚙ Settings")

	// Inline buttons.
	//
	// Pressing it will cause the client to
	// send the bot a callback.
	//
	// Make sure Unique stays unique as per button kind,
	// as it has to be for callback routing to work.
	//
	btnPrev = selector.Data("⬅", "prev")
	btnNext = selector.Data("➡", "next")
)

func Start(config types.Config) {
	var err error
	bot, err = tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		//URL: "http://195.129.111.17:8012",

		Token:  config.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: 1 * time.Second},
	})

	if err != nil {
		fmt.Printf("tgBot.Start tb.NewBot error: %s\n", err)
		return
	}

	// добавляем подписку на события в postgres
	pg.AddPgEventListener(pgListener)

	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)
	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)

	bot.Handle("/hello", func(m *tb.Message) {
		_, err := bot.Send(m.Sender, "Hello World!")
		if err != nil {
			fmt.Printf("tgBot send message error: %s\n", err)
		}
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		if strings.ToLower(m.Text) == "привет" {
			bot.Send(m.Sender, "Гамарджоба!")
			return
		}
		if strings.ToLower(m.Text) == "getid" || strings.ToLower(m.Text) == "get id" {
			bot.Send(m.Sender, fmt.Sprintf("Ваш id: %v", m.Sender.ID))
			return
		}
		if m.Text == "key" {
			bot.Send(m.Sender, "Hello!", menu)
			return
		}
		user, _ := userFindByTelegramId(m.Sender.ID)
		if user != nil {
			fmt.Printf("user: %s (%s) send '%s'\n", user.Fullname, m.Sender.Username, m.Text)
		} else {
			fmt.Printf("not auth user %s %v send '%s'\n", m.Sender.Username, m.Sender.ID, m.Text)
		}
	})

	bot.Handle(&btnHelp, func(m *tb.Message) {
		fmt.Printf("in btnHelp %s\n", m.Sender.Username)
	})

	bot.Handle(&btnSettings, func(m *tb.Message) {
		fmt.Printf("in btnSettings %s\n", m.Sender.Username)
	})

	// On inline button pressed (callback)
	//b.Handle(&btnPrev, func(c *tb.Callback) {
	//	// ...
	//	// Always respond!
	//	b.Respond(c, &tb.CallbackResponse{...})
	//})

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		err := bot.Download(&m.Photo.File, "test_photo.jpg")
		if err != nil {
			fmt.Printf("err %s\n", err)
		} else {
			bot.Send(m.Sender, "фото успешно сохранено")
		}
	})

	bot.Start()
}

func SendMsg(tgId, msg string) {
	if bot != nil && len(tgId) > 0 && len(msg) > 0 {
		msg = strings.Replace(msg, "\\n", "\n", -1)
		answer, err := bot.Send(&tgUser{tgId}, msg, tb.ModeHTML)
		if err != nil {
			fmt.Printf("bot.Send error: %s tgId:%s msg:'%s'\n", err, tgId, msg)
		}
		if answer != nil {
			fmt.Printf("bot.Send: tgId:%s msg:'%s' answer: %s\n", tgId, msg, answer.Text)
		}
	}
}

func SendSticker(tgId, fileId string) {
	if bot != nil && len(tgId) > 0 && len(fileId) > 0 {
		sticker := &tb.Sticker{
			File: tb.File{FileID: fileId},
		}
		_, err := bot.Send(&tgUser{tgId}, sticker, tb.ModeHTML)
		if err != nil {
			fmt.Printf("bot.Send error: %s tgId:%s msg:'%s'\n", err, tgId, fileId)
		}
	}
}

func pgListener(event string) {
	tableName := gjson.Get(event, "table").Str
	if tableName == "send_msg_to_user_telegram" {
		SendMsg(gjson.Get(event, "telegram_id").String(), gjson.Get(event, "msg").String())
	}
	if tableName == "send_sticker_to_user_telegram" {
		SendSticker(gjson.Get(event, "telegram_id").String(), gjson.Get(event, "file_id").String())
	}
}

func userFindByTelegramId(tgId interface{}) (user *types.User, err error) {
	cacheKey := fmt.Sprintf("telegram_id%v", tgId)
	// ищем пользователя в кэше
	userIntreface, _ := cacheUtil.GoCacheGet(cacheKey)
	if userIntreface != nil {
		var ok bool
		user, ok = userIntreface.(*types.User)
		if ok {
			return
		}
	}
	jsonStr, _ := json.Marshal(map[string]interface{}{"id": tgId})
	err = pg.CallPgFunc("user_get_by_telegram_id", jsonStr, &user, nil)
	// записываем в кэш
	if err == nil {
		cacheUtil.GoCacheSet(cacheKey, user, time.Minute*1)
	}
	return
}
