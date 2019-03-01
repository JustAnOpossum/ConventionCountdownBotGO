package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

var bot *tgAPI.Bot

func main() {
	bot = setUpBot("test")

	bot.Handle("/start", createMainMenu)

	bot.Start()

	fmt.Println("Telegram Bot is Started")
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-exitChan

	bot.Stop()
}
