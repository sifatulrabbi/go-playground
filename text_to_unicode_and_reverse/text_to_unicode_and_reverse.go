package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln("Please provide a text and the type")
	}

	unicodeInput := args[0] == "-u"

	if !unicodeInput {
		for _, r := range strings.Join(args[0:], " ") {
			fmt.Printf("%d ", r)
		}
		fmt.Println()
	} else {
		for _, strRuneCode := range args[1:] {
			runeCode, _ := strconv.ParseInt(strRuneCode, 10, 32)
			fmt.Printf("%c", runeCode)
		}
		fmt.Println()
	}
}
