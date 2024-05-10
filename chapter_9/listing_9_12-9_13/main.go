package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const (
	downloaders = 20
)

func main() {
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)
	pages := make([]<-chan string, downloaders)
	for i := 0; i < downloaders; i++ {
		pages[i] = downloadPages(quit, urls)
	}
	// Подключить горутину longestWords() к конвейеру просто после extractWords().
	results := longestWords(quit, extractWords(quit, FanIn(quit, pages...)))
	// Вывести одно сообщение, содержащее самые длинные слова.
	fmt.Println("Longest Words:", <-results)

}

func longestWords(quit <-chan int, words <-chan string) <-chan string {
	longWords := make(chan string)
	go func() {
		defer close(longWords)
		// Создать карту, чтобы сохранить уникальные слова.
		uniqueWordsMap := make(map[string]bool)
		// Создать срез, чтобы сохранить список уникальных
		// слов, чтобы упростить их сортировку после.
		uniqueWords := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			case word, moreData = <-words:
				// Если канал не закрыт и слово является новым,
				// то добавить новое слово в карту и срез.
				if moreData && !uniqueWordsMap[word] {
					uniqueWordsMap[word] = true
					uniqueWords = append(uniqueWords, word)
				}
			case <-quit:
				return
			}
		}
		// Как только входящий канал будет закрыт, отсортировать срез
		// с уникальными словами по длине слов.
		sort.Slice(uniqueWords, func(i, j int) bool {
			return len(uniqueWords[i]) > len(uniqueWords[j])
		})
		// Как только входящий канал будет закрыт, отправить строку,
		// которая будет содержать 10 самых длинных слов, перечисленных
		// через запятую.
		longWords <- strings.Join(uniqueWords[:10], ", ")
	}()
	return longWords
}

func FanIn[K any](quit <-chan int, pages ...<-chan K) <-chan K {
	wg := sync.WaitGroup{}
	wg.Add(len(pages))
	output := make(chan K)

	for _, c := range pages {
		go func(channel <-chan K) {
			defer wg.Done()
			for i := range channel {
				select {
				case output <- i:
				case <-quit:
					return
				}
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func generateUrls(quit <-chan int) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls)
		for i := 100; i <= 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d", i)
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
					page, _ := io.ReadAll(resp.Body)
					pages <- string(page)
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
					matches := regexp.MustCompile(`[a-zA-Z]+`)
					for _, word := range matches.FindAllString(pg, -1) {
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
