package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const downloaders = 20

func main() {
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)
	pages := make([]<-chan string, downloaders)
	for i := 0; i < downloaders; i++ {
		pages[i] = downloadPages(quit, urls)
	}
	// Соединяет все каналы с содержимым страниц в один канал,
	// используя fan-in шаблон.
	results := extractWords(quit, FanIn(quit, pages...))
	for result := range results {
		fmt.Println(result)
	}
}

func FanIn[K any](quit <-chan int, allChannels ...<-chan K) chan K {
	// Создать группу ожидания, задав размер,
	// равный количеству переданных каналов.
	wg := sync.WaitGroup{}
	wg.Add(len(allChannels))

	// Создать выходной канал.
	output := make(chan K)

	for _, c := range allChannels {
		// Создавать горутину для каждого канала.
		go func(channel <-chan K) {
			// Как только горутина завершит работу, отметить
			// её в группе ожидания, как завершенную.
			defer wg.Done()
			for i := range channel {
				select {
				// Пересылает каждое полученное сообщение в общий выходной канал.
				case output <- i:
				// Если закрыть канал quit, горутина прекратит работу.
				case <-quit:
					return
				}
			}
			// Передать один входящий канал в горутину.
		}(c)
	}

	go func() {
		// Подождать пока все горутины завершат работу
		// и тогда закрыть выходной канал output.
		wg.Wait()
		close(output)
	}()
	// Вернуть выходной канал.
	return output
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

func extractWords(quit <-chan int, pages <-chan string) <-chan string {
	words := make(chan string)
	go func() {
		defer close(words)
		moreData, pg := true, ""
		for moreData {
			select {
			case pg, moreData = <-pages:
				if moreData {
					pattern := regexp.MustCompile(`[a-zA-Z]+`)
					for _, word := range pattern.FindAllString(pg, -1) {
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return words
}
