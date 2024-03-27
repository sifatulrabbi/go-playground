package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetValidUsers(db *mongo.Database) CustomerList {
	cl := CustomerList{}
	payingUsers := ParseCustomsersCSVs()
	for _, u := range payingUsers {
		userColl := db.Collection("users")
		pipeline := mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.D{{Key: "email", Value: u.Email}}},
			},
			bson.D{
				{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 1},
					{Key: "email", Value: 1},
					{Key: "team", Value: 1},
				}},
			},
			bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "teams"},
					{Key: "localField", Value: "team"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "teamInfo"},
				}},
			},
			bson.D{
				{Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$teamInfo"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				}},
			},
		}
		cur, err := userColl.Aggregate(context.TODO(), pipeline)
		if err != nil {
			log.Fatalf("unable to get: %s\n", err)
		}
		defer cur.Close(context.TODO())
		results := []bson.M{}
		if err = cur.All(context.TODO(), &results); err != nil {
			log.Fatalln("error while decoding the user", err)
		}
		if len(results) < 1 {
			continue
		}

		email, _ := results[0]["email"].(string)
		if !slices.Contains(cl.Emails(), email) {
			cl = append(cl, Customer{Email: email, CustomerID: u.CustomerID, PlanName: u.PlanName})
		}
		fmt.Printf("%s | ", email)

		team, ok := results[0]["teamInfo"].(bson.M)
		if !ok {
			if b, err := json.MarshalIndent(team, "", "    "); err != nil {
				log.Fatalln("invalid team doc:", results[0]["teamInfo"])
			} else {
				fmt.Println(string(b))
				log.Fatalln("invalid team doc")
			}
		}
		members, ok := team["members"].(primitive.A)
		if !ok {
			log.Fatalln("invalid members list")
		}
		fmt.Printf("team len: %d\n", len(members))
		for _, m := range members {
			member, _ := m.(bson.M)
			email, _ := member["email"].(string)
			if !slices.Contains(cl.Emails(), email) {
				cl = append(cl, Customer{Email: email, PlanName: "Invited"})
			}
			fmt.Printf("- %s\n", email)
		}
	}

	return cl
}
