package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
)

func deleteDuplicatesFromTogai() {
	getAllTogaiCustomers()

	customers := ParseCustomsersCSVs()
	togaiCustomers := []TogaiCustomer{}
	customerIds := []string{}
	duplicates := map[string][]string{}
	customersToRemove := []string{}

	fileContent, err := os.ReadFile("all-togai-customers.json")
	if err != nil {
		log.Fatalln(err)
	}
	if err = json.Unmarshal(fileContent, &togaiCustomers); err != nil {
		log.Fatalln("unable to parse togai customer info:", err)
	}
	fmt.Printf("Customers need to be removed: %d\n",
		len(togaiCustomers)-len(customers))

	for i := 0; i < len(customers); i++ {
		customerIds = append(customerIds, customers[i].CustomerID)
	}

	for i := 0; i < len(togaiCustomers); i++ {
		tgcus := togaiCustomers[i]
		if _, exist := duplicates[tgcus.PrimaryEmail]; !exist {
			duplicates[tgcus.PrimaryEmail] = []string{tgcus.ID}
		} else {
			duplicates[tgcus.PrimaryEmail] = append(duplicates[tgcus.PrimaryEmail], tgcus.ID)
		}
	}
	fmt.Printf("total unique customers: %d. customers to remove: %d\n",
		len(duplicates), len(togaiCustomers)-len(duplicates))

	for _, cusIds := range duplicates {
		if len(cusIds) < 2 {
			continue
		}
		for i := 0; i < len(cusIds); i++ {
			if !slices.Contains(customerIds, cusIds[i]) {
				customersToRemove = append(customersToRemove, cusIds[i])
			}
		}
	}

	fmt.Printf("customers to remove: %d\n", len(customersToRemove))

	for i := 0; i < len(customersToRemove); i++ {
		tgcus := TogaiCustomer{ID: customersToRemove[i]}
		tgcus.Delete()
	}
}
