package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const AllLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	// Создается новый мьютекс.
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		// Передать ссылку на мьютекс всем горутинам.
		go CountLetters(url, frequency, &mutex)
	}
	// Подождать 30 секунд.
	time.Sleep(30 * time.Second)

	// Защитить чтение общей переменной с помощью мьютекса.
	mutex.Lock()
	for i, c := range AllLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}

// Неправильный случай блокировки/разблокировки мьютекса.

func CountLetters(url string, frequency []int, mutex *sync.Mutex) {
	// Заблокировать мьютекс для всего участка кода,
	// что сделает его последовательным.
	mutex.Lock()

	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(AllLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}

	fmt.Println("Completed: ", url, time.Now().Format("15:04:05"))
	// Разблокировать мьютекс.
	mutex.Unlock()
}
