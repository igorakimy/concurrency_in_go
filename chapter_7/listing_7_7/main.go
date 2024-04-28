package main

import (
	"fmt"
	"time"
)

// Объявить канал, который только принимает сообщения.
func receiver(messages <-chan int) {
	for {
		// Получать сообщения с канала.
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received", msg)
	}
}

// Объявить канал, который только отправляет сообщения.
func sender(messages chan<- int) {
	for i := 1; ; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending", i)
		// Отправлять сообщение в канал каждую секунду.
		messages <- i
		time.Sleep(1 * time.Second)
	}
}

func main() {
	msgChannel := make(chan int)
	go receiver(msgChannel)
	go sender(msgChannel)
	time.Sleep(5 * time.Second)
}
