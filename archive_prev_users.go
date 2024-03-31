package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           primitive.ObjectID `json:"_id"`
	Admin        primitive.ObjectID `json:"admin"`
	Subscription primitive.ObjectID `json:"subscription"`
	Email        string             `json:"email"`
	Fullname     string             `json:"fullname"`
	Sumoling     bool               `json:"sumoling"`
}

func FindAndListAllArchivedUsers(db *mongo.Database) {
	// db.yourCollectionName.find({
	//   email: { $regex: '^archived__', $options: 'i' }
	// })
	// filter := bson.M{"email": bson.M{"$regex": "^archived__", "$options": "i"}}
	filter := bson.M{"archived": true}
	cur, err := db.Collection("users").Find(context.TODO(), filter)
	if err != nil {
		log.Fatalln(err)
	}
	archivedEmails := []string{}
	for cur.Next(context.TODO()) {
		user := User{}
		if err := cur.Decode(&user); err != nil {
			log.Fatalln("Unable to decode user")
		}
		archivedEmails = append(archivedEmails, user.Email)
	}
	csvContent := strings.Join(archivedEmails, "\n")
	if err := os.WriteFile(".cache/archived-users-list.csv", []byte(csvContent), 0644); err != nil {
		log.Fatalln("Unable to write the file", err)
	}
}

func ArchivePrevUser(db *mongo.Database) {
	coll := db.Collection("users")
	cur, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalln("Unable to create cursor for the query.", err)
	}
	archiveList := []string{}
	validList := []string{}
	for cur.Next(context.TODO()) {
		user := User{}
		if err := cur.Decode(&user); err != nil {
			log.Fatalln("error while decoding doc", err)
		}
		if !user.Admin.IsZero() || !user.Subscription.IsZero() || user.Sumoling {
			validList = append(validList, user.Email)
			continue
		}
		archiveList = append(archiveList, user.Email)
		filter := bson.M{"email": user.Email}
		update := bson.M{"$set": bson.M{
			"email":    "archived__" + user.Email,
			"archived": true,
		}}
		if err := coll.FindOneAndUpdate(context.TODO(), filter, update).Err(); err != nil {
			log.Fatalf("unable to update %v\n%v\n", user.Email, err.Error())
		}
		fmt.Println("archived__" + user.Email)
	}
	fmt.Printf("\nvalid users: %d\narchived users: %d\n", len(validList), len(archiveList))
	fmt.Println("Done!")
}
