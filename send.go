package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/globalsign/mgo/bson"

	tgAPI "gopkg.in/telebot.v3"
)

func sendTelegramPhoto(img finalImg) error {
	photoCaption := intToEmoji(img.DaysLeft) + " Days Until " + config.Con + "!\n\n📸: [" + img.CreditName + "](" + img.CreditURL + ")"
	sendPhoto := tgAPI.Photo{
		File:    tgAPI.FromReader(img.ImgReader),
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
			if err.Error() == "api error: Forbidden: bot was blocked by the user" || err.Error() == "api error: Forbidden: user is deactivated" || err.Error() == "api error: Bad Request: chat not found" || err.Error() == "api error: Bad Request: have no rights to send a message" {
				users.removeOne(bson.M{"chatId": user.ChatID})
			}
		}
	}

	return nil
}

func sendMediaTweet(mediaID int64, tweetText string) error {
	twitterConfig := oauth1.NewConfig(config.Twitter.ConsumerKey, config.Twitter.ConsumerSecret)
	twitterToken := oauth1.NewToken(config.Twitter.AccessToken, config.Twitter.AccessSecret)
	httpClient := twitterConfig.Client(oauth1.NoContext, twitterToken)
	twitterClient := twitter.NewClient(httpClient)

	myMediaIds := []int64{mediaID}
	twitterClient.Statuses.Update(tweetText, &twitter.StatusUpdateParams{
		MediaIds: myMediaIds,
	})
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

//Gone, maybe used later
// func findRandomAnimalEmoji() string {
// 	animals := strings.Split(config.ImgSend.AnimalEmoji, ",")
// 	randSrc := rand.NewSource(time.Now().Unix())
// 	random := rand.New(randSrc)
// 	return animals[random.Intn(len(animals))]
// }

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
