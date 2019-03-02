package main

import (
	"fmt"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

var bot *tgAPI.Bot

func main() {
	bot = setUpBot("test")

	bot.Handle("/start", handleStart)

	createMainMenu()

	fmt.Println("Telegram Bot is Started")
	bot.Start()
}
