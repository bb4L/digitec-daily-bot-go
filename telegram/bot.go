// Package handling interaction with [Telegram] and defining the bot.
//
// [Telegram]: https://telegram.org
package telegram

import (
	"log"
	"os"
	"time"

	"github.com/bb4L/digitec-daily-bot-go/storage"
	"github.com/bb4L/digitec-daily-bot-go/task"
	tb "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

var logger = log.New(os.Stdout, "[bot] ", log.Ldate|log.Ltime|log.Lmsgprefix)
var storageHelper *storage.StorageHelper

// RunBot runs the bot itself
func RunBot(token string, offerChannel chan string, storage *storage.StorageHelper) {
	storageHelper = storage
	logger.Println("starting telegram bot")
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 2 * time.Second},
	})

	if err != nil {
		logger.Fatal(err)
	}

	b.Use(middleware.AutoRespond())

	addDefaultHandler(b)
	addCommandsHandler(b)
	go func() {
		for {
			s := <-offerChannel
			logger.Println("offer channel: ", s)
			for _, user := range storage.GetTaskSettings().Users {
				b.Send(&tb.User{ID: user}, s)
			}
		}
	}()
	b.Start()
}

func addDefaultHandler(b *tb.Bot) {
	logger.Println("add default handler")
	b.Handle(tb.OnText, func(c tb.Context) error {
		logger.Println("handle text")
		logger.Println(c.Text())
		return c.Send("unknown command: " + c.Text())
	})

}

func addCommandsHandler(b *tb.Bot) {
	logger.Println("add commands handler")

	b.Handle("/help", func(c tb.Context) error {
		c.Send(storageHelper.GetSettings().HelpText)
		return nil
	})
	b.Handle("/start", func(c tb.Context) error {
		storageHelper.AddUser(c.Sender().ID)
		logger.Printf("sending: \"%s\" to: \"%d\"", storageHelper.GetSettings().StartMessage, c.Sender().ID)
		c.Send(storageHelper.GetSettings().StartMessage)
		return nil
	})
	b.Handle("/stop", func(c tb.Context) error {
		storageHelper.RemoveUser(c.Sender().ID)
		c.Send(storageHelper.GetSettings().StopMessage)
		return nil
	})
	b.Handle("/current", func(c tb.Context) error {
		c.Send(task.GetMessageText(storageHelper))
		return nil
	})
}
