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

	words := extractWords(quit, FanIn(quit, pages...))
	// Создать горутину, которая будет транслировать
	// содержимое канала со словами в два выходных канала.
	wordsMulti := Broadcast(quit, words, 2)
	// Создать горутину, которая будет искать
	// самые длинные слова во входном канале.
	longestResults := longestWords(quit, wordsMulti[0])
	// Создать горутину, чтобы найти наиболее
	// часто встречающиеся слова во воходном канале.
	frequentResults := frequentWords(quit, wordsMulti[1])

	// Прочитать результат из горутины longestWords() и вывести его.
	fmt.Println("Longest Words:", <-longestResults)
	// Прочитать результат из горутины frequentWords() и вывести его.
	fmt.Println("Frequent Words:", <-frequentResults)

}

func Broadcast[K any](quit <-chan int, input <-chan K, n int) []chan K {
	// Создать n выходных каналов с типом K.
	outputs := CreateAll[K](n)
	go func() {
		// Как только горутина завершит работу,
		// закрыть все выходные каналы.
		defer CloseAll(outputs...)
		var msg K
		moreData := true
		for moreData {
			select {
			// Читать следующее сообщение из входящего канала.
			case msg, moreData = <-input:
				// Пока канал не закрылся, записывать сообщение
				// в каждый выходной канал.
				if moreData {
					for _, output := range outputs {
						output <- msg
					}
				}
			case <-quit:
				return
			}
		}
	}()
	// Вернуть множество выходных каналов.
	return outputs
}

// CreateAll создает n каналов с типом K.
func CreateAll[K any](n int) []chan K {
	channels := make([]chan K, n)
	for i, _ := range channels {
		channels[i] = make(chan K)
	}
	return channels
}

// CloseAll закрывает все указанные каналы.
func CloseAll[K any](outputs ...chan K) {
	for _, output := range outputs {
		close(output)
	}
}

func frequentWords(quit <-chan int, words <-chan string) <-chan string {
	mostFrequentWords := make(chan string)
	go func() {
		defer close(mostFrequentWords)
		// Создать карту, чтобы сохранять частоту
		// появления каждого уникального слова.
		freqMap := make(map[string]int)
		// Создать срез, чтобы сохранить список уникальных слов.
		freqList := make([]string, 0)
		moreData, word := true, ""
		for moreData {
			select {
			// Получать следующее сообщение из входящего канала.
			case word, moreData = <-words:
				if moreData {
					// Если сообщение содержит новое слово,
					// тогда добавить его в срес уникальных слов.
					if freqMap[word] == 0 {
						freqList = append(freqList, word)
					}
					// Увеличить частоту слова.
					freqMap[word] += 1
				}
			case <-quit:
				return
			}
		}
		// Как только все входящие сообщения будут получены,
		// отсортировать список слов по частоте появления.
		sort.Slice(freqList, func(i, j int) bool {
			return freqMap[freqList[i]] > freqMap[freqList[j]]
		})
		// Записать 10 наиболее часто встречаемых слов в выходной канал.
		mostFrequentWords <- strings.Join(freqList[:10], ", ")
	}()
	return mostFrequentWords
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
					if resp.StatusCode != http.StatusOK {
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

func longestWords(quit <-chan int, words <-chan string) <-chan string {
	longWords := make(chan string)

	go func() {
		defer close(longWords)
		uniqueWordsMap := make(map[string]bool)
		uniqueWords := make([]string, 0)
		moreData, word := true, ""

		for moreData {
			select {
			case word, moreData = <-words:
				if moreData && !uniqueWordsMap[word] {
					uniqueWordsMap[word] = true
					uniqueWords = append(uniqueWords, word)
				}
			case <-quit:
				return
			}
		}

		sort.Slice(uniqueWords, func(i, j int) bool {
			return len(uniqueWords[i]) > len(uniqueWords[j])
		})
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
