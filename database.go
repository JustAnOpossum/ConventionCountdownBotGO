package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type photo struct {
	_id   bson.ObjectId
	Photo string
	Used  bool
	Name  string
	URL   string
}

type user struct {
	_id    bson.ObjectId
	ChatID int `bson:"chatId"`
	Name   string
	Group  bool
}

type datastore struct {
	session       *mgo.Session
	collectioName string
}

func (datastore *datastore) findOne(query bson.M, result interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(datastore.collectioName).Find(query).One(result)
}

func (datastore *datastore) findAll(query bson.M, results interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(datastore.collectioName).Find(query).All(results)
}

func (datastore *datastore) insert(itemToIntert interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(datastore.collectioName).Insert(itemToIntert)
}

func (datastore *datastore) update(query, itemToUpdate bson.M) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(datastore.collectioName).Update(query, itemToUpdate)
}

func (datastore *datastore) removeOne(query bson.M) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(datastore.collectioName).Remove(query)
}

func (datastore *datastore) removeAll(query bson.M) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(datastore.collectioName).RemoveAll(query)
}

func (datastore *datastore) itemExists(query bson.M) bool {
	data := datastore.session.Copy()
	defer data.Close()
	var result []interface{}
	data.DB("").C(datastore.collectioName).Find(query).All(&result)
	if len(result) == 0 {
		return false
	}
	return true
}

func (datastore *datastore) distinct(query bson.M, distinctKey string) []string {
	data := datastore.session.Copy()
	defer data.Close()
	var tempResult []string
	data.DB("").C(datastore.collectioName).Find(query).Distinct(distinctKey, &tempResult)
	return tempResult
}

func setUpDB(dbName string) (*datastore, *datastore) {
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

	if err = statbotSession.DB("").C("photos").EnsureIndex(genIndex([]string{"photo", "name"})); err != nil {
		panic(err)
	}
	if err = statbotSession.DB("").C("users").EnsureIndex(genIndex([]string{"chatId"})); err != nil {
		panic(err)
	}
	return &datastore{session: session, collectioName: "users"}, &datastore{session: session, collectioName: "photos"}
}
