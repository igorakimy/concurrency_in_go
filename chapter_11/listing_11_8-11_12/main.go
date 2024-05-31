package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Arbitrator struct {
	// Хранит счета с их статусом доступности.
	accountsInUse map[string]bool

	// Условная переменная используется, чтобы приостановить горутины,
	// если счета недоступны.
	cond *sync.Cond
}

func NewArbitrator() *Arbitrator {
	return &Arbitrator{
		accountsInUse: make(map[string]bool),
		cond:          sync.NewCond(&sync.Mutex{}),
	}
}

func (a *Arbitrator) LockAccounts(ids ...string) {
	// Заблокировать мьютекс на условной переменной.
	a.cond.L.Lock()
	// Итерировать до тех пор, пока все счета не станут доступны.
	for allAvailable := false; !allAvailable; {
		allAvailable = true
		for _, id := range ids {
			// Если счет используется, приостановить выполнение горутины.
			if a.accountsInUse[id] {
				allAvailable = false
				a.cond.Wait()
			}
		}
	}
	// Как только все счета станут доступны,
	// отметить требуемые счета как использованные.
	for _, id := range ids {
		a.accountsInUse[id] = true
	}
	// Разблокировать мьютекс на условной переменной.
	a.cond.L.Unlock()
}

func (a *Arbitrator) UnlockAccounts(ids ...string) {
	// Заблокировать мьютекс на условной переменной.
	a.cond.L.Lock()
	// Отметить счета как доступные.
	for _, id := range ids {
		a.accountsInUse[id] = false
	}
	// Отправить сигнал, чтобы продолжить любые приостановленные горутины.
	a.cond.Broadcast()
	// Разблокировать мьютекс на условной переменной.
	a.cond.L.Unlock()
}

type BankAccount struct {
	id      string
	balance int
}

func NewBankAccount(name string) *BankAccount {
	return &BankAccount{
		id:      name,
		balance: 10000,
	}
}

func (src *BankAccount) Transfer(
	to *BankAccount, amount int,
	tellerId int, arb *Arbitrator) {
	fmt.Printf("%d Locking %s and %s\n", tellerId, src.id, to.id)
	// Заблокировать оба счета: отправителя и получателя.
	arb.LockAccounts(src.id, to.id)
	// Выполнить перевод как только оба счета будут заблокированы.
	src.balance -= amount
	to.balance += amount
	// Разблокировать оба счета после перевода.
	arb.UnlockAccounts(src.id, to.id)
	fmt.Printf("%d Unlocked %s and %s\n", tellerId, src.id, to.id)
}

func main() {
	accounts := []BankAccount{
		*NewBankAccount("Sam"),
		*NewBankAccount("Paul"),
		*NewBankAccount("Amy"),
		*NewBankAccount("Mia"),
	}
	total := len(accounts)
	// Создать нового арбитра, чтобы использовать его в переводах.
	arb := NewArbitrator()

	for i := 0; i < total; i++ {
		go func(tellerId int) {
			for j := 0; j < 100000; j++ {
				from, to := rand.Intn(total), rand.Intn(total)
				for from != to {
					to = rand.Intn(total)
				}
				accounts[from].Transfer(&accounts[to], 10, tellerId, arb)
			}
			fmt.Println(tellerId, "COMPLETE")
		}(i)
	}

	time.Sleep(60 * time.Second)
}
