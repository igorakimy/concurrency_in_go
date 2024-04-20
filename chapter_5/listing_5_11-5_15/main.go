package main

import (
	"fmt"
	"sync"
	"time"
)

type ReadWriteMutex struct {
	// Сохраняет количество читателей, на которых в данный
	// момент установлена блокировка чтения.
	readersCounter int

	// Сохраняет количество писателей, которые в данный момент ожидают.
	writersWaiting int

	// Индикатор, который показывает, удерживает ли
	// писатель блокировку на запись.
	writerActive bool

	cond *sync.Cond
}

// NewReadWriteMutex инициализирует новый ReadWriteMutex
// с новой "условной переменной" и связанным с ней мьютексом.
func NewReadWriteMutex() *ReadWriteMutex {
	return &ReadWriteMutex{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *ReadWriteMutex) ReadLock() {
	// Приобретает мьютекс.
	rw.cond.L.Lock()

	// "Условная переменная" ожидает, пока писатели
	// находятся в режиме ожидания или активны.
	for rw.writersWaiting > 0 || rw.writerActive {
		rw.cond.Wait()
	}
	// Увеличить счетчик читателей.
	rw.readersCounter++

	// Отпустить мьютекс.
	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	// Приобретает мьютекс.
	rw.cond.L.Lock()

	// Увеличивает счетчик ожидающих писателей.
	rw.writersWaiting++
	// "Условная переменная" ожидает до тех пор,
	// пока есть читатели или один активный писатель.
	for rw.readersCounter > 0 || rw.writerActive {
		rw.cond.Wait()
	}
	// Как только ожидание закончится, уменьшить на 1
	// счетчик ожидающих писателей.
	rw.writersWaiting--
	// Как только ожидание закончится, отметить флаг writerActive.
	rw.writerActive = true

	// Освобождает мьютекс.
	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
	// Приобретает мьютекс.
	rw.cond.L.Lock()

	// Уменьшает счетчик читателей на 1.
	rw.readersCounter--
	if rw.readersCounter == 0 {
		// Отправить сигнал, если горутина является
		// последним оставшимся читателем.
		rw.cond.Broadcast()
	}

	// Освобождает мьютекс.
	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	// Приобретает мьютекс.
	rw.cond.L.Lock()

	// Убирает флаг writerActive активного писателя.
	rw.writerActive = false
	// Отправляет сигнал.
	rw.cond.Broadcast()

	// Освобождает мьютекс.
	rw.cond.L.Unlock()
}

func main() {
	rwMutex := NewReadWriteMutex()

	for i := 0; i < 2; i++ {
		go func() {
			for {
				rwMutex.ReadLock()
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.ReadUnlock()
			}
		}()
	}

	time.Sleep(1 * time.Second)
	rwMutex.WriteLock()

	fmt.Println("Write finished")
}
