package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type BankAccount struct {
	id      string
	balance int
	mutex   sync.Mutex
}

func NewBankAccount(id string) *BankAccount {
	return &BankAccount{
		id:      id,
		balance: 100,
		mutex:   sync.Mutex{},
	}
}

func (src *BankAccount) Transfer(to *BankAccount, amount int, tellerId int) {
	// Разместить счета отправителя и получателя в слайсе.
	accounts := []*BankAccount{src, to}
	// Отсортировать счета по их ID.
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].id < accounts[j].id
	})
	fmt.Printf("%d Locking %s's account\n", tellerId, accounts[0].id)
	// Заблокировать счет с меньшим по ID порядком.
	accounts[0].mutex.Lock()
	fmt.Printf("%d Locking %s's account\n", tellerId, accounts[1].id)
	// Заблокировать счет с большим по ID порядком.
	accounts[1].mutex.Lock()
	src.balance -= amount
	to.balance += amount
	// Разблокировать оба счета.
	to.mutex.Unlock()
	src.mutex.Unlock()
	fmt.Printf("%d Unlocked %s and %s", tellerId, src.id, to.id)
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
		go func(tellerId int) {
			for j := 0; j < 1000; j++ {
				from, to := rand.Intn(total), rand.Intn(total)
				for from != to {
					to = rand.Intn(total)
				}
				accounts[from].Transfer(&accounts[to], 10, tellerId)
			}
			fmt.Println(tellerId, "COMPLETE")
		}(i)
	}

	time.Sleep(10 * time.Second)
}
