package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func AutomateLogin() {
	data := map[string]string{
		"email":    "jerry@jerrytrisya.com",
		"name":     "Oxygen",
		"origin":   "http://localhost:8080",
		"nickname": "Oxygen",
		"sub":      "bypass|fasdkj3248fdsahfdfjas213qapqpa",
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3001/api/authenticate/bypass", bytes.NewReader(b))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("request failed", err)
	}
	defer res.Body.Close()

	resBody := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		log.Fatalln("unable to decode respones data", err)
	}
	if res.StatusCode != 200 {
		log.Fatalln("unable to authenticate", resBody)
	}
	b, _ = json.MarshalIndent(resBody, "", "  ")
	fmt.Println("logged in:", string(b))
}
