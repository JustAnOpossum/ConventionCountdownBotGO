package helper

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//Credit Credit Struct
type Credit struct {
	_id   bson.ObjectId
	Photo string
	Used  bool
	Name  string
}

//User user Struct
type User struct {
	_id    bson.ObjectId
	ChatID int64 `bson:"chatId"`
	Name   string
	Group  bool
}

// Datastore Is the Handlaer for methods
type Datastore struct {
	session *mgo.Session
}

//FindOne Finds one item
func (datastore *Datastore) FindOne(collectionName string, query bson.M, result interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Find(query).One(result)
}

//FindAll Finds all items
func (datastore *Datastore) FindAll(collectionName string, query bson.M, results interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Find(query).All(results)
}

//Insert Inserts item
func (datastore *Datastore) Insert(collectionName string, itemToIntert interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Insert(itemToIntert)
}

//Update Updates item
func (datastore *Datastore) Update(collectionName string, query, itemToUpdate bson.M) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Update(query, itemToUpdate)
}

//RemoveOne Removes one item
func (datastore *Datastore) RemoveOne(collectionName string, query bson.M) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Remove(query)
}

//ItemExists See if item is there
func (datastore *Datastore) ItemExists(collectionName string, query bson.M) bool {
	data := datastore.session.Copy()
	defer data.Close()
	var result []interface{}
	data.DB("").C(collectionName).Find(query).All(&result)
	if len(result) == 0 {
		return false
	}
	return true
}

//SetUpDB sets up Database
func SetUpDB(dbName string) *Datastore {
	session, err := mgo.Dial("localhost/" + dbName)
	if err != nil {
		panic(err)
	}

	genIndex := func(keys []string) mgo.Index {
		return mgo.Index{
			Key:        keys,
			Unique:     true,
			Background: false,
			Sparse:     true,
		}
	}

	statbotSession := session.Copy()
	defer statbotSession.Close()

	if err = statbotSession.DB("").C("credits").EnsureIndex(genIndex([]string{"photo", "name"})); err != nil {
		panic(err)
	}
	if err = statbotSession.DB("").C("users").EnsureIndex(genIndex([]string{"chatId"})); err != nil {
		panic(err)
	}
	return &Datastore{session: session}
}
