package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatalln("MONGODB_URI env var not found but is required.")
	}
	db, disconnectFn := connectToDb(uri)
	defer disconnectFn()
	fmt.Println(db.Name())

	teams := []Team{}
	b, err := os.ReadFile(".cache/team-need-updates.json")
	if err != nil {
		log.Fatalln("unable to read file", err)
	}
	if err = json.Unmarshal(b, &teams); err != nil {
		log.Fatalln("unable to parse team data", err)
	}

	statusList := map[string]string{}
	wg := sync.WaitGroup{}
	for _, t := range teams {
		for _, m := range t.Members {
			statusList[m.Status] = m.Status
		}
		wg.Add(1)
		go func() {
			FixInvitedUsers(t)
			wg.Done()
		}()
	}
	wg.Wait()
}

func getUser(db *mongo.Database, email string) primitive.M {
	coll := db.Collection("users")
	user := bson.M{}
	filter := bson.D{{Key: "email", Value: email}}
	if err := coll.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		log.Fatalln("Error while finding the user:", err)
	}
	return user
}

func getSubscription(db *mongo.Database, email string) bson.M {
	coll := db.Collection("subscriptions")
	sub := bson.M{}
	filter := bson.D{{Key: "userEmail", Value: email}}
	if err := coll.FindOne(context.TODO(), filter).Decode(&sub); err != nil {
		log.Fatalln("Error while finding the user:", err)
	}
	return sub
}

func getUserTeamByEmail(db *mongo.Database, email string) bson.M {
	user := getUser(db, email)
	coll := db.Collection("teams")
	team := bson.M{}
	filter := bson.D{{Key: "_id", Value: user["team"]}}
	if err := coll.FindOne(context.TODO(), filter).Decode(&team); err != nil {
		log.Fatalln("Error while finding the user:", err)
	}
	return team
}

func connectToDb(uri string) (*mongo.Database, func()) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalln(err)
	}
	disconnect := func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalln("Error while disconnecting\n", err)
		}
	}

	db := client.Database("helloscribe-ai")
	return db, disconnect
}
