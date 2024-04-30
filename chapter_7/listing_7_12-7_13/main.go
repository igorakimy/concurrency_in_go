package main

import "fmt"

// findFactors находит все делители входящего числа.
func findFactors(number int) []int {
	result := make([]int, 0)

	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}

	return result
}

func main() {
	// Создать новый канал с типом целочисленного среза.
	resultCh := make(chan []int)

	go func() {
		// Вызывать функцию в анонимной горутине и поместить
		// результат в канал.
		resultCh <- findFactors(3419110721)
	}()

	fmt.Println(findFactors(4033836233))
	// Собрать результат из канала.
	fmt.Println(<-resultCh)
}
