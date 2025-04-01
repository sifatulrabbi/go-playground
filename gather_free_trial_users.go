package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GatherFreeTrialUsers(db *mongo.Database) {
	coll := db.Collection("subscriptions")
	filter := bson.M{"planName": "Free Trial"}
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatalln(err)
	}
	list := []map[string]interface{}{}
	for cur.Next(context.TODO()) {
		doc := bson.M{}
		if err = cur.Decode(&doc); err != nil {
			log.Fatalln(err)
		}
		list = append(list, map[string]interface{}{
			"email":    doc["userEmail"],
			"joinedAt": doc["createdAt"],
			"userId":   doc["userId"],
			"planName": doc["planName"],
		})
	}
	b, err := json.MarshalIndent(list, "", "    ")
	if err = os.WriteFile("./.cache/free_trial_users.json", b, 0644); err != nil {
		log.Fatalln(err)
	}

	usersCsv := `"Email","Plan Name","Joined at"
`
	for _, entry := range list {
		usersCsv += fmt.Sprintf("\"%s\",\"%s\"\n", entry["email"], entry["planName"])
	}
	if err = os.WriteFile("./.cache/free_trial_users.csv", []byte(usersCsv), 0644); err != nil {
		log.Fatalln(err)
	}
}
