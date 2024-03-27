package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GatherUsersWithTeam(db *mongo.Database) {
	customers := ParseCustomsersCSVs()
	fmt.Printf("Total valid users: %d\n", len(customers))

	c := make(chan bson.M)
	wg := sync.WaitGroup{}

	go func() {
		wg.Wait()
		close(c)
	}()

	for _, cus := range customers {
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
		if members, ok := members.(primitive.A); ok && len(members) > 0 {
			teamsNeedUpdates = append(teamsNeedUpdates, team)
		}
		wg.Done()
	}

	fmt.Printf("total users count: %d. total teams need check: %d\n",
		usersCount, len(teamsNeedUpdates))

	b, err := json.MarshalIndent(teamsNeedUpdates, "", "    ")
	if err != nil {
		log.Fatalln("Unable to marshal team data", err)
	}
	if err = os.WriteFile(".cache/team-need-updates.json", b, 0644); err != nil {
		log.Fatalln("unable to write file", err)
	}
}
