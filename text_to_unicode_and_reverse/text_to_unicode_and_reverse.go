package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getTestInput() string {
	f, err := os.Open("./test_input.txt")
	if err != nil {
		log.Fatalln("Unable to find the file\n", err)
	}
	buff := make([]byte, 1024)
	if _, err = f.Read(buff); err != nil {
		log.Fatalln("Unable to read the file\n", err)
	}
	return bytes.NewBuffer(buff).String()
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln("Please provide a text and the type")
	}

	var result string
	switch args[0] {
	case "test":
		result = strToUnicode(strings.Split(getTestInput(), " "))
		break
	case "-u":
		result = unicodeToStr(args[1:])
		break
	default:
		result = strToUnicode(args[0:])
		break
	}
	fmt.Println(result)
}

func strToUnicode(text []string) string {
	result := ""
	for _, r := range strings.Join(text, " ") {
		result += fmt.Sprintf("%d ", r)
	}
	return result
}

func unicodeToStr(unicodeText []string) string {
	result := ""
	for _, strRuneCode := range unicodeText {
		runeCode, _ := strconv.ParseInt(strRuneCode, 10, 32)
		result += fmt.Sprintf("%c", runeCode)
	}
	return result
}
