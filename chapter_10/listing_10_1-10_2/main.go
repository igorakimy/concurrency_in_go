package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func FHash(filepath string) []byte {
	// Открыть файл.
	file, _ := os.Open(filepath)
	defer func() { _ = file.Close() }()

	// Вычислить хэш-код, используя библиотеку crypto/sha256.
	sha := sha256.New()
	_, _ = io.Copy(sha, file)

	// Вернуть результат хэширования.
	return sha.Sum(nil)
}

func main() {
	dir := os.Args[1]
	// Получить список файлов из указанной директории.
	files, _ := os.ReadDir(dir)
	wg := sync.WaitGroup{}

	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			// Запустить горутину, котрая будет вычислять хэш-сумму
			// для файла на каждой итерации.
			go func(filename string) {
				fPath := filepath.Join(dir, filename)
				// Вычислить и вывести хэш-сумму.
				hash := FHash(fPath)
				fmt.Printf("%s - %x\n", filename, hash)
				wg.Done()
			}(file.Name())
		}
	}

	// Подождать, пока все задачи, вычисляющие хэш, будут завершены.
	wg.Wait()
}
