package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type BankAccount struct {
	id      string
	balance int
	mutex   sync.Mutex
}

func NewBankAccount(id string) *BankAccount {
	// Создает новый экземпляр банковского счета
	// с 100$ и новым мьютексом.
	return &BankAccount{
		id:      id,
		balance: 100,
		mutex:   sync.Mutex{},
	}
}

func (src *BankAccount) Transfer(to *BankAccount, amount int, exId int) {
	fmt.Printf("%d Locking %s's account\n", exId, src.id)
	// Заблокировать мьютекс на аккаунте отправителя.
	src.mutex.Lock()
	fmt.Printf("%d Locking %s's account\n", exId, to.id)
	// Заблокировать мьютекс на аккаунте получателя.
	to.mutex.Lock()
	// Уменьшить деньги у отправителя и добавить их на аккаунт получателя.
	src.balance -= amount
	to.balance += amount
	// Разблокировать оба мьютекса на аккаунтах отправителя и получателя.
	to.mutex.Unlock()
	src.mutex.Unlock()
	fmt.Printf("%d Unlocked %s and %s\n", exId, src.id, to.id)
}

func main() {
	accounts := []BankAccount{
		*NewBankAccount("Sam"),
		*NewBankAccount("Paul"),
		*NewBankAccount("Amy"),
		*NewBankAccount("Mia"),
	}

	total := len(accounts)
	for i := 0; i < total; i++ {
		// Создать горутину с уникальным ID выполнения.
		go func(eId int) {
			// Выполнить 1000 рандомно сгенерированных переводов.
			for j := 1; j < 1000; j++ {
				// Случайным образом выбрать аккаунты оптравителя и получателя.
				from, to := rand.Intn(total), rand.Intn(total)
				for from == to {
					to = rand.Intn(total)
				}
				// Выполнить перевод.
				accounts[from].Transfer(&accounts[to], 10, eId)
			}
			// Как только все 1000 переводов будут завершены,
			// вывести соответствующее сообщение.
			fmt.Println(eId, "COMPLETE")
		}(i)
	}
	// Подождать 60 секунд до завершения программы.
	time.Sleep(60 * time.Second)
}
