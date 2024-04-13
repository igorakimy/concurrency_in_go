package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency []int) {
	// Загрузить веб-страничку с указанного URL.
	resp, _ := http.Get(url)
	// Закрыть ответ в конце функции.
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	// Итерировать по каждому загруженному символу.
	for _, b := range body {
		c := strings.ToLower(string(b))
		// Найти индекс символа в алфавите.
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			// Если символ является частью алфавита,
			// увеличить количество на 1.
			frequency[cIndex] += 1
		}
	}

	fmt.Println("Completed: ", url)
}

func main() {
	// Инициализация слайса для таблицы частот встречающихся букв.
	var frequency = make([]int, 26)

	// Итерировать от 1000 до 1030 ID документа,
	// чтобы загрузить 31 документ.
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		// Последовательный вызов фукнции countLetters().
		countLetters(url, frequency)
	}

	// Вывести каждую букву и частоту, с которой она встречается.
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
}
