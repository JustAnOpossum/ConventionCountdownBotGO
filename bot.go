package main

import (
	"os"
	"time"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

func setUpBot(botMode string) *tgAPI.Bot {
	var tempBot *tgAPI.Bot
	var err error
	switch botMode {
	case "test":
		tempBot, err = tgAPI.NewBot(tgAPI.Settings{
			Token:  os.Getenv("TOKEN"),
			Poller: &tgAPI.LongPoller{Timeout: 10 * time.Second},
		})
		break
	case "prod":
		break
	case "send":
		tempBot, err = tgAPI.NewBot(tgAPI.Settings{
			Token: os.Getenv("TOKEN"),
		})
		break
	}
	if err != nil {
		panic(err)
	}
	return tempBot
}
