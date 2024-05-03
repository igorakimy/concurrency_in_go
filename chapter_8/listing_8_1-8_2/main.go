package main

import (
	"fmt"
	"time"
)

func writeEvery(msg string, seconds time.Duration) <-chan string {
	// Создать канал строкового типа.
	messages := make(chan string)

	// Создать новую, анонимную горутину.
	go func() {
		for {
			// Заснуть на указанный период.
			time.Sleep(seconds)
			// Отправить указанное сообщение в канал.
			messages <- msg
		}
	}()

	// Вернуть недавно созданные канал сообщений.
	return messages
}

func main() {
	// Создать горутину, которорая отправляет
	// каждую секунду сообщение в канал A.
	messagesFromA := writeEvery("Tick", 1*time.Second)
	// Создать горутину, которорая отправляет
	// каждые 3 секунды сообщение в канал B.
	messagesFromB := writeEvery("Tock", 3*time.Second)

	// Бесконечный цикл.
	for {
		select {
		// Выводит сообщение из канала A, если оно доступно.
		case msg1 := <-messagesFromA:
			fmt.Println(msg1)
		// Выводит сообщение из канала B, если оно доступно.
		case msg2 := <-messagesFromB:
			fmt.Println(msg2)
		}
	}
}
