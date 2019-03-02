package main

import (
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

var keyboards = make(map[string][][]tgAPI.InlineButton, 0)

func createMainMenu() {
	var finalKeyboard = make([][]tgAPI.InlineButton, 0)

	finalKeyboard = append(finalKeyboard, []tgAPI.InlineButton{createBtn("Test1", "Test")})

	keyboards["mainMenu"] = finalKeyboard
}

func createBtn(uniqueName, text string) tgAPI.InlineButton {
	btn := tgAPI.InlineButton{
		Unique: uniqueName,
		Text:   text,
	}
	bot.Handle(&btn)
	return btn
}
