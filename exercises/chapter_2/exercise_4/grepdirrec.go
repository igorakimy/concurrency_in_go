package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func grepFile(path string, entry os.DirEntry, search string) {
	fullPath := filepath.Join(path, entry.Name())

	if entry.IsDir() {
		files, err := os.ReadDir(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			go grepFile(fullPath, file, search)
		}
	} else {
		content, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatal(err)
		}

		matched, _ := regexp.Match(search, content)

		if matched {
			fmt.Printf("%s filename contains a match\n", fullPath)
		}
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

	for _, file := range files {
		go grepFile(dirname, file, word)
	}

	time.Sleep(1 * time.Second)
}
