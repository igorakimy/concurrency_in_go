package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func generateUrls(quit <-chan int) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls)
		for i := 100; i <= 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			case urls <- url:
			case <-quit:
				return
			}
		}
	}()
	return urls
}

func downloadPages(quit <-chan int, urls <-chan string) <-chan string {
	pages := make(chan string)
	go func() {
		defer close(pages)
		moreData, url := true, ""
		for moreData {
			select {
			case url, moreData = <-urls:
				if moreData {
					resp, _ := http.Get(url)
					if resp.StatusCode != 200 {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					_ = resp.Body.Close()
				}
			case <-quit:
				return
			}
		}
	}()
	return pages
}

func extractWords(quit <-chan int, pages <-chan string) <-chan string {
	// Создать выходной канал, который будет содержать извлеченные слова.
	words := make(chan string)
	go func() {
		defer close(words)
		// Создать регулярное выражение для извлечения слов.
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		moreData, pg := true, ""
		for moreData {
			select {
			// Обновлять переменные с новым сообщением и флагом,
			// который отображает наличие данных в канале.
			case pg, moreData = <-pages:
				if moreData {
					for _, word := range wordRegex.FindAllString(pg, -1) {
						// Когда новый текст со страницы получен, извлечь
						// все слова при помощи регулярного выражения и
						// отправить их в выходной канал.
						words <- strings.ToLower(word)
					}
				}
			// Когда сообщение достигает канала quit, горутина завершает работу.
			case <-quit:
				return
			}
		}
	}()
	// Вернуть выходной канал.
	return words
}

func main() {
	quit := make(chan int)
	defer close(quit)
	results := extractWords(quit, downloadPages(quit, generateUrls(quit)))
	for result := range results {
		fmt.Println(result)
	}
}
