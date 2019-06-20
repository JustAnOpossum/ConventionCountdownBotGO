package main

import (
	"time"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

func setUpBot(botMode string) *tgAPI.Bot {
	var tempBot *tgAPI.Bot
	var err error
	switch botMode {
	case "test":
		tempBot, err = tgAPI.NewBot(tgAPI.Settings{
			Token:  config.Token,
			Poller: &tgAPI.LongPoller{Timeout: 10 * time.Second},
		})
		break
	case "prod":
		tempBot, err = tgAPI.NewBot(tgAPI.Settings{
			Token: config.Token,
			Poller: &tgAPI.Webhook{
				Listen: ":" + config.WebhookPort,
				Endpoint: &tgAPI.WebhookEndpoint{
					PublicURL: config.WebhookURL,
				},
			},
		})
		break
	case "send":
		tempBot, err = tgAPI.NewBot(tgAPI.Settings{
			Token: config.Token,
		})
		break
	}
	if err != nil {
		panic(err)
	}
	return tempBot
}
