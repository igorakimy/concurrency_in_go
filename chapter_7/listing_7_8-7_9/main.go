package main

import (
	"fmt"
	"time"
)

// Объявить канал, который будет только принимать сообщения.
func receiver(messages <-chan int) {
	for {
		// Прочитать сообщение из канала.
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		// Подождать 1 секунду.
		time.Sleep(1 * time.Second)
	}
}

func main() {
	msgChannel := make(chan int)

	go receiver(msgChannel)

	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgChannel <- i
		time.Sleep(1 * time.Second)
	}
	close(msgChannel)

	time.Sleep(3 * time.Second)
}
