package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalln(err)
	// }
	// uri := os.Getenv("MONGODB_URI")
	// if uri == "" {
	// 	log.Fatalln("MONGODB_URI env var not found but is required.")
	// }

	GroupCitiesByCountry()
	// ExtractCities()
	// db, disconnectFn := connectToDb(uri)
	// defer disconnectFn()

	// AutomateLogin()

	// FixTier2Users(db)
	// ArchivePrevUser(db)
	// FindAndListAllArchivedUsers(db)
}

func findingValidUsers() {
	var (
		validUsers         []string
		invalidUsers       []string
		updatedInvalidList []string
		updatedInvalidLen  int
		prevInvalidLen     int
		updatedCsv         string
	)
	if b, err := os.ReadFile(".cache/valid/valid-user-email-list.csv"); err != nil {
		log.Fatalln("valid user email list file not found", err)
	} else {
		validUsers = strings.Split(string(b), "\n")
	}
	if b, err := os.ReadFile(".cache/invalid-users-list.csv"); err != nil {
		log.Fatalln("invalid user email list file not found", err)
	} else {
		invalidUsers = strings.Split(string(b), "\n")
	}

	prevInvalidLen = len(invalidUsers)
	for _, e := range invalidUsers {
		if slices.Contains(validUsers, e) {
			continue
		}
		updatedInvalidList = append(updatedInvalidList, e)
		updatedCsv += fmt.Sprintf("%s\n", e)
	}
	updatedInvalidLen = len(updatedInvalidList)

	if err := os.WriteFile(".cache/valid/valid-user-email-list.csv", []byte(updatedCsv), 0644); err != nil {
		log.Fatalln("Unable to update file content", err)
	}

	fmt.Printf("prev length: %d, updated length: %d\n", prevInvalidLen, updatedInvalidLen)
}

func gatherAndFilterUserEmails(db *mongo.Database) {
	validCusEmails := ParseCustomsersCSVs().Emails()
	allCustomers := GetAllUserEmails(db)
	notPayingCustomers := CustomerList{}
	removeQueue := CustomerList{}
	for _, cus := range allCustomers {
		if slices.Contains(validCusEmails, cus.Email) {
			continue
		}
		if slices.Contains(notPayingCustomers.Emails(), cus.Email) {
			continue
		}
		notPayingCustomers = append(notPayingCustomers, cus)
		if cus.PlanName == "" {
			removeQueue = append(removeQueue, cus)
		}
	}
	fmt.Printf("Total user docs found: %d. total not paying users: %d. total removable users: %d\n",
		len(allCustomers), len(notPayingCustomers), len(removeQueue))

	b, err := json.MarshalIndent(notPayingCustomers, "", "    ")
	if err != nil {
		log.Fatalln("Unable to marshal data", err)
	}
	if err = os.WriteFile(".cache/not-paying-users.json", b, 0644); err != nil {
		log.Fatalln("unable to write file", err)
	}

	notPayingUserEmails := ""
	for _, e := range notPayingCustomers.Emails() {
		notPayingUserEmails += fmt.Sprintf("%s,\n", e)
	}
	if err = os.WriteFile(".cache/not-paying-users-list.csv", []byte(notPayingUserEmails), 0644); err != nil {
		log.Fatalln("unable to write file", err)
	}

	invalidUserEmails := ""
	for _, e := range removeQueue.Emails() {
		invalidUserEmails += fmt.Sprintf("%s,\n", e)
	}
	if err = os.WriteFile(".cache/invalid-users-list.csv", []byte(invalidUserEmails), 0644); err != nil {
		log.Fatalln("unable to write file", err)
	}
}

func gatherTeamInfoOfTheUsers() {
	// teams := []Team{}
	// b, err := os.ReadFile(".cache/team-need-updates.json")
	// if err != nil {
	// 	log.Fatalln("unable to read file", err)
	// }
	// if err = json.Unmarshal(b, &teams); err != nil {
	// 	log.Fatalln("unable to parse team data", err)
	// }

	// FixAnnualUsers(db)

	// statusList := map[string]string{}
	// wg := sync.WaitGroup{}
	// for _, t := range teams {
	// 	for _, m := range t.Members {
	// 		statusList[m.Status] = m.Status
	// 	}
	// 	wg.Add(1)
	// 	go func() {
	// 		FixInvitedUsers(t)
	// 		wg.Done()
	// 	}()
	// }
	// wg.Wait()
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
