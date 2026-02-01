package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

func handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("Error while closing the connection:", err)
		}
	}()

	addr := conn.RemoteAddr()
	fmt.Printf("New connection from: %s\n", addr.String())

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

func mockHTTPServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		go mockHTTPServer()

		time.Sleep(time.Second * 30)
		wg.Done()
	}()

	time.Sleep(time.Second * 2)

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to:", conn.RemoteAddr())

	if _, err = conn.Write([]byte("JSON:")); err != nil {
		fmt.Println("unable to write to the open connection!")
	}

	if err = conn.Close(); err != nil {
		fmt.Println("Failed to gracefully close the connection.", err)
	}

	wg.Wait()
}
