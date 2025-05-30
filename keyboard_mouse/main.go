package main

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

func main() {
	low()
}

func low() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
	}
}
