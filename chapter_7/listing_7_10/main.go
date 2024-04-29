package main

import (
	"fmt"
	"time"
)

func receiver(messages <-chan int) {
	for {
		// Прочитать сообщение и флаг открытого канала,
		// который будет установлен в значение false, если канал закрыт.
		msg, more := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg, more)
		time.Sleep(1 * time.Second)
		// Когда сообщений больше нет, получение из канала прекращается.
		if !more {
			return
		}
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
