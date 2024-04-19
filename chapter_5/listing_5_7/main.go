package main

import (
	"fmt"
	"sync"
)

func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	// Заблокировать мьютекс до отправки сигнала.
	cond.L.Lock()
	// Отправить сигнал "условной переменной".
	cond.Signal()
	// Разблокировать мьютекс после отправки сигнала.
	cond.L.Unlock()
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})

	cond.L.Lock()

	for i := 0; i < 50000; i++ {
		go doWork(cond)

		fmt.Println("Waiting for child goroutine")
		cond.Wait()
		fmt.Println("Child goroutine finished")
	}

	cond.L.Unlock()
}