package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"time"

	tgAPI "gopkg.in/telebot.v3"
)

type configStruct struct {
	Con          string
	Date         time.Time
	DBName       string
	DBUrl        string
	Token        string
	WebhookURL   string
	WebhookPort  string
	MessageOnSub bool
	MainBot      mainBotStruct
	ImgSend      imgSendStruct
	Twitter      twitterStruct
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

type twitterStruct struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

var bot *tgAPI.Bot
var config configStruct
var users *datastore
var photos *datastore
var dataDir = os.Getenv("DATADIR")
var imgDir = dataDir + "/img"
var out = ioutil.Discard

func logError(err error) {
	if err != nil {
		fmt.Printf("+%v", err)
	}
}

func main() {
	if os.Getenv("DEBUG") == "true" {
		out = os.Stdout
	}

	if os.Getenv("MODE") == "" {
		askQuestions()
		loadConfig(dataDir, &config)
		fmt.Println("Please Wait... Connecting to Database")
		users, photos = setUpDB(config.DBName, config.DBUrl)
		fmt.Println("Connected to Database!")
		err := uploadZip()
		if err == nil {
			fmt.Println("Added All Photos!")
		} else {
			fmt.Println(err)
		}
		users.client.Disconnect(context.Background())
		return
	}

	loadConfig(dataDir, &config)
	users, photos = setUpDB(config.DBName, config.DBUrl)

	switch os.Getenv("MODE") {
	case "test":
		bot = setUpBot("test")
	case "prod":
		bot = setUpBot("prod")
	case "send":
		bot = setUpBot("send")
		if getDays(config.Date) > config.ImgSend.DayToStart || getDays(config.Date) < 0 {
			return
		}
		checkForAPI()
		returnedImg, err := createImg()
		if err != nil {
			logError(err)
			return
		}
		fmt.Fprintln(out, "Got Image")
		err = sendTelegramPhoto(returnedImg)
		logError(err)
		if config.Twitter.ConsumerKey != "" {
			mediaID, err := uploadTwitterMedia(returnedImg.FilePath, "image/jpeg")
			if err != nil {
				logError(err)
				return
			}

			twitterCaption := intToEmoji(returnedImg.DaysLeft) + " Days Until " + config.Con + "!\n\n📸: " + returnedImg.CreditName + " " + returnedImg.CreditURL
			sendMediaTweet(mediaID, twitterCaption)
		}

		users.client.Disconnect(context.Background())
		return
	}

	bot.Handle("/start", handleStart)
	bot.Handle("/menu", handleStart)
	bot.Handle(tgAPI.OnAddedToGroup, handleGroupAdd)
	bot.Handle(tgAPI.OnMigration, handleMigration)

	//Creates all available keyboards for later use
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
	timeUntil := time.Until(day)
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
