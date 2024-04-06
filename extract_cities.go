package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func ExtractCities() {
	b, err := os.ReadFile("./assets/worldcities.csv")
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(b), "\n")
	locations := []map[string]string{}
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		values := strings.Split(line, ",")
		loc := map[string]string{
			"city_ascii": strings.ReplaceAll(values[1], "\"", ""),
			"lat":        strings.ReplaceAll(values[2], "\"", ""),
			"lng":        strings.ReplaceAll(values[3], "\"", ""),
			"country":    strings.ReplaceAll(values[4], "\"", ""),
		}
		locations = append(locations, loc)
	}
	b, err = json.MarshalIndent(locations, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	if err = os.WriteFile("assets/worldcities.json", b, 0644); err != nil {
		log.Fatalln(err)
	}
}
