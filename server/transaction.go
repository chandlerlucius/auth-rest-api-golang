package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User is a simple object to store authentication properties
type User struct {
	Created   time.Time
	LastLogin time.Time
	Username  string
	Password  string `json:"-"`
}

var ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
var mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
var database = "authentication"

func addUserToDb(username string, hashedPassword string) (*mongo.InsertOneResult, error) {
	user := User{time.Now(), time.Now(), username, hashedPassword}
	res, err := insertIntoDb(user, "user")
	return res, err
}

func findUserInDb(username string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	mongoCollection := mongoClient.Database(database).Collection("user")
	err := mongoCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return user, err
}

func insertIntoDb(data interface{}, collection string) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoCollection := mongoClient.Database(database).Collection(collection)
	res, err := mongoCollection.InsertOne(ctx, data)
	return res, err
}
