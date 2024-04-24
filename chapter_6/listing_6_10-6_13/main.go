package main

import (
	"fmt"
	"sync"
	"time"
)

type Barrier struct {
	// Общее количество участников барьера.
	size int

	// Переменная-счетчик, которая отображает
	// количество приостановленных выполнений.
	waitCount int

	// "Условная переменная", которая используется в барьере.
	cond *sync.Cond
}

func NewBarrier(size int) *Barrier {
	// Создать новую условную переменную.
	condVar := sync.NewCond(&sync.Mutex{})
	// Создать и вернуть ссылку на новый барьер.
	return &Barrier{size, 0, condVar}
}

func (b *Barrier) Wait() {
	// Защитить доступ к переменной waitCount, использую мьютекс.
	b.cond.L.Lock()
	// Увеличить количество ожиданий на 1.
	b.waitCount += 1

	// Если waitCount достигает размера барьрера (size),
	// сбросить waitCount и сообщить "условной переменной".
	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		// Если waitCount еще не достигла размера барьера,
		// подождать "условную переменную".
		b.cond.Wait()
	}

	// Защитить доступ к переменной waitCount, использую мьютекс.
	b.cond.L.Unlock()
}

func workAndWait(name string, timeToWork int, barrier *Barrier) {
	start := time.Now()
	for {
		fmt.Println(time.Since(start), name, "is running")
		// Симулирует совершение работы, засыпая на несколько секунд.
		time.Sleep(time.Duration(timeToWork) * time.Second)
		fmt.Println(time.Since(start), name, "is waiting on barrier")
		// Подождать, пока другие горутины подтянуться.
		barrier.Wait()
	}
}

func main() {
	// Создает новый барьер, с двумя участниками.
	barrier := NewBarrier(2)

	// Запустить горутину с названием Red и временем работы 4 секунды.
	go workAndWait("Red", 4, barrier)

	// Запустить горутину с названием Blue и временем работы 10 секунд.
	go workAndWait("Blue", 10, barrier)

	// Подождать 100 секунд.
	time.Sleep(100 * time.Second)
}
