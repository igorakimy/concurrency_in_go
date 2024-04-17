package main

import "sync"

type ReadWriteMutex struct {
	// Целочисленная переменная для подсчета количества
	// читающих горутин, которые конкурентно читают
	// "критическую секцию".
	readersCounter int

	// Мьютекс для синхронизации доступа читателей.
	readersLock sync.Mutex

	// Мьютекс для блокировки любого доступа писателям.
	globalLock sync.Mutex
}

func (rw *ReadWriteMutex) ReadLock() {
	// Синхронизирует доступ таким образом, что в любой
	// момент времени разрешено выполнение только одной горутине.
	rw.readersLock.Lock()
	// Читающая горутина увеличивает счетчик readersCounter на 1.
	rw.readersCounter++
	if rw.readersCounter == 1 {
		// Если читающая горутина запускается первой, она
		// пытатеся заблокировать globalLock на запись.
		rw.globalLock.Lock()
	}
	// Синхронизирует доступ таким образом, что в любой
	// момент времени разрешено выполнение только одной горутине.
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	// Любая пишущая горутина требует блокировки globalLock на запись.
	rw.globalLock.Lock()
}

func (rw *ReadWriteMutex) ReadUnlock() {
	// Синхронизирует доступ таким образом, что в любой
	// момент времени разрешено выполнение только одной горутине.
	rw.readersLock.Lock()
	// Читающая горутина уменьшает readersCounter на 1.
	rw.readersCounter--
	if rw.readersCounter == 0 {
		// Если читающая горутина является последней,
		// она разблокирует глобальную блокировку(на запись).
		rw.globalLock.Unlock()
	}
	// Синхронизирует доступ таким образом, что в любой
	// момент времени разрешено выполнение только одной горутине.
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	// Пишущая горутина, завершая запись в критическую секцию,
	// освобождает глобальную блокировку.
	rw.globalLock.Unlock()
}
