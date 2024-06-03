package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func handleDirectories(dirs <-chan string, files chan<- string) {
	// Создать срез, чтобы сохранять файлы, которые нужно
	// добавить в канал для обработки файлов.
	toPush := make([]string, 0)

	appendAllFiles := func(path string) {
		fmt.Println("Reading all files from", path)
		filesInDir, _ := os.ReadDir(path)
		// Добавлять все файлы из директории в срез.
		for _, f := range filesInDir {
			toPush = append(toPush, filepath.Join(path, f.Name()))
		}
	}

	for {
		// Если нет файлов на добавление, прочитать директорию
		// из входящего канала и добавить все файлы из директории.
		if len(toPush) == 0 {
			appendAllFiles(<-dirs)
		} else {
			select {
			// Читать следующую директорию из входящего канала
			// и добавлять все файлы из директории.
			case fullPath := <-dirs:
				appendAllFiles(fullPath)
			// Добавляет первый файл из среза в канал.
			case files <- toPush[0]:
				// Удаляет первый файл из среза.
				toPush = toPush[1:]
			}
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
