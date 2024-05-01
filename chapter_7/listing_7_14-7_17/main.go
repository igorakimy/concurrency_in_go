package main

import (
	"container/list"
	"sync"
)

type Semaphore struct {
	permits int
	cond    *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	for s.permits <= 0 {
		s.cond.Wait()
	}
	s.permits--
	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.permits++
	s.cond.Signal()
	s.cond.L.Unlock()
}

type Channel[M any] struct {
	// Семафор емкости, необходимый для блокировки отправителя,
	// когда буфер заполнен.
	capacitySema *Semaphore

	// Семафор размера буфера, необходимый для блокировки
	// получателя, когда буфер пуст.
	sizeSema *Semaphore

	// Мьютекс, который будет защищать общую списковую структуру данных.
	mutex sync.Mutex

	// Связанный список, который будет использоваться
	// как структура данных очереди.
	buffer *list.List
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		// Создает новый семафор с количеством разрешений,
		// равным передаваемому размеру емкости.
		capacitySema: NewSemaphore(capacity),
		// Создает новый семафор с количеством разрешений, равным 0.
		sizeSema: NewSemaphore(0),
		// Создает новый, пустой связанный список.
		buffer: list.New(),
	}
}

func (c *Channel[M]) Send(message M) {
	// Получает одно разрешение от семафора емкости.
	c.capacitySema.Acquire()

	// Добавляет сообщение в очередь буфера,
	// защищая от состояния гонки с помощью мьютекса.
	c.mutex.Lock()
	c.buffer.PushBack(message)
	c.mutex.Unlock()

	// Освобождает одно разрешение от семафора размера буфера.
	c.sizeSema.Release()
}

func (c *Channel[M]) Receive() M {
	// Освобождает одно разрешение из семафора емкости.
	c.capacitySema.Release()

	// Получает одно разрешение из семафора размера буфера.
	c.sizeSema.Acquire()

	// Удаляет одно сообщение из буфера, защищая
	// от состояния гонки, используя мьютекс.
	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.mutex.Unlock()

	// Возвращает значение сообщения
	return v
}
