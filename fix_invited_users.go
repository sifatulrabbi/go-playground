package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID          primitive.ObjectID `json:"_id"`
	Admin       primitive.ObjectID `json:"admin"`
	LicenseInfo LicenseInfo        `json:"licenseInfo"`
	Members     []TeamMember       `json:"members"`
}

type TeamMember struct {
	User     primitive.ObjectID `json:"user"`
	Email    string             `json:"email"`
	Fullname string             `json:"fullname"`
	// "Pending" | "Joined" | "Accepted"
	Status string `json:"status"`
}

type LicenseInfo struct {
	Billed int `json:"billed"`
	Free   int `json:"free"`
}

func FixInvitedUsers(team Team) {
	fmt.Printf("fixing team: %s\n", team.ID)
}
