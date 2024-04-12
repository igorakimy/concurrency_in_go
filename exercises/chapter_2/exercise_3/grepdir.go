package main

import (
	"fmt"
	"log"
	"os"
	"path"
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

	if len(args) < 2 {
		log.Fatal("you must specify directory name as second argument")
	}

	dirname := args[1]

	files, err := os.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(files); i++ {
		go grepFile(word, path.Join(dirname, files[i].Name()))
	}

	time.Sleep(1 * time.Second)
}
