package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func handleDirectories(dirs <-chan string, files chan<- string) {
	// Читать полный путь директории со входного канала.
	for fullPath := range dirs {
		fmt.Println("Reading all files from", fullPath)
		// Прочитать содержимое директории.
		filesInDir, _ := os.ReadDir(fullPath)
		fmt.Printf("Pushing %d files from %s\n", len(filesInDir), fullPath)
		// Подавать каждый элемент директории в выходящий канал.
		for _, file := range filesInDir {
			files <- filepath.Join(fullPath, file.Name())
		}
	}
}

func handleFiles(files chan string, dirs chan string) {
	// Читать полный путь файла.
	for path := range files {
		file, _ := os.Open(path)
		// Читать информацию о файле.
		fileInfo, _ := file.Stat()
		if fileInfo.IsDir() {
			// Если файл является директорией, записать её в выходящий канал.
			fmt.Printf("Pushing %s directory\n", fileInfo.Name())
			dirs <- path
		} else {
			// Если файл не является директорией, отобразить информацию в терминал.
			fmt.Printf("File %s, size: %dMB, last modified: %s\n",
				fileInfo.Name(), fileInfo.Size()/(1024*1024),
				fileInfo.ModTime().Format("15:04:05"))
		}
	}
}

func main() {
	// Создать каналы для файлов и директорий.
	filesChannel := make(chan string)
	dirsChannel := make(chan string)

	// Создать обрабатывающие горутины для файлов и директорий.
	go handleFiles(filesChannel, dirsChannel)
	go handleDirectories(dirsChannel, filesChannel)

	// Передать название директории из аргументов
	// командной строки в канал директорий.
	dirsChannel <- os.Args[1]
	// Заснуть на 60 секунд.
	time.Sleep(60 * time.Second)
}
