package main

import (
	"strconv"
	"strings"

	"github.com/globalsign/mgo/bson"
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

//Handlers From Telegram
func findSubOrUnsubKeyboard(chatID int64) [][]tgAPI.InlineButton {
	var keyboardToSend [][]tgAPI.InlineButton
	if users.itemExists(bson.M{"chatId": chatID}) == true {
		keyboardToSend = keyboards["mainUnsub"]
	} else {
		keyboardToSend = keyboards["mainSub"]
	}
	return keyboardToSend
}

func handleChatUser(c *tgAPI.Callback) bool {
	chatMember, err := bot.ChatMemberOf(c.Message.Chat, c.Sender)
	if err != nil {
		handleBtnClick("An Error Occured", keyboards["back"], c)
		handleErr(err)
		return false
	}
	isAdmin := checkForAdmin(chatMember)
	if isAdmin == false {
		handleBtnClick(config.MainBot.GroupNotAdminMsg, keyboards["back"], c)
		return false
	}
	return true
}

func checkForAdmin(chatMember *tgAPI.ChatMember) bool {
	if chatMember.Role == "creator" || chatMember.Role == "administrator" {
		return true
	}
	return false
}

func handleStart(msg *tgAPI.Message) {
	if msg.Chat.Type == "channel" || msg.Chat.Type == "privatechannel" {
		return
	}
	bot.Send(&tgAPI.User{
		ID: int(msg.Chat.ID),
	}, config.MainBot.WelcomeMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: findSubOrUnsubKeyboard(msg.Chat.ID),
	})
}

func handleGroupAdd(msg *tgAPI.Message) {
	bot.Send(msg.Chat, config.MainBot.GroupAddMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: findSubOrUnsubKeyboard(msg.Chat.ID),
	})
}

func handleMigration(from, to int64) {
	if users.itemExists(bson.M{"chatId": from}) == false {
		return
	}
	users.update(bson.M{"chatId": from}, bson.M{"$set": bson.M{"chatId": to}})
}

//Handalers For Keybaord
func handleSubBtn(c *tgAPI.Callback) {
	if c.Message.FromGroup() == true {
		if shouldContinue := handleChatUser(c); shouldContinue == false {
			return
		}
	}
	status := handleSub(c.Message)
	if status == true {
		handleBtnClick(config.MainBot.SubMsg, keyboards["back"], c)
		owners := strings.Split(config.MainBot.Owners, ",")
		for i := range owners {
			idToSend, _ := strconv.Atoi(owners[i])
			bot.Send(&tgAPI.User{
				ID: idToSend,
			}, c.Message.Chat.Username+" Subscribed!")
		}
	} else {
		handleBtnClick(config.MainBot.AlreadySubMsg, keyboards["back"], c)
	}
}

func handleUnsubBtn(c *tgAPI.Callback) {
	if c.Message.FromGroup() == true {
		if shouldContinue := handleChatUser(c); shouldContinue == false {
			return
		}
	}
	status := handleUnsub(c.Message)
	if status == true {
		handleBtnClick(config.MainBot.UnsubMsg, keyboards["back"], c)
	} else {
		handleBtnClick(config.MainBot.NotSubMsg, keyboards["back"], c)
	}
}

func handleCommandBtn(c *tgAPI.Callback) {
	handleBtnClick(config.MainBot.CmdMsg, keyboards["cmd"], c)
}

func handleHomeBtn(c *tgAPI.Callback) {
	handleBtnClick(config.MainBot.WelcomeMsg, findSubOrUnsubKeyboard(c.Message.Chat.ID), c)
}

func handleInfoBtn(c *tgAPI.Callback) {
	var totalUsers []user
	users.findAll(bson.M{}, &totalUsers)
	sendString := config.MainBot.InfoMsg + "\n\nUsers Subscribed: " + strconv.Itoa(len(totalUsers))
	handleBtnClick(sendString, keyboards["back"], c)
}

func handleDaysBtn(c *tgAPI.Callback) {
	dayStr := strconv.Itoa(getDays(config.Date)) + " Days Until " + config.Con + "!"
	handleBtnClick(dayStr, keyboards["back"], c)
}

func handleSub(msg *tgAPI.Message) bool {
	if users.itemExists(bson.M{"chatId": msg.Chat.ID}) == true {
		return false
	}
	itemToInsert := user{
		ChatID: int(msg.Chat.ID),
		Group:  msg.FromGroup(),
	}
	users.insert(itemToInsert)
	return true
}

func handleUnsub(msg *tgAPI.Message) bool {
	if users.itemExists(bson.M{"chatId": msg.Chat.ID}) == false {
		return false
	}
	users.removeOne(bson.M{"chatId": msg.Chat.ID})
	return true
}
