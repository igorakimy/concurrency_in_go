package main

import "sync"

type WaitGrp struct {
	// Размер группы ожидания, по умолчанию равен 0.
	groupSize int

	// Условная переменная, которая будет использована в группе ожидания.
	cond *sync.Cond
}

func NewWaitGrp() *WaitGrp {
	return &WaitGrp{
		// Инициализировать условную переменную с новым мьютексом.
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (wg *WaitGrp) Add(delta int) {
	// Защитить обновление groupSize при помощи
	// блокировки мьютекса на "условной переменной".
	wg.cond.L.Lock()
	// Увеличить groupSize на значение delta.
	wg.groupSize += delta
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Wait() {
	// Защитить чтение groupSize при помощи
	// блокировки мьютекса на "условной переменной".
	wg.cond.L.Lock()
	for wg.groupSize > 0 {
		// Подождать и атомарно совобождать мьютекс,
		// пока groupSize больше 0.
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}

func (wg *WaitGrp) Done() {
	// Защитить обновление groupSize при помощи
	// блокировки мьютекса на "условной переменной".
	wg.cond.L.Lock()
	// Уменьшить groupSize на 1.
	wg.groupSize--
	if wg.groupSize == 0 {
		// Если это последняя горутина, которая должна быть завершена
		// в группе ожидания, она сообщает об этом "условной переменной".
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}
