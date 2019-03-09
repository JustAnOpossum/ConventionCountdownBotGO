package main

import (
	"os"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

func createMainMenu(isSubscribed bool) [][]tgAPI.InlineButton {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)
	if isSubscribed == true {
		finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("unsub", "Unsubscribe From Countdown", handleUnsub)})
	} else {
		finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("sub", "Subscribe To Countdown", handleSub)})
	}
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("cmd", "Commands", handleCommand)})
	return finalKeyboard
}

func createCmdKeybaord() [][]tgAPI.InlineButton {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)

	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("sub", "Subscribe To Countdown", handleSub)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("unsub", "Unsubscribe From Countdown", handleUnsub)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("info", "Bot Information", handleInfo)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("days", "Days Until "+os.Getenv("CON"), handleDays)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("home", "Back To Main Menu", handleHome)})

	return finalKeyboard
}

func createBackKeyboard() [][]tgAPI.InlineButton {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("home", "Back To Main Menu", handleHome)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("cmd", "Commands", handleCommand)})
	return finalKeyboard
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
