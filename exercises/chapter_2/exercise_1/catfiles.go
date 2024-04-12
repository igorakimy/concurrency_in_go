package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func readFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(content))
}

func main() {
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		go readFile(args[i])
	}

	time.Sleep(1 * time.Second)
}
