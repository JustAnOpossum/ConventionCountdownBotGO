//tgHandle.go
//Handles all callbacks for the bot. Including the slash command and inline keyboard buttons.

package main

import (
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	tgAPI "gopkg.in/telebot.v3"
)

// Handlers From Telegram
func findSubOrUnsubKeyboard(chatID int64) [][]tgAPI.InlineButton {
	var keyboardToSend [][]tgAPI.InlineButton
	if users.itemExists(bson.M{"chatId": chatID}) {
		keyboardToSend = keyboards["mainUnsub"]
	} else {
		keyboardToSend = keyboards["mainSub"]
	}
	return keyboardToSend
}

// Handles the /start command and returns an inline keyboard to the user
func handleStart(ctx tgAPI.Context) error {
	if ctx.Message().Chat.Type == "channel" || ctx.Message().Chat.Type == "privatechannel" {
		return nil
	}
	ctx.Send(config.MainBot.WelcomeMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: findSubOrUnsubKeyboard(ctx.Message().Chat.ID),
	})
	return nil
}

// Handles when the bot is added to a group.
// TODO: Unify above function into one.
func handleGroupAdd(ctx tgAPI.Context) error {
	ctx.Send(config.MainBot.GroupAddMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: findSubOrUnsubKeyboard(ctx.Message().Chat.ID),
	})
	return nil
}

// If a chat is migrated, handle migration by updating group ID in database
func handleMigration(ctx tgAPI.Context) error {
	to, from := ctx.Migration()
	if !users.itemExists(bson.M{"chatId": from}) {
		return nil
	}
	users.update(bson.M{"chatId": from}, bson.M{"$set": bson.M{"chatId": to}})
	return nil
}

//May be used at a later time
// //Checks if user using commmand is an owner. And sends a test image to check image genration, etc.
// func handleTest(ctx tgAPI.Context) error {
// 	owners := strings.Split(config.MainBot.Owners, ",")
// 	for i := range owners {
// 		ownerID, _ := strconv.Atoi(owners[i])
// 		if ownerID == int(ctx.Sender().ID) {
// 			img := createImg()
// 			photo := tgAPI.Photo{File: tgAPI.FromReader(img.ImgReader)}
// 			err := ctx.Send(photo)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 		}
// 	}
// 	return nil
// }

// Handles subscribe button event
func handleSubBtn(ctx tgAPI.Context) error {
	//Calls helper to make sure user has permissions to subscribe the group to the bot
	if ctx.Chat().Type == tgAPI.ChatGroup {
		if shouldContinue := handleChatUser(ctx); !shouldContinue {
			return nil
		}
	}

	status := handleSub(ctx.Chat().ID, ctx.Message().FromGroup())
	if status {
		handleBtnClick(config.MainBot.SubMsg, keyboards["back"], ctx)
		//Messages the owners in the config when subscribe is successful
		if config.MessageOnSub {
			owners := strings.Split(config.MainBot.Owners, ",")
			for i := range owners {
				idToSend, _ := strconv.Atoi(owners[i])
				bot.Send(&tgAPI.User{
					ID: int64(idToSend),
				}, ctx.Message().Chat.Username+" Subscribed!")
			}
		}
	} else {
		handleBtnClick("I'm sorry, you are already subscribed.", keyboards["back"], ctx)
	}
	return nil
}

// Handles an unsubscribe button click
func handleUnsubBtn(ctx tgAPI.Context) error {
	//Calls helper to make sure user has permissions to subscribe the group to the bot
	if ctx.Message().FromGroup() {
		if shouldContinue := handleChatUser(ctx); !shouldContinue {
			return nil
		}
	}
	status := handleUnsub(ctx.Message().Chat.ID)
	if status {
		handleBtnClick("Unsubscribed", keyboards["back"], ctx)
	} else {
		handleBtnClick("I'm sorry you are not unsubscribed", keyboards["back"], ctx)
	}
	return nil
}

// Button handle for command
func handleCommandBtn(ctx tgAPI.Context) error {
	handleBtnClick("Commands:", keyboards["cmd"], ctx)
	return nil
}

// Button handle for Home
func handleHomeBtn(ctx tgAPI.Context) error {
	handleBtnClick(config.MainBot.WelcomeMsg, findSubOrUnsubKeyboard(ctx.Message().Chat.ID), ctx)
	return nil
}

// Button handle for info
func handleInfoBtn(ctx tgAPI.Context) error {
	var totalUsers []user
	users.findAll(bson.M{}, &totalUsers)
	sendString := "Bot created by @JustAnOpossum" + "\n\nUsers Subscribed: " + strconv.Itoa(len(totalUsers))
	handleBtnClick(sendString, keyboards["back"], ctx)
	return nil
}

// Button handle for days
func handleDaysBtn(ctx tgAPI.Context) error {
	dayStr := strconv.Itoa(getDays(config.Date)) + " Days Until " + config.Con + "!"
	if getDays(config.Date) == 1 {
		dayStr = strconv.Itoa(getDays(config.Date)) + " Day Until " + config.Con + "!"
	}
	if getDays(config.Date) < 1 {
		dayStr = "Have fun at " + config.Con + "!"
	}
	handleBtnClick(dayStr, keyboards["back"], ctx)
	return nil
}

// Handles a user or group subscription to that database
func handleSub(chatID int64, isGroup bool) bool {
	if users.itemExists(bson.M{"chatId": chatID}) {
		return false
	}
	itemToInsert := user{
		ChatID: chatID,
		Group:  isGroup,
	}
	users.insert(itemToInsert)
	return true
}

func handleUnsub(chatID int64) bool {
	if !users.itemExists(bson.M{"chatId": chatID}) {
		return false
	}
	users.removeOne(bson.M{"chatId": chatID})
	return true
}

// Helper function for subscibing a user of a group to the bot
func handleChatUser(ctx tgAPI.Context) bool {
	chatMember, err := bot.ChatMemberOf(ctx.Chat(), ctx.Sender())
	if err != nil {
		handleBtnClick("An Error Occured", keyboards["back"], ctx)
		handleErr(err)
		return false
	}
	isAdmin := checkForAdmin(chatMember)
	if !isAdmin {
		handleBtnClick("You don't have permissoins to subscribe this group. Please ask a group admin or owner to subscribe the bot.", keyboards["back"], ctx)
		return false
	}
	return true
}

// Helper function to make sure the user trying to subscribe is an admin of the group.
func checkForAdmin(chatMember *tgAPI.ChatMember) bool {
	if chatMember.Role == "creator" || chatMember.Role == "administrator" {
		return true
	}
	return false
}
