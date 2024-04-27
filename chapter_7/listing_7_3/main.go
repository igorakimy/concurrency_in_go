package main

import (
	"fmt"
	"time"
)

// Поскольку горутина receiver() завершается через 5 секунд,
// никакая другая горутина не доступна для получения сообщения из канала.
// Среда выполнения Go понимает это и выбрасывает фатальную ошибку,
// говорящую о том, что возникло состояние deadlock.

func main() {
	msgChannel := make(chan string)
	go receiver(msgChannel)
	fmt.Println("Sending HELLO...")
	msgChannel <- "HELLO"
	fmt.Println("Sending THERE...")
	msgChannel <- "THERE"
	fmt.Println("Sending STOP...")
	msgChannel <- "STOP"
}

func receiver(messages chan string) {
	// Подождать 5 секунд, вместо чтения из канала.
	time.Sleep(5 * time.Second)
	fmt.Println("Receiver slept for 5 seconds")
}
