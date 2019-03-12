package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
	"time"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

var bot *tgAPI.Bot
var config configStruct
var db *datastore
var dataDir = os.Getenv("DATADIR")

type configStruct struct {
	Con           string
	WelcomeMsg    string
	SubMsg        string
	AlreadySubMsg string
	UnsubMsg      string
	CmdMsg        string
	InfoMsg       string
	Date          time.Time
	DBName        string
}

func main() {
	bot = setUpBot("test")

	bot.Handle("/start", handleStart)

	configFile, err := ioutil.ReadFile(path.Join(dataDir, "config.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		panic(nil)
	}
	db = setUpDB(config.DBName)

	createMainMenu(true)
	createMainMenu(false)
	createBackKeyboard()
	createCmdKeybaord()

	fmt.Println("Telegram Bot is Started")
	bot.Start()
}

func getDays() string {
	timeUntil := config.Date.Sub(time.Now())
	daysUntil := timeUntil.Hours() / 24
	daysRounded := math.Round(daysUntil)
	return strconv.Itoa(int(daysRounded))
}
