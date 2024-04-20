package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Создать новую "условную переменную".
	cond := sync.NewCond(&sync.Mutex{})
	// Инициализировать общее количество игроков, равное 4.
	playersInGame := 4

	for playerId := 0; playerId < 4; playerId++ {
		// Запустить горутину, передавая ей общую "условную переменную",
		// количество игроков и идентификатор игрока.
		go playerHandler(cond, &playersInGame, playerId)
		// Засыпать на 1 секунду, до тогоа как следующий игрок приконнектится.
		time.Sleep(1 * time.Second)
	}
}

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int) {
	// Заблокировать мьютекс "условной переменной"
	// чтобы избежать состояний гонки.
	cond.L.Lock()

	fmt.Println(playerId, ": Connected")
	// Вычесть 1 из общего количества оставшихся игроков.
	*playersRemaining--
	if *playersRemaining == 0 {
		// Отправить трансляцию, когда все игроки подключатся.
		cond.Broadcast()
	}

	for *playersRemaining > 0 {
		fmt.Println(playerId, ": Waiting for more players")
		// Ожидает в зависимости от "условной переменной"
		// до тех пор, пока не подключится больше игроков.
		cond.Wait()
	}

	// Разблокировать мьютекс, чтобы все горутины могли
	// возобновить выполнение и запустить игру.
	cond.L.Unlock()

	fmt.Println("All players connected. Ready player", playerId)
	// Игра началась
}
