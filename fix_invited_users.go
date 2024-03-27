package main

import "fmt"

type Team struct {
	ID          string       `json:"_id"`
	Admin       string       `json:"admin"`
	LicenseInfo LicenseInfo  `json:"licenseInfo"`
	Members     []TeamMember `json:"members"`
}

type TeamMember struct {
	User     string `json:"user"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
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
