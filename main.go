package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"time"

	tgAPI "gopkg.in/tucnak/telebot.v2"
)

type configStruct struct {
	Con     string
	Date    time.Time
	DBName  string
	Token   string
	MainBot mainBotStruct
	ImgSend imgSendStruct
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

type imgSendStruct struct {
	DayToStart int
	FontSize   float64
	Font       string
}

var bot *tgAPI.Bot
var config configStruct
var db *datastore
var dataDir = os.Getenv("DATADIR")
var imgDir = dataDir + "/img"

func main() {
	if os.Getenv("MODE") == "" {
		askQuestions()
		loadConfig(dataDir, &config)
		fmt.Println("Please Wait... Connecting to Database")
		db = setUpDB(config.DBName)
		fmt.Println("Connected to Database!")
		err := uploadZip()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Added All Photos!")
		db.session.Close()
		return
	}

	loadConfig(dataDir, &config)
	db = setUpDB(config.DBName)

	switch os.Getenv("MODE") {
	case "test":
		bot = setUpBot("test")
		break
	case "main":
		break
	case "send":
		if getDays(config.Date) > config.ImgSend.DayToStart || getDays(config.Date) < 0 {
			return
		}
		imgBlob, err := createImg()
		if err != nil {
			fmt.Printf("+%v", err)
		}
		checkForAPI()
		sendPhoto(imgBlob)
		db.session.Close()
		return
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

func getDays(day time.Time) int {
	timeUntil := day.Sub(time.Now())
	daysUntil := timeUntil.Hours() / 24
	daysRounded := math.Round(daysUntil)
	return int(daysRounded)
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
