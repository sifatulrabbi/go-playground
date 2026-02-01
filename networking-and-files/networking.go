package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func readFromConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Error while closing the connection:", err)
		}
	}()

	addr := conn.RemoteAddr()
	fmt.Printf("Reading from: %s, network: %s\n", addr.String(), addr.Network())

	content := []byte{}
	buff := make([]byte, 1024)

	for {
		n, err := conn.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Done getting all the content")
			} else {
				fmt.Println("Error during saving the incoming content.", err)
			}
			break
		}
		fmt.Printf("Reading from %s ...\n", addr)
		content = append(content, buff[:n]...)
	}

	fmt.Println("Total bytes read:", len(content))
}

func mockHTTPServer(ready chan<- struct{}) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Hello world!")); err != nil {
			fmt.Println("Error writing on the HTTP conn:", err)
		}
	})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	close(ready)

	if err := http.Serve(listener, mux); err != nil {
		fmt.Println(err)
	}
}

func main() {
	ready := make(chan struct{})

	go mockHTTPServer(ready)
	<-ready

	addr := ":8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect to address: %s | %s\n", addr, err)
	}

	go readFromConn(conn)

	sampleGetReqContent := "GET / HTTP/1.0\r\nHost: localhost\r\n\r\n"

	if _, err = conn.Write([]byte(sampleGetReqContent)); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Done writing HTTP request.")

	time.Sleep(time.Second * 30)
}
