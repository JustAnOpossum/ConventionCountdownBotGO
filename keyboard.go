package main

import (
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

var keyboards = make(map[string][][]tgAPI.InlineButton)

func createMainMenu(isSubscribed bool) {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)
	var nameToAppend string
	if isSubscribed == true {
		finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("unsub", "Unsubscribe From Countdown", handleUnsubBtn)})
		nameToAppend = "Unsub"
	} else {
		finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("sub", "Subscribe To Countdown", handleSubBtn)})
		nameToAppend = "Sub"
	}
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("cmd", "Commands", handleCommandBtn)})
	keyboards["main"+nameToAppend] = finalKeyboard
}

func createCmdKeybaord() {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)

	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("sub", "Subscribe To Countdown", handleSubBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("unsub", "Unsubscribe From Countdown", handleUnsubBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("info", "Bot Information", handleInfoBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("days", "Days Until "+config.Con, handleDaysBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("home", "Back To Main Menu", handleHomeBtn)})

	keyboards["cmd"] = finalKeyboard
}

func createBackKeyboard() {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("home", "Back To Main Menu", handleHomeBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("cmd", "Commands", handleCommandBtn)})
	keyboards["back"] = finalKeyboard
}

func createBtn(uniqueName, text string, callback func(*tgAPI.Callback)) tgAPI.InlineButton {
	btn := tgAPI.InlineButton{
		Unique: uniqueName,
		Text:   text,
	}
	bot.Handle(&btn, callback)
	return btn
}

func handleBtnClick(msgTxt string, msgKeyboard [][]tgAPI.InlineButton, c *tgAPI.Callback) {
	bot.Edit(c.Message, msgTxt, &tgAPI.ReplyMarkup{
		InlineKeyboard: msgKeyboard,
	})
	bot.Respond(c, &tgAPI.CallbackResponse{})
}
