package main

import (
	"github.com/globalsign/mgo/bson"
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

//Handlers From Telegram

func findSubOrUnsubKeyboard(chatID int64) [][]tgAPI.InlineButton {
	var keyboardToSend [][]tgAPI.InlineButton
	if db.itemExists("users", bson.M{"chatId": chatID}) == true {
		keyboardToSend = keyboards["mainUnsub"]
	} else {
		keyboardToSend = keyboards["mainSub"]
	}
	return keyboardToSend
}

func handleStart(msg *tgAPI.Message) {
	if msg.Chat.Type == "channel" || msg.Chat.Type == "privatechannel" {
		return
	}
	bot.Send(msg.Sender, config.WelcomeMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: findSubOrUnsubKeyboard(msg.Chat.ID),
	})
}

//Handalers For Keybaord
func handleSubBtn(c *tgAPI.Callback) {
	status := handleSub(c.Message.Chat)
	if status == true {
		handleBtnClick(config.SubMsg, keyboards["back"], c)
	} else {
		handleBtnClick(config.AlreadySubMsg, keyboards["back"], c)
	}
}

func handleUnsubBtn(c *tgAPI.Callback) {
	handleBtnClick(config.SubMsg, keyboards["back"], c)
}

func handleCommandBtn(c *tgAPI.Callback) {
	handleBtnClick(config.CmdMsg, keyboards["cmd"], c)
}

func handleHomeBtn(c *tgAPI.Callback) {
	handleBtnClick(config.WelcomeMsg, keyboards["mainUnsub"], c)
}

func handleInfoBtn(c *tgAPI.Callback) {
	handleBtnClick(config.InfoMsg, keyboards["back"], c)
}

func handleDaysBtn(c *tgAPI.Callback) {
	dayStr := getDays() + " Days Until " + config.Con + " !"
	handleBtnClick(dayStr, keyboards["back"], c)
}

func handleSub(chat *tgAPI.Chat) bool {
	if db.itemExists("users", bson.M{"chatId": chat.ID}) == true {
		return false
	}
	var isGroup bool
	if chat.Type != "ChatPrivate" {
		isGroup = true
	}
	itemToInsert := user{
		ChatID: chat.ID,
		Name:   chat.Username,
		Group:  isGroup,
	}
	db.insert("users", itemToInsert)
	return true
}

func handleUnsub() {

}
