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

func main() {
	// Создать новую группу ожидания.
	wg := sync.WaitGroup{}
	// Добавить дельту в размере 31 (по одной на каждую
	// веб-страницу, загружаемую конкурентно).
	wg.Add(31)
	mutex := sync.Mutex{}
	frequency := make([]int, 26)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		// Создать горутину при помощи анонимной функции.
		go func() {
			countLetters(url, frequency, &mutex)
			// Вызвать Done() после завершения подсчета букв.
			wg.Done()
		}()
	}
	// Подождать, пока все горутины завершат работу.
	wg.Wait()

	mutex.Lock()
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}

func countLetters(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
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
