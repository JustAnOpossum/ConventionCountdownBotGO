package main

import (
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

//Handlers From Telegram
func handleStart(msg *tgAPI.Message) {
	bot.Send(msg.Sender, config.WelcomeMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: keyboards["mainSub"],
	})
}

//Handalers For Keybaord
func handleSub(c *tgAPI.Callback) {
	handleBtnClick(config.SubMsg, keyboards["back"], c)
}

func handleUnsub(c *tgAPI.Callback) {
	handleBtnClick(config.SubMsg, keyboards["back"], c)
}

func handleCommand(c *tgAPI.Callback) {
	handleBtnClick(config.CmdMsg, keyboards["cmd"], c)
}

func handleHome(c *tgAPI.Callback) {
	handleBtnClick(config.WelcomeMsg, keyboards["mainUnsub"], c)
}

func handleInfo(c *tgAPI.Callback) {
	handleBtnClick(config.InfoMsg, keyboards["back"], c)
}

func handleDays(c *tgAPI.Callback) {
	dayStr := getDays() + " Days Until " + config.Con + " !"
	handleBtnClick(dayStr, keyboards["back"], c)
}
