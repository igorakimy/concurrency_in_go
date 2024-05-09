package main

import (
	"fmt"
	"io"
	"net/http"
)

const downloaders = 20

func main() {
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)
	// Создать срез, чтобы сохранить выходные каналы из загружающих горутин.
	pages := make([]<-chan string, downloaders)
	// Создать 20 горутин, чтобы загрузить веб-страницы
	// и сохранить выходные каналы
	for i := 0; i < downloaders; i++ {
		pages[i] = downloadPages(quit, urls)
	}
}

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
						panic("Server's error:" + resp.Status)
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
