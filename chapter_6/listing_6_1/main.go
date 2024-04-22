package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// Создать новую группу ожидания.
	wg := sync.WaitGroup{}
	// Добавить 4 в группу ожидания, т.к. имеется 4 горутины,
	// которые должны отработать.
	wg.Add(4)

	// Создать 4 горутины, передавая ссылку на группу ожидания.
	for i := 1; i <= 4; i++ {
		go doWork(i, &wg)
	}

	// Ждать, пока все горутины не завершат свою работу.
	wg.Wait()
	fmt.Println("All complete")
}

func doWork(id int, wg *sync.WaitGroup) {
	i := rand.Intn(5)
	// Заснуть на рандомное кол-во времени(до 5 секунд).
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Println(id, "Done working after", i, "seconds")
	// Сигнализировать, что горутина завершила задачу.
	wg.Done()
}
