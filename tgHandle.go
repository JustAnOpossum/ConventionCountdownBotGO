package main

import (
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

func createMainMenu(msg tgAPI.Message) {
	bot.Send(msg.Sender, "Hello World")
}
