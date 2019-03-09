package main

import (
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

//Handlers From Telegram
func handleStart(msg *tgAPI.Message) {
	bot.Send(msg.Sender, config.WelcomeMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: createMainMenu(false),
	})
}

//Handalers For Keybaord
func handleSub(c *tgAPI.Callback) {
	handleBtnClick(config.SubMsg, createBackKeyboard(), c)
}

func handleUnsub(c *tgAPI.Callback) {

}

func handleCommand(c *tgAPI.Callback) {
	handleBtnClick(config.CmdMsg, createCmdKeybaord(), c)
}

func handleHome(c *tgAPI.Callback) {
	handleBtnClick(config.WelcomeMsg, createMainMenu(true), c)
}

func handleInfo(c *tgAPI.Callback) {

}

func handleDays(c *tgAPI.Callback) {

}
