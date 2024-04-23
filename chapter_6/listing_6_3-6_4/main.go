package main

import (
	"fmt"
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

type WaitGrp struct {
	// Сохранить ссылку на семафор.
	sema *Semaphore
}

func NewWaitGrp(size int) *WaitGrp {
	return &WaitGrp{
		// Инициализировать новый семафор с (1 - размер) разрешений.
		sema: NewSemaphore(1 - size),
	}
}

func (wg *WaitGrp) Wait() {
	// Вызвать Acquire() у семафора в функции Wait().
	wg.sema.Acquire()
}

func (wg *WaitGrp) Done() {
	// Когда происходит завершение, у семафора вызывается Release().
	wg.sema.Release()
}

func main() {
	// Создать группу ожидания с размерностью 4.
	wg := NewWaitGrp(4)

	for i := 1; i <= 4; i++ {
		// Создать горутину, передавая ссылку на группу ожидания.
		go doWork(i, wg)
	}
	// Подождать в группе ожидания, пока работа будет завершена.
	wg.Wait()

	fmt.Println("All complete")
}

func doWork(id int, wg *WaitGrp) {
	fmt.Println(id, "Done working")
	// Когда горутина завершает работу,
	// она вызывает Done() у группы ожидания.
	wg.Done()
}
