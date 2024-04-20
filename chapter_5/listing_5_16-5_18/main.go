package main

import (
	"fmt"
	"sync"
)

type Semaphore struct {
	// Разрешения, оставшиеся на семафоре.
	permits int

	// "Условная переменная", которая используется для ожидания,
	// когда не достаточно разрешений.
	cond *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		// Инициализировать количество разрешений на новом семафоре.
		permits: n,
		// Инициализировать новую "условную переменную", и связанным
		// с мьютексом новый семафор.
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Semaphore) Acquire() {
	// Завладеть мьютексом, чтобы защитить переменную permits.
	s.cond.L.Lock()
	for s.permits <= 0 {
		// Подождать, пока появится доступное разрешение.
		s.cond.Wait()
	}
	// Уменьшить количество доступных разрешения на 1.
	s.permits--
	// Отпустить мьютекс.
	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	// Завладеть мьютексом, чтобы защитить переменную permits.
	s.cond.L.Lock()
	// Увеличить количество доступных разрешений на 1.
	s.permits++
	// Отправить сигнал "условной переменной" о том,
	// что доступно еще одно разрешение.
	s.cond.Signal()
	// Отпустить мьютекс.
	s.cond.L.Unlock()
}

func main() {
	// Создать новый семафор.
	semaphore := NewSemaphore(0)

	// Повторить 50000 раз.
	for i := 0; i < 50000; i++ {
		// Запустить горутину, передав ссылку на семафор.
		go doWork(semaphore)
		fmt.Println("Waiting for child goroutine")
		// Подождать доступное разрешение на семафоре,
		// указывающего на выполнение задачи.
		semaphore.Acquire()
		fmt.Println("Child goroutine finished")
	}
}

func doWork(semaphore *Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	// Когда горутина завершает работу, она выдает разрешение
	// на уведомление основной main() горутины.
	semaphore.Release()
}
