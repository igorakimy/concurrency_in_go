package main

import (
	"fmt"
	"time"
)

func sendMsgAfter(seconds time.Duration) <-chan string {
	messages := make(chan string)
	go func() {
		time.Sleep(seconds)
		messages <- "Hello"
	}()
	return messages
}

func main() {
	// Отправляет сообщение в канал после 3 секунд.
	messages := sendMsgAfter(3 * time.Second)

	for {
		select {
		// Прочитать сообщение из канала, если оно есть.
		case msg := <-messages:
			fmt.Println("Message received:", msg)
			// Когда сообщение получено прервать выполнение.
			return
		// Если сообщение не поступало, выполняется случай по умолчанию.
		default:
			fmt.Println("No messages waiting")
			time.Sleep(1 * time.Second)
		}
	}
}
