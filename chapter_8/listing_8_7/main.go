package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// sendMsgAfter отправляет сообщение "Hello" в возвращаемый канал
// после указанного количества секунд.
func sendMsgAfter(seconds time.Duration) <-chan string {
	messages := make(chan string)

	go func() {
		time.Sleep(seconds)
		messages <- "Hello"
	}()

	return messages
}

func main() {
	// Прочитать значение таймаута из аргументов командной строки.
	t, _ := strconv.Atoi(os.Args[1])

	// Запустить горутину, которая отправит сообщение в возвращаемый
	// канал по прошествии 3-х секунд.
	messages := sendMsgAfter(3 * time.Second)
	timeoutDuration := time.Duration(t) * time.Second

	fmt.Printf("Waiting for message for %d seconds...\n", t)

	select {
	// Читать сообщение из из канала messages, если таковое есть.
	case msg := <-messages:
		fmt.Println("Message received:", msg)
	// Создать канал и таймер, который получает сообщение
	// после указанного значения времени.
	case tNow := <-time.After(timeoutDuration):
		fmt.Println("Timed out. Waited until:", tNow.Format("15:04:05"))
	}
}
