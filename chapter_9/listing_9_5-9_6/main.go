package main

import (
	"fmt"
	"io"
	"net/http"
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
	// Создать выходной канал, который будет содержать
	// загруженные веб-страницы.
	pages := make(chan string)
	go func() {
		// Закрыть выходной канал, когда горутина завершит работу.
		defer close(pages)
		moreData, url := true, ""
		// Продолжать выбирать, если входящий канал
		// со ссылками всё ещё открыт и в нем есть данные.
		for moreData {
			select {
			// Обновляет переменные новым сообщением и флагом,
			// чтобы показать, есть ли ещё данные в канале.
			case url, moreData = <-urls:
				// Когда поступает новое сообщение с URL-ом,
				// загрузить страницу и отправить текст со страницы
				// в канал pages.
				if moreData {
					resp, _ := http.Get(url)
					if resp.StatusCode != 200 {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					_ = resp.Body.Close()
				}
			// Когда сообщение поступит в канал quit,
			// завершить работу горутины.
			case <-quit:
				return
			}
		}
	}()
	// Вернуть выходной канал.
	return pages
}

func main() {
	quit := make(chan int)
	defer close(quit)
	// Добавить новую горутину, которая загружает страницы в существующий пайплайн.
	results := downloadPages(quit, generateUrls(quit))
	for result := range results {
		fmt.Println(result)
	}
}
