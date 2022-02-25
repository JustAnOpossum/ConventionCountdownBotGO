//tg_test.go
//Makes sure some of the core functions for the bot work as intended. eg. checking to make sure that users are able to subscribe and unsubscribe.

package main

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	tgAPI "gopkg.in/telebot.v3"
)

func TestSetUpDB(t *testing.T) {
	if photos == nil || users == nil {
		users, photos = setUpDB("localhost", "test")
	}
	users.removeAll()
	photos.removeAll()
}

func TestCheckForAdmin(t *testing.T) {
	testMember := &tgAPI.ChatMember{
		Role: "creator",
	}
	if result := checkForAdmin(testMember); result != true {
		t.Error("Got False")
	}

	testMember = &tgAPI.ChatMember{
		Role: "administrator",
	}
	if result := checkForAdmin(testMember); result != true {
		t.Error("Got False")
	}

	testMember = &tgAPI.ChatMember{
		Role: "member",
	}
	if result := checkForAdmin(testMember); result != false {
		t.Error("Got True")
	}
}

func TestSubscibeGroup(t *testing.T) {
	TestSetUpDB(t)
	inputMsgGroup := &tgAPI.Message{
		Chat: &tgAPI.Chat{
			ID:       123,
			Username: "Test",
			Type:     tgAPI.ChatGroup,
		},
	}

	if result := handleSub(inputMsgGroup.Chat.ID, true); result == false {
		t.Error("Result Is False")
	}
	var item user
	users.findOne(bson.M{"chatId": inputMsgGroup.Chat.ID}, &item)
	if item.ChatID != 0 && item.Group == false {
		t.Error("Chat Id Not Correct or Group Not Correct")
		t.Error(item.ChatID)
		t.Error(item.Group)
	}

	if result := handleSub(inputMsgGroup.Chat.ID, true); result == true {
		t.Error("Result Is True")
	}
}

func TestSubscibeUser(t *testing.T) {
	TestSetUpDB(t)
	inputMsgUser := &tgAPI.Message{
		Chat: &tgAPI.Chat{
			ID:       123,
			Username: "Test",
			Type:     tgAPI.ChatPrivate,
		},
	}

	if result := handleSub(inputMsgUser.Chat.ID, false); result == false {
		t.Error("Result Is False")
	}
	var item user
	users.findOne(bson.M{"chatId": inputMsgUser.Chat.ID}, &item)
	if item.ChatID != 0 && item.Group == true {
		t.Error("Chat Id Not Correct or Group Not Correct")
		t.Error(item.ChatID)
		t.Error(item.Group)
	}

	if result := handleSub(inputMsgUser.Chat.ID, false); result == true {
		t.Error("Result Is True")
	}
}

func TestUnsubscribe(t *testing.T) {
	TestSetUpDB(t)
	fakeMsg := &tgAPI.Message{
		Chat: &tgAPI.Chat{
			ID:       123,
			Username: "Test",
			Type:     tgAPI.ChatPrivate,
		},
	}

	if result := handleUnsub(fakeMsg.Chat.ID); result == true {
		t.Error("Result is True")
	}

	handleSub(fakeMsg.Chat.ID, false)

	if result := handleUnsub(fakeMsg.Chat.ID); result == false {
		t.Error("Result is false")
	}
	if result := users.itemExists(bson.M{"chatId": fakeMsg.Chat.ID}); result == true {
		t.Error("Item Not Deleted")
	}
}

// func TestHandleMigration(t *testing.T) {
// 	TestSetUpDB(t)
// 	inputMsgGroup := &tgAPI.Message{
// 		Chat: &tgAPI.Chat{
// 			ID:       123,
// 			Username: "Test",
// 			Type:     tgAPI.ChatGroup,
// 		},
// 	}

// 	ctx := tgAPI.Context{}

// 	handleSub(inputMsgGroup.Chat.ID, true)

// 	handleMigration(123, 456)
// 	if result := users.itemExists(bson.M{"chatId": 456}); result == false {
// 		t.Error("Update Failed")
// 	}
// }
