package main

import "fmt"

func main() {
	// Создать новый канал с типом string.
	msgChannel := make(chan string)

	// Запустить новую горутину с сылкой на канал.
	go receiver(msgChannel)

	// Отправить 3 строковых сообщения через канал.
	fmt.Println("Sending HELLO...")
	msgChannel <- "HELLO"
	fmt.Println("Sending THERE...")
	msgChannel <- "THERE"
	fmt.Println("Sending STOP...")
	msgChannel <- "STOP"
}

func receiver(messages chan string) {
	msg := ""
	// Получать сообщения до тех пор, пока
	// не будет получено сообщение "STOP".
	for msg != "STOP" {
		// Прочитать следующее сообщение из канала.
		msg = <-messages
		// Вывести полученное сообщение в консоли.
		fmt.Println("Received:", msg)
	}
}
