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
	mutex := sync.Mutex{}
	frequency := make([]int, 26)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go CountLetters(url, frequency, &mutex)
	}
	time.Sleep(3 * time.Second)

	mutex.Lock()
	for i, c := range AllLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}

func CountLetters(url string, frequency []int, mutex *sync.Mutex) {
	// Оставить медленную часть функции(которая загружает страничку)
	// конкурентной, чтобы не блокировать выполнение других горутин.
	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	// Заблокировать только быстро выполняющуюся секцию функции.
	mutex.Lock()
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(AllLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()

	fmt.Println("Completed: ", url, time.Now().Format("15:04:05"))
}
