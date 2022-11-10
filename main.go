//main.go
//Contains the entry for the program, and some helper functions to accomplish getting the bot up and running

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"time"

	tgAPI "gopkg.in/telebot.v3"
)

type configStruct struct {
	Con          string
	Date         string
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
	DayToStart  int
	FontSize    float64
	Font        string
	AnimalEmoji string
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
var out = io.Discard

func logError(err error) {
	if err != nil {
		log.Println("Error: " + err.Error())
	}
}

func main() {
	//Silent debug output for uploading a zip file and other methods that use fmt
	if os.Getenv("DEBUG") == "true" {
		out = os.Stdout
	}

	//Creates directories if they don't exist
	if _, err := os.Stat(path.Join(dataDir, "countdown/")); os.IsNotExist(err) {
		os.Mkdir(path.Join(dataDir, "countdown/"), 0644)
	}
	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		os.Mkdir(imgDir, 0644)
	}

	loadConfig(dataDir, &config)
	users, photos = setUpDB(config.DBName, config.DBUrl)

	//Options for the bot to start
	switch os.Getenv("MODE") {
	//Starts the bot in long polling mode. Usefull if no webhooks
	case "longPoll":
		bot = setUpBot("longPoll")
	//Starts the bot in webhook mode
	case "webhook":
		bot = setUpBot("webhook")
	//Starts the bot in send image mode
	case "send":
		log.Println("Loading bot in send mode")
		bot = setUpBot("send")
		if getDays(config.Date) > config.ImgSend.DayToStart || getDays(config.Date) < 1 {
			return
		}
		checkForAPI()
		returnedImg := createImg()
		log.Println("Generated image for day " + strconv.Itoa(returnedImg.DaysLeft))
		err := sendTelegramPhoto(returnedImg)
		logError(err)
		if config.Twitter.ConsumerKey != "" {
			mediaID, err := uploadTwitterImg(returnedImg.FilePath)
			if err != nil {
				logError(err)
				return
			}

			twitterCaption := intToEmoji(returnedImg.DaysLeft) + " Days Until " + config.Con + "!\n\n📸: " + returnedImg.CreditName + " " + returnedImg.CreditURL
			sendMediaTweet(mediaID, twitterCaption)
		}

		users.client.Disconnect(context.Background())
		return
	//Starts the bot in upload mode
	case "upload":
		imgDir = dataDir + "/img"
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
	default:
		fmt.Println("Error: Please specify a mode though MODE=(longPoll, webhook, send, upload)")
		return
	}

	//Sets up handleres for telegram
	bot.Handle("/start", handleStart)
	bot.Handle("/menu", handleStart)
	// bot.Handle("/test", handleTest)
	bot.Handle(tgAPI.OnAddedToGroup, handleGroupAdd)
	bot.Handle(tgAPI.OnMigration, handleMigration)

	//Creates inline keyboards
	createMainMenu(true)
	createMainMenu(false)
	createBackKeyboard()
	createCmdKeybaord()

	log.Println("Telegram Bot is Started")
	bot.Start()
}

// Helper function to print errors
func handleErr(err error) {
	fmt.Printf("%+v", err)
}

// Gets how many days are left until the current con
// Counts number of full days before con. Set con date at midnight for best results
func getDays(day string) int {
	//Timezone aware parsed string
	parsedTime, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", day)
	if err != nil {
		panic(err)
	}
	timeUntil := time.Until(parsedTime)
	daysUntil := timeUntil.Hours() / 24
	daysRounded := math.Ceil(daysUntil)
	return int(daysRounded)
}

// Loads the config file
func loadConfig(dataDir string, configVar *configStruct) {
	configFile, err := os.ReadFile(path.Join(dataDir, "config.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configFile, configVar)
	if err != nil {
		panic(err)
	}
}
