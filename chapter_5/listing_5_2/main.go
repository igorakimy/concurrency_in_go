package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func stingy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		mutex.Lock()
		*money += 10
		mutex.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 200000; i++ {
		mutex.Lock()
		// Продолжает пытаться, если денег на счету недостаточно.
		for *money < 50 {
			// Разблокирует мьютекс, позволяя другой горутине
			// получить доступ к переменной money.
			mutex.Unlock()
			// Заснуть на короткое время.
			time.Sleep(10 * time.Millisecond)
			// Заблокировать мьютекс стнова, чтобы обеспечить
			// доступ к последнему значению переменной money.
			mutex.Lock()
		}
		*money -= 50
		if *money < 0 {
			fmt.Println("Money is negative!")
			os.Exit(1)
		}
		mutex.Unlock()
	}
	fmt.Println("Spendy Done")
}

func main() {
	mutex := sync.Mutex{}
	money := 100

	go stingy(&money, &mutex)
	go spendy(&money, &mutex)

	time.Sleep(2 * time.Second)

	fmt.Println("Money in bank account:", money)
}
