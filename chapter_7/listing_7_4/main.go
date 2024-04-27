package main

import (
	"fmt"
	"time"
)

// Также возникнет состояние, называемое deadlock, т.к. горутина,
// которая должна отправлять сообщение, но не делает этого, засыпая
// на 5 секунд.

func main() {
	// Создать новый канал с типом string.
	msgChannel := make(chan string)

	go sender(msgChannel)

	fmt.Println("Reading message from channel...")
	// Попробовать прочитать сообщение из канала.
	msg := <-msgChannel
	fmt.Println("Received:", msg)
}

func sender(messages chan string) {
	// Заснуть на 5 секунд, вместо отправки сообщения.
	time.Sleep(5 * time.Second)
	fmt.Println("Sender slept for 5 seconds")
}
