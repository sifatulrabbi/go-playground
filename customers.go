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

type CustomerList []Customer

func (cl CustomerList) Emails() []string {
	emails := []string{}
	for _, c := range cl {
		emails = append(emails, c.Email)
	}
	return emails
}

func (cl CustomerList) CusIds() []string {
	cusIds := []string{}
	for _, c := range cl {
		cusIds = append(cusIds, c.CustomerID)
	}
	return cusIds
}

func ParseCustomsersCSVs() CustomerList {
	lines := []string{}
	lines = append(lines, readAndParseUsers("paying-users-list.csv")...)
	lines = append(lines, readAndParseUsers("appsumo-users-list.csv")...)
	customers := CustomerList{}
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
