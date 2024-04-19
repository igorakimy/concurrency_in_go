package main

import (
	"fmt"
	"sync"
)

// Неправильный способ отправки сигналов.

func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	// Горутина сигнализирует о завершении работы.
	cond.Signal()
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()

	// Повторить 50000 раз.
	for i := 0; i < 50000; i++ {
		// Запустить горутину, симулируя выполнение какой-то работы.
		go doWork(cond)
		fmt.Println("Waiting for child goroutine")
		// Подождать завершающий сигнал от горутины.
		cond.Wait()
		fmt.Println("Child goroutine finished")
	}

	cond.L.Unlock()
}
