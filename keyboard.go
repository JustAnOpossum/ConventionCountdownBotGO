package main

import (
	"fmt"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

func createMainMenu(isSubscribed bool) [][]tgAPI.InlineButton {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)
	if isSubscribed == true {
		finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("unsub", "Unsubscribe To Countdown", handleUnsub)})
	} else {
		finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("sub", "Subscribe To Countdown", handleSub)})
	}
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

func handleUnsub(c *tgAPI.Callback) {
}

func handleSub(c *tgAPI.Callback) {
	fmt.Println(c.MessageID)
	bot.Edit(c.Message, c.Message.Text, "Testing2")
	bot.Edit(c.Message, tgAPI.ReplyMarkup{}, createMainMenu(true))
	bot.Respond(c, &tgAPI.CallbackResponse{})
}

func handleCommand(c *tgAPI.Callback) {

}
