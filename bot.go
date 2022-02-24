package main

import (
	"time"

	tgAPI "gopkg.in/telebot.v3"
)

func setUpBot(botMode string) *tgAPI.Bot {
	var tempBot *tgAPI.Bot
	var err error
	switch botMode {
	case "longPoll":
		tempBot, err = tgAPI.NewBot(tgAPI.Settings{
			Token:  config.Token,
			Poller: &tgAPI.LongPoller{Timeout: 10 * time.Second},
		})
	case "webhook":
		tempBot, err = tgAPI.NewBot(tgAPI.Settings{
			Token: config.Token,
			Poller: &tgAPI.Webhook{
				Listen: ":" + config.WebhookPort,
				Endpoint: &tgAPI.WebhookEndpoint{
					PublicURL: config.WebhookURL,
				},
			},
		})
	case "send":
		tempBot, err = tgAPI.NewBot(tgAPI.Settings{
			Token: config.Token,
		})
	}
	if err != nil {
		panic(err)
	}
	return tempBot
}
