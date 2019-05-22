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
	Twitter twitterStruct
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
	DayToStart   int
	FontSize     float64
	Font         string
	AnimalEmoji  string
	Music        string
	VideoCaption string
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
var countdownDir = dataDir + "/countdown"
var out = ioutil.Discard

func main() {
	if os.Getenv("DEBUG") == "true" {
		out = os.Stdout
	}

	if os.Getenv("MODE") == "" {
		askQuestions()
		loadConfig(dataDir, &config)
		fmt.Println("Please Wait... Connecting to Database")
		users, photos = setUpDB(config.DBName)
		fmt.Println("Connected to Database!")
		err := uploadZip()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Added All Photos!")
		users.session.Close()
		return
	}

	loadConfig(dataDir, &config)
	users, photos = setUpDB(config.DBName)

	switch os.Getenv("MODE") {
	case "test":
		bot = setUpBot("test")
		break
	case "main":
		break
	case "send":
		var err error
		bot = setUpBot("send")
		if getDays(config.Date) > config.ImgSend.DayToStart || getDays(config.Date) < 0 {
			return
		}
		checkForAPI()
		if getDays(config.Date) == 0 {
			if slideshow, err := createSlideShow(); err == nil {
				err = sendVideo(slideshow)
			}
		} else {
			returnedImg, err := createImg()
			fmt.Fprintln(out, "Got Image")
			if err != nil {
				fmt.Printf("+%v", err)
			}
			err = sendTelegramPhoto(returnedImg)
			if config.Twitter.ConsumerKey != "" {
				err = sendTwitterPhoto(returnedImg)
			}
		}
		if err != nil {
			fmt.Printf("+%v", err)
		}
		users.session.Close()
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

	createSlideShow()

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
