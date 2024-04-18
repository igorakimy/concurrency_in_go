package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func stingy(money *int, cond *sync.Cond) {
	for i := 0; i < 1_000_000; i++ {
		// Использовать мьютекс для "условной переменной".
		cond.L.Lock()
		*money += 10
		// Сигнализировать "условной переменной" каждый раз,
		// когда добавляется значение к общей переменной money.
		cond.Signal()
		// Использовать мьютекс для "условной переменной".
		cond.L.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, cond *sync.Cond) {
	for i := 0; i < 200000; i++ {
		// Использовать мьютекс для "условной переменной".
		cond.L.Lock()
		// Ждать, пока не будет достаточно средств, чтобы
		// освободить мьютекс и приостановить выполнение.
		for *money < 50 {
			cond.Wait()
		}
		// Возвращаясь из функции Wait(), повторно завладевает
		// мьютексом и уменьшает значение money, как только их
		// наберется достаточно.
		*money -= 50
		if *money < 0 {
			fmt.Println("Money is negative!")
			os.Exit(1)
		}
		cond.L.Unlock()
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	// Создать новый мьютекс.
	mutex := sync.Mutex{}
	// Создать новую "условную" переменную, используя мьютекс.
	cond := sync.NewCond(&mutex)

	// Передать "условную" переменную в обе горутины.
	go stingy(&money, cond)
	go spendy(&money, cond)

	time.Sleep(2 * time.Second)

	mutex.Lock()
	fmt.Println("Money is bank account:", money)
	mutex.Unlock()
}
