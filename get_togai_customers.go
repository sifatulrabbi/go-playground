package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type TogaiCustomer struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PrimaryEmail string `json:"primaryEmail"`
	Status       string `json:"status"`
}

type Response struct {
	Message   string          `json:"message"`
	Data      []TogaiCustomer `json:"data"`
	NextToken string          `json:"nextToken"`
}

func (cus TogaiCustomer) Delete() {
	accountDeleted, customerDeleted := false, false
	for !accountDeleted {
		if cus.deleteAccount() != 429 {
			accountDeleted = true
		}
	}
	for !customerDeleted {
		if cus.deleteCustomer() != 429 {
			customerDeleted = true
		}
	}
	fmt.Printf("removed %s\n", cus.ID)
}

func (cus TogaiCustomer) deleteAccount() int {
	togaiApiKey := os.Getenv("TOGAI_API_KEY")
	url := fmt.Sprintf("https://api.togai.com/accounts/%s", cus.ID)
	req, err := http.NewRequest(http.MethodDelete, url, http.NoBody)
	if err != nil {
		log.Fatalln("error while creating the request:", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", togaiApiKey))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error: account %s - %v\n", cus.PrimaryEmail, err)
	}
	defer res.Body.Close()
	return res.StatusCode
}

func (cus TogaiCustomer) deleteCustomer() int {
	togaiApiKey := os.Getenv("TOGAI_API_KEY")
	url := fmt.Sprintf("https://api.togai.com/customers/%s", cus.ID)
	req, err := http.NewRequest(http.MethodDelete, url, http.NoBody)
	if err != nil {
		log.Fatalln("error while creating the request:", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", togaiApiKey))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error: customer %s - %v\n", cus.PrimaryEmail, err)
	}
	defer res.Body.Close()
	return res.StatusCode
}

func getAllTogaiCustomers() {
	customers := []TogaiCustomer{}
	getTogaiCustomers(&customers, "")
	fmt.Println("Total customers found:", len(customers))

	data, err := json.MarshalIndent(customers, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	if err = os.WriteFile("all-togai-customers.json", data, 0644); err != nil {
		log.Fatalln(err)
	}
}

func getTogaiCustomers(customers *[]TogaiCustomer, nextToken string) {
	togaiApiKey := os.Getenv("TOGAI_API_KEY")
	req, err := http.NewRequest(http.MethodGet, "https://api.togai.com/customers", http.NoBody)
	if err != nil {
		log.Fatalln("error while creating the request:", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", togaiApiKey))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("nextToken", nextToken)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	if res.StatusCode == 429 {
		time.Sleep(1 * time.Second)
		getTogaiCustomers(customers, nextToken)
		return
	}

	result := Response{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalln("error while converting the response:", err)
	}
	*customers = append(*customers, result.Data...)
	if result.NextToken == "" {
		return
	}
	getTogaiCustomers(customers, result.NextToken)
}
