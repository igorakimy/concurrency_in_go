package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		panic("Server response error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	mutex.Lock()
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()

	fmt.Println("Completed: ", url, time.Now().Format("15:04:05"))
}

func main() {
	mutex := sync.Mutex{}
	frequency := make([]int, 26)

	for i := 2000; i <= 2200; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency, &mutex)
	}

	for i := 0; i < 100; i++ {
		// Заснуть на 100 миллисекунд.
		time.Sleep(100 * time.Millisecond)

		// Попытаться завладеть блокировкой мьютекса.
		if mutex.TryLock() {
			// Если блокировка мьютекса доступна, то выводится
			// буква и частота, с которой она встречается.
			for idx, c := range allLetters {
				fmt.Printf("%c-%d ", c, frequency[idx])
			}
			// После этого мьютекс освобождается.
			mutex.Unlock()
		} else {
			// Если же мьютекс недоступен, то выводится сообщение
			// о том, что мьютекс уже используется, а после попытается
			// снова им завладеть.
			fmt.Println("Mutex already being used")
		}
	}
}
