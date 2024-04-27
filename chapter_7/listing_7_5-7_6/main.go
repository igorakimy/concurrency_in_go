package main

import (
	"fmt"
	"sync"
	"time"
)

func receiver(messages chan int, wg *sync.WaitGroup) {
	msg := 0

	// Продолжать читать сообщения из канала,
	// пока не будет получено значение -1.
	for msg != -1 {
		// Подождать одну секунду.
		time.Sleep(1 * time.Second)
		// Прочитать следующее сообщение из канала.
		msg = <-messages
		fmt.Println("Received:", msg)
	}

	// Вызвать Done() на группе ожидания
	// после прочтения всех сообщений.
	wg.Done()
}

func main() {
	// Создать новый канал с емкостью буфера - 3 сообщения.
	msgChannel := make(chan int, 3)
	// Создать группу ожидания с размером 1.
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Запустить горутину, передав буферизованный канал и групу ожидания.
	go receiver(msgChannel, &wg)

	for i := 1; i <= 6; i++ {
		// Получить количество сообщений в буферизованном канале.
		size := len(msgChannel)
		fmt.Printf(
			"%s Sending : %d. Buffer Size: %d\n",
			time.Now().Format("15:04:05"), i, size,
		)
		// Отправлять шесть целочисленных сообщений от 1 до 6.
		msgChannel <- i
	}
	// Отправить сообщении, содержащее -1.
	msgChannel <- -1
	// Подождать в группе ожидания, пока receiver() не завершит работу.
	wg.Wait()
}
