package main

import (
	"conBot/helper"
	"fmt"
	"os"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

var bot *tgAPI.Bot
var config helper.ConfigStruct
var db *helper.Datastore
var dataDir = os.Getenv("DATADIR")

func main() {
	bot = setUpBot("test")

	bot.Handle("/start", handleStart)
	bot.Handle(tgAPI.OnAddedToGroup, handleGroupAdd)

	helper.LoadConfig(dataDir, &config)
	db = helper.SetUpDB(config.DBName)

	createMainMenu(true)
	createMainMenu(false)
	createBackKeyboard()
	createCmdKeybaord()

	fmt.Println("Telegram Bot is Started")
	bot.Start()
}

func handleErr(err error) {
	fmt.Printf("%+v", err)
}
