package main

import (
	"fmt"
	"time"
)

// Функции принимают указатель на переменную, которая
// является суммой на банковском счете.

func stingy(money *int) {
	for i := 0; i < 1_000_000; i++ {
		// Функция stingy() добавляет 10 долларов.
		*money += 10
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int) {
	for i := 0; i < 1_000_000; i++ {
		// Функция spendy() отнимает 10 долларов.
		*money -= 10
	}
	fmt.Println("Spendy Done")
}

func main() {
	// Инициализация денежной суммы в 100 долларов
	// на банковском счете.
	money := 100

	// Запустить две горутины, передав каждой ссылку
	// на переменную значения суммы.
	go stingy(&money)
	go spendy(&money)

	// Подождать 2 секунды, пока горутины завершат работу.
	time.Sleep(2 * time.Second)
	fmt.Println("Money in bank account: ", money)
}
