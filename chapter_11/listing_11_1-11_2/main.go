package main

import (
	"fmt"
	"sync"
	"time"
)

func red(lock1, lock2 *sync.Mutex) {
	for {
		// Завладеть и удерживать обе блокировки.
		fmt.Println("Red: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Red: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Red: Both locks Acquired")
		// Освободить обе блокировки.
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Red Locks Released")
	}
}

func blue(lock1, lock2 *sync.Mutex) {
	for {
		// Завладеть и удерживать обе блокировки.
		fmt.Println("Blue: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Blue: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Blue: Both locks Acquired")
		// Освободить обе блокировки.
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Blue: Locks Released")
	}
}

func main() {
	lockA := sync.Mutex{}
	lockB := sync.Mutex{}

	// Запуск горутины red.
	go red(&lockA, &lockB)
	// Запуск горутины blue.
	go blue(&lockA, &lockB)

	// Позволяет горутинам red и blue завершиться, ожидая 20 секунд.
	time.Sleep(20 * time.Second)
	fmt.Println("Done")
}
