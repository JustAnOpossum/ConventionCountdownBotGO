package main

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

func TestSetUpDB(t *testing.T) {
	if db == nil {
		testDB := setUpDB("testing")
		db = testDB
	}
	db.removeAll("users", bson.M{})
	db.removeAll("credits", bson.M{})
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

	if result := handleSub(inputMsgGroup); result == false {
		t.Error("Result Is False")
	}
	var item user
	db.findOne("users", bson.M{"chatId": inputMsgGroup.Chat.ID}, &item)
	if item.ChatID != 0 && item.Group == false {
		t.Error("Chat Id Not Correct or Group Not Correct")
		t.Error(item.ChatID)
		t.Error(item.Group)
	}

	if result := handleSub(inputMsgGroup); result == true {
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

	if result := handleSub(inputMsgUser); result == false {
		t.Error("Result Is False")
	}
	var item user
	db.findOne("users", bson.M{"chatId": inputMsgUser.Chat.ID}, &item)
	if item.ChatID != 0 && item.Group == true {
		t.Error("Chat Id Not Correct or Group Not Correct")
		t.Error(item.ChatID)
		t.Error(item.Group)
	}

	if result := handleSub(inputMsgUser); result == true {
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

	if result := handleUnsub(fakeMsg); result == true {
		t.Error("Result is True")
	}

	handleSub(fakeMsg)

	if result := handleUnsub(fakeMsg); result == false {
		t.Error("Result is false")
	}
	if result := db.itemExists("users", bson.M{"chatId": fakeMsg.Chat.ID}); result == true {
		t.Error("Item Not Deleted")
	}
}

func TestHandleMigration(t *testing.T) {
	TestSetUpDB(t)
	inputMsgGroup := &tgAPI.Message{
		Chat: &tgAPI.Chat{
			ID:       123,
			Username: "Test",
			Type:     tgAPI.ChatGroup,
		},
	}
	handleSub(inputMsgGroup)

	handleMigration(123, 456)
	if result := db.itemExists("users", bson.M{"chatId": 456}); result == false {
		t.Error("Update Failed")
	}
}
