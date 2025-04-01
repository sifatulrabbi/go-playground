package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	filename := "target-file.csv"
	filename = "../"
	readAndPrintDirEntries(filename, 0)
	// readFileInChunks(filename)
}

func readFileInChunks(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	fmt.Println(f.Name())

	r := bufio.NewReader(f)
	for {
		buf := make([]byte, 4*1024)
		_, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Done reading file")
				break
			} else {
				log.Fatalln("Error during file reading:", err)
			}
		}
	}
}

func readAndPrintDirEntries(dirname string, prefixSize int) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		log.Fatalln(err)
	}
	for _, e := range entries {
		prefix := ""
		for i := 0; i < prefixSize; i++ {
			prefix += "  "
		}
		prefix += "-"
		fmt.Println(prefix, e.Name())
		if e.IsDir() {
			readAndPrintDirEntries(dirname+"/"+e.Name(), prefixSize+1)
		}
	}
}
