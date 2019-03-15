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

type configStruct struct {
	Con     string
	Date    time.Time
	DBName  string
	Token   string
	MainBot mainBotStruct
}

type mainBotStruct struct {
	WelcomeMsg       string
	SubMsg           string
	AlreadySubMsg    string
	GroupAddMsg      string
	GroupNotAdminMsg string
	NotSubMsg        string
	UnsubMsg         string
	CmdMsg           string
	InfoMsg          string
	Owners           string
}

var bot *tgAPI.Bot
var config configStruct
var db *datastore
var dataDir = os.Getenv("DATADIR")

func main() {
	if os.Getenv("MODE") == "" {
		askQuestions()
		loadConfig(dataDir, &config)
		err := uploadZip()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	loadConfig(dataDir, &config)

	switch os.Getenv("MODE") {
	case "test":
		db = setUpDB(config.DBName)

		bot = setUpBot("test")
		break
	case "main":
		break
	case "send":
		break
	}

	bot.Handle("/start", handleStart)
	bot.Handle("/menu", handleStart)
	bot.Handle(tgAPI.OnAddedToGroup, handleGroupAdd)
	bot.Handle(tgAPI.OnMigration, handleMigration)

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

func getDays(day time.Time) string {
	timeUntil := day.Sub(time.Now())
	daysUntil := timeUntil.Hours() / 24
	daysRounded := math.Round(daysUntil)
	return strconv.Itoa(int(daysRounded))
}

func loadConfig(dataDir string, configVar *configStruct) {
	configFile, err := ioutil.ReadFile(path.Join(dataDir, "config.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configFile, configVar)
	if err != nil {
		panic(err)
	}
}
