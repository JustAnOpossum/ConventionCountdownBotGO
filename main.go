package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

var bot *tgAPI.Bot
var config configStruct
var dataDir = os.Getenv("DATADIR")

type configStruct struct {
	WelcomeMsg string
	SubMsg     string
	UnsubMsg   string
	CmdMsg     string
}

func main() {
	bot = setUpBot("test")

	bot.Handle("/start", handleStart)

	configFile, err := ioutil.ReadFile(path.Join(dataDir, "config.json"))
	if err != nil {
		panic(err)
	}
	json.Unmarshal(configFile, &config)

	fmt.Println("Telegram Bot is Started")
	bot.Start()
}
