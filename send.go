//send.go
//Sends both the telegram image and twitter image when the bot is running in send mode

package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"go.mongodb.org/mongo-driver/bson"

	tgAPI "gopkg.in/telebot.v3"
)

func sendTelegramPhoto(img finalImg) error {
	log.Println("Sending telegram messages")
	var photoCaption string
	if img.DaysLeft == 1 {
		photoCaption = "Tomorrow is " + config.Con + "!"
	} else {
		photoCaption = intToEmoji(img.DaysLeft) + " Days Until " + config.Con + "! " + findRandomAnimalEmoji() + "\n\n📸: [" + img.CreditName + "](" + img.CreditURL + ")"
	}
	sendPhoto := tgAPI.Photo{
		File:    tgAPI.FromDisk(img.FilePath),
		Caption: photoCaption,
	}
	var allUsers []user
	users.findAll(bson.M{}, &allUsers)

	for i := range allUsers {
		sendPhoto.Caption = photoCaption
		user := allUsers[i]
		tgUser := tgAPI.User{
			ID: user.ChatID,
		}
		_, err := sendPhoto.Send(bot, &tgUser, &tgAPI.SendOptions{
			ParseMode: tgAPI.ModeMarkdown,
		})
		if err != nil {
			//List of common errors to prevent users from being accidentaly removed by the bot
			if strings.Contains(err.Error(), "bot was blocked by the user") || strings.Contains(err.Error(), "user is deactivated") || strings.Contains(err.Error(), "chat not found") || strings.Contains(err.Error(), "have no rights to send a message") {
				users.removeOne(bson.M{"chatId": user.ChatID})
			}
			log.Println(err)
		}
	}

	return nil
}

func sendMediaTweet(mediaID int64, tweetText string) error {
	log.Println("Sending twitter image")
	os.Setenv("GOTWI_API_KEY", config.Twitter.ConsumerKey)
	os.Setenv("GOTWI_API_KEY_SECRET", config.Twitter.ConsumerSecret)
	token := gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           config.Twitter.AccessToken,
		OAuthTokenSecret:     config.Twitter.AccessSecret,
	}
	client, err := gotwi.NewClient(&token)
	if err != nil {
		return err
	}

	req := &types.CreateInput{
		Text: gotwi.String(tweetText),
		Media: &types.CreateInputMedia{
			MediaIDs: []string{strconv.Itoa(int(mediaID))},
		},
	}
	_, err = managetweet.Create(context.Background(), client, req)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func intToEmoji(input int) string {
	parsedInt := strconv.Itoa(input)
	splitString := strings.Split(parsedInt, "")
	var finalString string

	for i := range splitString {
		switch splitString[i] {
		case "0":
			finalString += string("\u0030\uFE0F\u20E3")
		case "1":
			finalString += string("\u0031\uFE0F\u20E3")
		case "2":
			finalString += string("\u0032\uFE0F\u20E3")
		case "3":
			finalString += string("\u0033\uFE0F\u20E3")
		case "4":
			finalString += string("\u0034\uFE0F\u20E3")
		case "5":
			finalString += string("\u0035\uFE0F\u20E3")
		case "6":
			finalString += string("\u0036\uFE0F\u20E3")
		case "7":
			finalString += string("\u0037\uFE0F\u20E3")
		case "8":
			finalString += string("\u0038\uFE0F\u20E3")
		case "9":
			finalString += string("\u0039\uFE0F\u20E3")
		}
	}
	return finalString
}

// Finds a random animal emoji from the config file
func findRandomAnimalEmoji() string {
	animals := strings.Split(config.ImgSend.AnimalEmoji, ",")
	randSrc := rand.NewSource(time.Now().Unix())
	random := rand.New(randSrc)
	return animals[random.Intn(len(animals))]
}

func checkForAPI() {
	for {
		resp, err := http.Get("https://api.telegram.org")
		if err != nil {
			time.Sleep(time.Minute * 2)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			time.Sleep(time.Minute * 2)
			continue
		}
		return
	}
}
