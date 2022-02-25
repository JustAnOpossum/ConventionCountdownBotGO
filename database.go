//database.go
//Containts the driver code for the database. Wraps database functions inside a method so that underlying API can be changed easier.

package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type photo struct {
	Photo string
	Used  bool
	Name  string
	URL   string
}

type user struct {
	ChatID int64 `bson:"chatId"`
	Group  bool
}

//Data type for a collection
type datastore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (datastore *datastore) findOne(query bson.M, result interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	datastore.collection.FindOne(ctx, query).Decode(result)
}

func (datastore *datastore) findAll(query bson.M, result interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, _ := datastore.collection.Find(ctx, query)
	cursor.All(ctx, result)
}

func (datastore *datastore) insert(itemToIntert interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	datastore.collection.InsertOne(ctx, itemToIntert)
}

func (datastore *datastore) update(query, itemToUpdate bson.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	datastore.collection.UpdateOne(ctx, query, itemToUpdate)
}

func (datastore *datastore) removeOne(query bson.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	datastore.collection.DeleteOne(ctx, query)
}

func (datastore *datastore) removeAll() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	datastore.collection.DeleteMany(ctx, bson.D{})
}

func (datastore *datastore) itemExists(query bson.M) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var results []interface{}
	err := datastore.collection.FindOne(ctx, query).Decode(results)
	return err != mongo.ErrNoDocuments
}

func (datastore *datastore) distinct(query bson.M, distinctKey string) []interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tempResult []interface{}
	distinct, _ := datastore.collection.Distinct(ctx, distinctKey, bson.M{})
	tempResult = append(tempResult, distinct...)
	return tempResult
}

func setUpDB(dbName string, dbURL string) (*datastore, *datastore) {
	//Connects to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))

	if err != nil {
		panic(err)
	}

	usersCollection := client.Database(dbName).Collection("users")
	photosCollection := client.Database(dbName).Collection("photos")

	noRepeats := true

	//Creates indexes for the models so that everything will be unique
	usersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"chatId": 1}, Options: &options.IndexOptions{Unique: &noRepeats}})
	photosCollection.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"photo": 1, "name": 1}, Options: &options.IndexOptions{Unique: &noRepeats}})

	//Wraps the collections inside a wrapper so that the underlying API can be changed
	return &datastore{collection: usersCollection, client: client}, &datastore{collection: photosCollection, client: client}
}
