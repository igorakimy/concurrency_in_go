package main

import "fmt"

func printNumbers(numbers <-chan int, quit chan int) {
	go func() {
		// Принимать 10 чисел из канала numbers.
		for i := 0; i < 10; i++ {
			fmt.Println(<-numbers)
		}
		// Закрыть канал quit.
		close(quit)
	}()
}

func main() {
	// Создать каналы для чисел и выхода.
	numbers := make(chan int)
	quit := make(chan int)

	// Вызвать функцию printNumbers, передав каналы.
	printNumbers(numbers, quit)

	next := 0
	for i := 1; ; i++ {
		// Генерировать следующее триугольное число.
		next += i
		select {
		// Отправить число в канал numbers.
		case numbers <- next:
		// Когда канал quit будет разблокирован,
		// вывести сообщение и прекратить выполнение.
		case <-quit:
			fmt.Println("Quitting number generation")
			return
		}
	}
}
