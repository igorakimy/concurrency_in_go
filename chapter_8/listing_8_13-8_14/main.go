package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string) <-chan []int {
	// Создать выходной канал с типом целочисленного среза.
	result := make(chan []int)

	go func() {
		defer close(result)
		// Создать локальную переменную среза частоты встречаемых букв.
		frequency := make([]int, 26)
		resp, _ := http.Get(url)
		if resp.StatusCode != 200 {
			panic("Server returning error code: " + resp.Status)
		}
		body, _ := io.ReadAll(resp.Body)
		for _, b := range body {
			c := strings.ToLower(string(b))
			cIndex := strings.Index(allLetters, c)
			if cIndex > 0 {
				// Обновляет количество каждого символа,
				// в локальной переменной среза частоты.
				frequency[cIndex] += 1
			}
		}
		fmt.Println("Completed:", url)
		// Как только завершит работу, в результирующий канал будет
		// отправлен срез с частотой встречаемых букв.
		result <- frequency
	}()

	return result
}

func main() {
	// Создать срез, который будет содержать все выходные каналы.
	results := make([]<-chan []int, 0)
	// Создать срез, чтобы сохранить частоту каждой
	// буквы английского алфавита.
	totalFrequencies := make([]int, 26)

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		// Создать горутину для каждой веб-страницы и сохранить
		// выходной канал в результирующий срез results.
		results = append(results, countLetters(url))
	}

	// Итерировать по каждому выходному каналу.
	for _, c := range results {
		// Получать сообщение из каждого выходного канала,
		// содержащего частоты для каждой отдельной веб-страницы.
		frequencyResult := <-c
		// Добавлять частоту, полученную из канала
		// к общей итоговой частоте.
		for i := 0; i < 26; i++ {
			totalFrequencies[i] += frequencyResult[i]
		}
	}

	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, totalFrequencies[i])
	}
}
