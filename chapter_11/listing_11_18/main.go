package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func handleDirectories(dirs <-chan string, files chan<- string) {
	for fullPath := range dirs {
		fmt.Println("Reading all files from", fullPath)
		filesInDir, _ := os.ReadDir(fullPath)
		fmt.Printf("Pushing %d files from %s\n", len(filesInDir), fullPath)
		for _, file := range filesInDir {
			// Запустить новую горутину, которая отправляет
			// каждый файл в канал для файлов.
			go func(fp string) {
				files <- fp
			}(filepath.Join(fullPath, file.Name()))
		}
	}
}

func handleFiles(files chan string, dirs chan string) {
	for path := range files {
		file, _ := os.Open(path)
		fileInfo, _ := file.Stat()
		if fileInfo.IsDir() {
			fmt.Printf("Pushing %s directory\n", fileInfo.Name())
			dirs <- path
		} else {
			fmt.Printf("File %s, size: %dMB, last modified: %s\n",
				fileInfo.Name(), fileInfo.Size()/(1024*1024),
				fileInfo.ModTime().Format("15:04:05"))
		}
	}
}

func main() {
	filesChannel := make(chan string)
	dirsChannel := make(chan string)

	go handleFiles(filesChannel, dirsChannel)
	go handleDirectories(dirsChannel, filesChannel)

	dirsChannel <- os.Args[1]
	time.Sleep(60 * time.Second)
}
