// Package tdg does something
package tgd

import (
	"fmt"
	"log"
	"net/http"
)

func init() {
	LoadConfigs()
}

func StartAPI() {
	cb := NewChatHub()
	go cb.start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(cb, w, r)
	})

	addr := fmt.Sprintf(":%d", AppConfig.PORT)
	fmt.Printf("Starting the server; addr=%s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
