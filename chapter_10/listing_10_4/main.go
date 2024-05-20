package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	sha := sha256.New()
	var prev, next chan int

	for _, file := range files {
		if !file.IsDir() {
			// Создать канал next, который будет использован
			// горутиной чтобы сообщить, что она готова.
			next = make(chan int)
			go func(filename string, prev, next chan int) {
				fPath := filepath.Join(dir, filename)
				// Вычислить хэш файла.
				hashOnFile := FHash(fPath)
				// Если эта горутина не на первой итерации,
				// подождать пока предыдущая итерация отправит сигнал.
				if prev != nil {
					<-prev
				}
				// Вычислить частичный хэш директории.
				sha.Write(hashOnFile)
				// Сигнализировать, что следующая итерация завершилась.
				next <- 0
			}(file.Name(), prev, next)
			// Сделать следующий канал предыдущим; Следующая горутина
			// будет ждать сигнала от текущей итерации.
			prev = next
		}
	}
	// Подождать последней итерации, чтобы завершить завершить
	// работу до вывода результата.
	<-next
	fmt.Printf("%x\n", sha.Sum(nil))
}

func FHash(filename string) []byte {
	file, _ := os.Open(filename)
	defer func() { _ = file.Close() }()

	sha := sha256.New()
	_, _ = io.Copy(sha, file)

	return sha.Sum(nil)
}
