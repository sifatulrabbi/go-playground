package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FixAnnualUsers(db *mongo.Database) {
	emails := getAnnualUsers()
	for _, email := range emails {
		// get the user profile
		// get their subscription & team info
		// get their togai info

		userColl := db.Collection("users")
		pipeline := mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.D{{Key: "email", Value: email}}},
			},
			bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "subscriptions"},
					{Key: "localField", Value: "subscription"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "subscriptionDetails"},
				}},
			},
			bson.D{
				{Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$subscriptionDetails"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				}},
			},
		}
		cur, err := userColl.Aggregate(context.TODO(), pipeline)
		if err != nil {
			log.Fatalf("unable to get: %s\n", err)
		}
		defer cur.Close(context.TODO())

		users := []bson.M{}
		if err = cur.All(context.TODO(), &users); err != nil {
			log.Fatalln("error while decoding user", err)
		}

		fmt.Println(users[0]["subscriptionDetails"])
	}
}

func getAnnualUsers() []string {
	b, err := os.ReadFile(".cache/users.txt")
	if err != err {
		log.Fatalln("unable to read the file", err)
	}
	lines := strings.Split(string(b), "\n")
	return lines[:len(lines)-1]
}
