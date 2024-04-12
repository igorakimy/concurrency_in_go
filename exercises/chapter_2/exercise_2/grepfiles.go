package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

func grepFile(search string, filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	matched, _ := regexp.Match(search, content)

	if matched {
		fmt.Printf("%s filename contains a match\n", filename)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("you must specify search word at first argument")
	}

	word := args[0]
	fileNames := args[1:]

	for i := 0; i < len(fileNames); i++ {
		go grepFile(word, fileNames[i])
	}

	time.Sleep(1 * time.Second)
}
