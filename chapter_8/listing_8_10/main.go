package main

import "fmt"

func main() {
	// Создать нулевой канал.
	var ch chan string = nil

	// Заблокировать выполнение при попытке
	// оптравить сообщение в нулевой канал.
	ch <- "message"

	fmt.Println("This is never printed")
}
