package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

func sendPhoto(img finalImg) error {
	sendPhoto := tgAPI.Photo{
		File:    tgAPI.FromReader(img.ImgReader),
		Caption: intToEmoji(img.DaysLeft) + " Days Until " + config.Con + "! " + findRandomAnimalEmoji(),
	}
	var allUsers []user
	users.findAll(bson.M{}, &allUsers)

	for i := range allUsers {
		user := allUsers[i]
		tgUser := tgAPI.User{
			ID: user.ChatID,
		}
		_, err := sendPhoto.Send(bot, &tgUser, nil)
		if err != nil {
			if err.Error() == "api error: Forbidden: bot was blocked by the user" {
				users.removeOne(bson.M{"chatId": user.ChatID})
			}
		}
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
			break
		case "1":
			finalString += string("\u0031\uFE0F\u20E3")
			break
		case "2":
			finalString += string("\u0032\uFE0F\u20E3")
			break
		case "3":
			finalString += string("\u0033\uFE0F\u20E3")
			break
		case "4":
			finalString += string("\u0034\uFE0F\u20E3")
			break
		case "5":
			finalString += string("\u0035\uFE0F\u20E3")
			break
		case "6":
			finalString += string("\u0036\uFE0F\u20E3")
			break
		case "7":
			finalString += string("\u0037\uFE0F\u20E3")
			break
		case "8":
			finalString += string("\u0038\uFE0F\u20E3")
			break
		case "9":
			finalString += string("\u0039\uFE0F\u20E3")
			break
		}
	}
	return finalString
}

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
