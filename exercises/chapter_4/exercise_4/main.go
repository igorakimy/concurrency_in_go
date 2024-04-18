package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func calculateWords(url string, frequency map[string]int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		fmt.Printf("Server error response with status: %s", resp.Status)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`[a-z]+`)
	words := re.FindAllString(strings.ToLower(string(body)), -1)

	mutex.Lock()
	for _, word := range words {
		if _, ok := frequency[word]; !ok {
			frequency[word] = 1
		} else {
			frequency[word] += 1
		}
	}
	mutex.Unlock()
}

func main() {
	mutex := sync.Mutex{}
	frequency := make(map[string]int)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go calculateWords(url, frequency, &mutex)
	}
	time.Sleep(2 * time.Second)

	mutex.Lock()
	for w, f := range frequency {
		fmt.Printf("%s -> %d\n", w, f)
	}
	mutex.Unlock()
}
