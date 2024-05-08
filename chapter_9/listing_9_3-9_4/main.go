package main

import "fmt"

// generateUrls принимает канал quit и возвращает выходной канал.
func generateUrls(quit <-chan int) <-chan string {
	// Создать выходной канал.
	urls := make(chan string)
	go func() {
		// Как только работа будет завершена, закрыть выходной канал.
		defer close(urls)
		for i := 100; i <= 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			// Записать 50 урлов в выходной канал.
			case urls <- url:
			case <-quit:
				return
			}
		}
	}()
	// Вернуть выходной канал.
	return urls
}

func main() {
	// Создать завершающий канал quit.
	quit := make(chan int)
	defer close(quit)
	// Вызвать функцию, чтобы запустить горутину,
	// которая возвращает ссылки в результирующем канале.
	results := generateUrls(quit)
	// Прочитать все сообщения из результирующего канала.
	for result := range results {
		// Вывести результат.
		fmt.Println(result)
	}
}
