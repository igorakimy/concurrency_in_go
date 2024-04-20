package main

import (
	"fmt"
	"sync"
	"time"
)

type ReadWriteMutex struct {
	readersCounter int
	readersLock    sync.Mutex
	globalLock     sync.Mutex
}

func (rw *ReadWriteMutex) ReadLock() {
	rw.readersLock.Lock()
	rw.readersCounter++
	if rw.readersCounter == 1 {
		rw.globalLock.Lock()
	}
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	rw.globalLock.Lock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
	rw.readersLock.Lock()
	rw.readersCounter--
	if rw.readersCounter == 0 {
		rw.globalLock.Unlock()
	}
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.globalLock.Unlock()
}

func main() {
	// Использовать кастомный reader/writer мьютекс.
	rwMutex := ReadWriteMutex{}

	// Запустить две горутины.
	for i := 0; i < 2; i++ {
		go func() {
			// Повторять бесконечно.
			for {
				rwMutex.ReadLock()
				// Спать одну секунду, пока удерживыается блокировка на чтение.
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.ReadUnlock()
			}
		}()
	}

	time.Sleep(1 * time.Second)
	// Попытаться завладеть блокировкой на запись из main() горутины.
	rwMutex.WriteLock()
	// После того, как блокировка на запись будет получена,
	// вывести сообщение и завершить программу.
	fmt.Println("Write finished")
}
