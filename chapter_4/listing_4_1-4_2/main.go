package main

import (
	"fmt"
	"sync"
	"time"
)

// Функции принимают указатель на общую структуру мьютекса.

func stingy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1_000_000; i++ {
		// Блокировка мьютекса перед входом с "критическую секцию".
		mutex.Lock()

		*money += 10

		// Снятие блокировки мьютекса после выхода из "критической секции".
		mutex.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1_000_000; i++ {
		// Блокировка мьютекса перед входом с "критическую секцию".
		mutex.Lock()

		*money -= 10

		// Снятие блокировки мьютекса после выхода из "критической секции".
		mutex.Unlock()
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	// Создать новый мьютекс. Когда мьютекс создается,
	// он всегда находится в разблокированном состоянии.
	mutex := sync.Mutex{}

	// Передать ссылку на новый мьютекс всем горутинам.
	go stingy(&money, &mutex)
	go spendy(&money, &mutex)

	time.Sleep(2 * time.Second)

	// Защитить чтение общей переменной с помощью мьютекса.
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println("Money in bank account: ", money)
}
