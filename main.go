package main

import (
	"context"
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

	customers := ParseCustomsersCSVs()
	fmt.Printf("Total valid users: %d\n", len(customers))

	c := make(chan bson.M)
	wg := sync.WaitGroup{}

	go func() {
		wg.Wait()
		close(c)
	}()

	for _, cus := range customers[4:10] {
		if cus.PlanName == "HelloScribe Tier 2" || cus.PlanName == "Premium" {
			go func() {
				wg.Add(1)
				c <- getUserTeamByEmail(db, cus.Email)
			}()
		}
	}

	usersCount := 0
	teamsNeedUpdates := []bson.M{}
	for team := range c {
		usersCount++
		members := team["members"]
		if members, ok := members.(primitive.A); ok {
			for _, m := range members {
				fmt.Println(m)
			}
		}
		wg.Done()
	}

	fmt.Printf("total users count: %d. total teams need check: %d\n",
		usersCount, len(teamsNeedUpdates))
}

// func deleteDuplicatesFromTogai() {
// 	// getAllTogaiCustomers()
//
// 	customers := ParseCustomsersCSVs()
// 	togaiCustomers := []TogaiCustomer{}
// 	customerIds := []string{}
// 	duplicates := map[string][]string{}
// 	customersToRemove := []string{}
//
// 	fileContent, err := os.ReadFile("all-togai-customers.json")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	if err = json.Unmarshal(fileContent, &togaiCustomers); err != nil {
// 		log.Fatalln("unable to parse togai customer info:", err)
// 	}
// 	fmt.Printf("Customers need to be removed: %d\n",
// 		len(togaiCustomers)-len(customers))
//
// 	for i := 0; i < len(customers); i++ {
// 		customerIds = append(customerIds, customers[i].CustomerID)
// 	}
//
// 	for i := 0; i < len(togaiCustomers); i++ {
// 		tgcus := togaiCustomers[i]
// 		if _, exist := duplicates[tgcus.PrimaryEmail]; !exist {
// 			duplicates[tgcus.PrimaryEmail] = []string{tgcus.ID}
// 		} else {
// 			duplicates[tgcus.PrimaryEmail] = append(duplicates[tgcus.PrimaryEmail], tgcus.ID)
// 		}
// 	}
// 	fmt.Printf("total unique customers: %d. customers to remove: %d\n",
// 		len(duplicates), len(togaiCustomers)-len(duplicates))
//
// 	for _, cusIds := range duplicates {
// 		if len(cusIds) < 2 {
// 			continue
// 		}
// 		for i := 0; i < len(cusIds); i++ {
// 			if !slices.Contains(customerIds, cusIds[i]) {
// 				customersToRemove = append(customersToRemove, cusIds[i])
// 			}
// 		}
// 	}
//
// 	fmt.Printf("customers to remove: %d\n", len(customersToRemove))
//
// 	for i := 0; i < len(customersToRemove); i++ {
// 		tgcus := TogaiCustomer{ID: customersToRemove[i]}
// 		tgcus.Delete()
// 	}
// }
