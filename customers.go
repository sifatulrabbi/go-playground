package main

import (
	"log"
	"os"
	"strings"
)

type Customer struct {
	Email      string `json:"email"`
	PlanName   string `json:"plan_name"`
	CustomerID string `json:"customer_id"`
}

func ParseCustomsersCSVs() []Customer {
	lines := []string{}
	lines = append(lines, readAndParseUsers("paying-users-list.csv")...)
	lines = append(lines, readAndParseUsers("appsumo-users-list.csv")...)
	customers := []Customer{}
	for _, line := range lines {
		parts := strings.Split(line, ",")
		cus := Customer{
			Email:      parts[0],
			PlanName:   parts[1],
			CustomerID: parts[2],
		}
		customers = append(customers, cus)
	}
	return customers
}

func readAndParseUsers(file string) []string {
	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln("unable to read paying users:", err)
	}
	strContent := string(b)
	lines := strings.Split(strContent, "\n")[1:]
	lines = lines[:len(lines)-1]
	return lines
}
