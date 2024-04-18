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
		// Получает по 10 монет.
		*money += 10
		mutex.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 200000; i++ {
		mutex.Lock()
		// Но тратит 50 монет.
		*money -= 50
		// Когда переменная money достигнет отрицательного значения,
		// вывести сообщение и завершить программу.
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

	time.Sleep(1 * time.Second)
}
