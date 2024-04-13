package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency []int) {
	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}

	fmt.Println("Completed: ", url)
}

func main() {
	var frequency = make([]int, 26)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		// Запустить горутину, которая вызывает функцию countLetters().
		go countLetters(url, frequency)
	}
	// Подожать, пока горутины завершат свою работу.
	time.Sleep(4 * time.Second)

	// Вывести букви и частоту, с которой она встречается.
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
}
