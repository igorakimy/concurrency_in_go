package main

import (
	"fmt"
	"sync"
)

func lockBoth(lock1, lock2 *sync.Mutex, wg *sync.WaitGroup) {
	for i := 0; i < 10_000; i++ {
		// Блокировка и разблокировка обоих мьютексов.
		lock1.Lock()
		lock2.Lock()
		lock1.Unlock()
		lock2.Unlock()
	}
	// Отметить, что группа ожидания завершилась.
	wg.Done()
}

func main() {
	lockA, lockB := sync.Mutex{}, sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Запустить две горутины, которые блокируют
	// оба мьютекса в одно и то же время.
	go lockBoth(&lockA, &lockB, &wg)
	go lockBoth(&lockB, &lockA, &wg)

	// Подождать, пока горутины завершат свою работу.
	wg.Wait()

	fmt.Println("Done")
}
