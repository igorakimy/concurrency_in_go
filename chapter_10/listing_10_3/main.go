package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func FHash(filename string) []byte {
	file, _ := os.Open(filename)
	defer func() { _ = file.Close() }()

	sha := sha256.New()
	_, _ = io.Copy(sha, file)

	return sha.Sum(nil)
}

func main() {
	dir := os.Args[1]
	// Получить список файлов из указанной директории.
	files, _ := os.ReadDir(dir)
	// Создать новый, пустой хэш-контейнер для директории.
	sha := sha256.New()

	for _, file := range files {
		if !file.IsDir() {
			fPath := filepath.Join(dir, file.Name())
			// Вычислить хэш-сумму для каждого файла в директории.
			hashOnFile := FHash(fPath)
			// Конкатенировать вычесленную хэш-сумму к хеш-сумме директории.
			sha.Write(hashOnFile)
		}
	}

	// Вывести финальную хэш-сумму.
	fmt.Printf("%s - %x\n", dir, sha.Sum(nil))
}
