package main

import (
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

func handleStart(msg *tgAPI.Message) {
	bot.Send(msg.Sender, "Testing", &tgAPI.ReplyMarkup{
		InlineKeyboard: createMainMenu(false),
	})
}
