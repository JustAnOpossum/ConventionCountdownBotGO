package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type credit struct {
	_id   bson.ObjectId
	Photo string
	Used  bool
	Name  string
}

type user struct {
	_id    bson.ObjectId
	ChatID int64 `bson:"chatId"`
	Name   string
	Group  bool
}

type datastore struct {
	session *mgo.Session
}

func (datastore *datastore) findOne(collectionName string, query bson.M, result interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Find(query).One(result)
}

func (datastore *datastore) findAll(collectionName string, query bson.M, results interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Find(query).All(results)
}

func (datastore *datastore) insert(collectionName string, itemToIntert interface{}) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Insert(itemToIntert)
}

func (datastore *datastore) update(collectionName string, query, itemToUpdate bson.M) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Update(query, itemToUpdate)
}

func (datastore *datastore) removeOne(collectionName string, query bson.M) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).Remove(query)
}

func (datastore *datastore) removeAll(collectionName string, query bson.M) {
	data := datastore.session.Copy()
	defer data.Close()

	data.DB("").C(collectionName).RemoveAll(query)
}

func (datastore *datastore) itemExists(collectionName string, query bson.M) bool {
	data := datastore.session.Copy()
	defer data.Close()
	var result []interface{}
	data.DB("").C(collectionName).Find(query).All(&result)
	if len(result) == 0 {
		return false
	}
	return true
}

func setUpDB(dbName string) *datastore {
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
	return &datastore{session: session}
}
