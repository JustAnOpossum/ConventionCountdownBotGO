//keyboard.go
//Deals with the inline keyboard functions. Including creating and managing button clicks.

package main

import (
	tgAPI "gopkg.in/telebot.v3"
)

//Global to keep track of the created keyboards
var keyboards = make(map[string][][]tgAPI.InlineButton)

//Creates the main menu for the keyboard. This is the one referenced by /start and /menu
func createMainMenu(isSubscribed bool) {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)
	var nameToAppend string
	if isSubscribed {
		finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("unsub", "Unsubscribe From Countdown", handleUnsubBtn)})
		nameToAppend = "Unsub"
	} else {
		finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("sub", "Subscribe To Countdown", handleSubBtn)})
		nameToAppend = "Sub"
	}
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("cmd", "Commands", handleCommandBtn)})
	keyboards["main"+nameToAppend] = finalKeyboard
}

//Creates the command menu keyboard.
func createCmdKeybaord() {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)

	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("sub", "Subscribe To Countdown", handleSubBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("unsub", "Unsubscribe From Countdown", handleUnsubBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("info", "Bot Information", handleInfoBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("days", "Days Until "+config.Con, handleDaysBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("home", "Back To Main Menu", handleHomeBtn)})

	keyboards["cmd"] = finalKeyboard
}

//Creates the keyboard used after you run a command
func createBackKeyboard() {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("home", "Back To Main Menu", handleHomeBtn)})
	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("cmd", "Commands", handleCommandBtn)})
	keyboards["back"] = finalKeyboard
}

//Creates each button, registering the callback
func createBtn(uniqueName, text string, callback func(tgAPI.Context) error) tgAPI.InlineButton {
	btn := tgAPI.InlineButton{
		Unique: uniqueName,
		Text:   text,
	}
	bot.Handle(&btn, callback)
	return btn
}

//Sends the final response when a button is clicked
func handleBtnClick(msgTxt string, msgKeyboard [][]tgAPI.InlineButton, ctx tgAPI.Context) {
	ctx.Edit(msgTxt, &tgAPI.ReplyMarkup{
		InlineKeyboard: msgKeyboard,
	})
	ctx.Respond(&tgAPI.CallbackResponse{})
}
