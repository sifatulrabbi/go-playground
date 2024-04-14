package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type CityEntry struct {
	CityAscii string `json:"city_ascii"`
	Lat       string `json:"lat"`
	Lon       string `json:"lng"`
	Country   string `json:"country"`
}

type CCityEntry struct {
	Name    string `json:"name"`
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
	Country string `json:"country"`
}

type CountryEntry map[string][]CCityEntry

func GroupCitiesByCountry() {
	list := CountryEntry{}
	b, err := os.ReadFile("./assets/worldcities.json")
	if err != nil {
		log.Fatalln(err)
	}
	data := []CityEntry{}
	if err := json.Unmarshal(b, &data); err != nil {
		log.Fatalln(err)
	}
	for _, d := range data {
		entry := CCityEntry{
			Name:    d.CityAscii,
			Lat:     d.Lat,
			Lon:     d.Lon,
			Country: d.Country,
		}
		if _, ok := list[entry.Country]; !ok {
			list[entry.Country] = []CCityEntry{entry}
		} else {
			exists := false
			for _, c := range list[entry.Country] {
				if c.Name == entry.Name {
					exists = true
					break
				}
			}
			if !exists {
				list[entry.Country] = append(list[entry.Country], entry)
			}
		}
	}
	if b, err = json.MarshalIndent(list, "", "  "); err != nil {
		log.Fatalln(err)
	} else if err := os.WriteFile("./assets/citieslist.json", b, 0644); err != nil {
		log.Fatalln(err)
	}
}

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
