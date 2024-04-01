package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TogaiLicenseEntry struct {
	Quantity int `json:"quantity"`
}

type TogaiLicenseResponse struct {
	Data      []TogaiLicenseEntry `json:"data"`
	NextToken string              `json:"nextToken"`
}

func getBilledLicenseCount(cusId string) int {
	req, err := http.NewRequest(http.MethodGet, "https://api.togai.com/license_updates", http.NoBody)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TOGAI_API_KEY")))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Set("license_id", "")
	q.Set("account_id", cusId)
	req.URL.RawQuery = q.Encode()

	fmt.Printf("%s | ", cusId)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("unable to complete the req", err)
	}
	if res.StatusCode == 429 {
		time.Sleep(time.Second * 1)
		return getBilledLicenseCount(cusId)
	}
	defer res.Body.Close()
	data := TogaiLicenseResponse{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		log.Fatalln("unable to decode the data", err)
	}
	fmt.Println(data)
	if len(data.Data) < 1 {
		return 0
	}
	return data.Data[0].Quantity
}

func FixTier2Users(db *mongo.Database) {
	usersNeedFix := []string{}
	coll := db.Collection("subscriptions")
	filter := bson.M{"planName": "HelloScribe Tier 2"}
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatalln(err)
	}
	for cur.Next(context.TODO()) {
		sub := bson.M{}
		if err := cur.Decode(&sub); err != nil {
			log.Fatalln(err)
		}
		cusId, _ := sub["stripeCustomerId"].(string)
		if q := getBilledLicenseCount(cusId); q > 0 {
			usersNeedFix = append(usersNeedFix, cusId)
		}
	}

	b, err := json.MarshalIndent(usersNeedFix, "", "  ")
	if err != nil {
		log.Fatalln("unable to decode the user list", err)
	}
	if err := os.WriteFile(".cache/teir2-user-with-lic-issue.json", b, 0644); err != nil {
		log.Fatalln("Unable to write in the file", err)
	}
}
