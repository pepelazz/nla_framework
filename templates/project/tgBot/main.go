package tgBot

import (
	"fmt"
	"[[.Config.LocalProjectPath]]/types"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type tgUser struct {
	Id string
}

func (u *tgUser) Recipient() string  {
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
		if m.Text == "key" {
			bot.Send(m.Sender, "Hello!", menu)
			return
		}
		fmt.Printf("%s send '%s'\n", m.Sender.Username, m.Text)
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

func SendMsg(tgId, msg string)  {
	if bot != nil && len(tgId) > 0 && len(msg) > 0 {
		bot.Send(&tgUser{tgId}, msg)
	}
}