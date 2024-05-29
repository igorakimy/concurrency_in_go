package main

import (
	"fmt"
	"sync"
	"time"
)

func lockBoth(lock1, lock2 *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10_000; i++ {
		lock1.Lock()
		lock2.Lock()
		lock1.Unlock()
		lock2.Unlock()
	}
}

func main() {
	lockA, lockB := sync.Mutex{}, sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go lockBoth(&lockA, &lockB, &wg)
	go lockBoth(&lockB, &lockA, &wg)

	// Создать горутину, которая ожидает
	// на группе ожидания до вывода сообщения.
	go func() {
		wg.Wait()
		fmt.Println("Done waiting on waitgroup")
	}()
	// Подождать 10 секунд.
	time.Sleep(10 * time.Second)
	// Вывести сообщение и тогда программа завершит работу.
	fmt.Println("Done")
}
