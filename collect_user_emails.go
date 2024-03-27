package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUserEmails(db *mongo.Database) []Customer {
	findOpts := options.Find().SetProjection(bson.M{"email": 1, "_id": 1, "subscription": 1, "team": 1})
	cur, err := db.Collection("users").Find(context.TODO(), bson.M{}, findOpts)
	if err != nil {
		log.Fatalln("Error while getting users", err)
	}
	defer cur.Close(context.TODO())

	customers := []Customer{}
	for cur.Next(context.TODO()) {
		user := bson.M{}
		if err = cur.Decode(&user); err != nil {
			log.Fatalln("Unable to decode user", err)
		}
		email := user["email"]
		sub := user["subscription"]
		cus := Customer{}
		if email, ok := email.(string); ok {
			cus.Email = email
		} else {
			continue
		}
		if sub, ok := sub.(string); ok {
			cus.PlanName = sub
		}
		customers = append(customers, cus)
	}

	return customers
}
