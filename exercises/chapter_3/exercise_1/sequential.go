package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func calculateWords(url string, frequency map[string]int) {
	resp, _ := http.Get(url)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		fmt.Printf("Server error response with status: %s", resp.Status)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`[a-z']+`)
	words := re.FindAllString(strings.ToLower(string(body)), -1)

	for _, word := range words {
		if _, ok := frequency[word]; !ok {
			frequency[word] = 1
		} else {
			frequency[word] += 1
		}
	}
}

func main() {
	frequency := make(map[string]int)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		calculateWords(url, frequency)
	}

	for w, f := range frequency {
		fmt.Printf("%s -> %d\n", w, f)
	}
}
