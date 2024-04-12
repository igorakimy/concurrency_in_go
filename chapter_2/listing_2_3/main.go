package main

import (
	"fmt"
	"time"
)

func doWork(id int) {
	fmt.Printf("Work %d started at %s\n", id, time.Now().Format("15:04:05"))
	// Имитирует выполнение вычислительной работы, засыпая на 1 секунду.
	time.Sleep(1 * time.Second)
	fmt.Printf("Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
}

func main() {
	for i := 0; i < 5; i++ {
		// Запускает новую горутину, которая вызывает функцию doWork().
		go doWork(i)
	}
	// Ждет завершения работы всех горутин, засыпая
	// на определенный промежуток времени.
	time.Sleep(2 * time.Second)
}

// *Итоговый вывод может отличаться.

// Work 0 started at 19:51:42
// Work 4 started at 19:51:42
// Work 1 started at 19:51:42
// Work 2 started at 19:51:42
// Work 3 started at 19:51:42
// Work 4 finished at 19:51:43
// Work 3 finished at 19:51:43
// Work 0 finished at 19:51:43
// Work 2 finished at 19:51:43
// Work 1 finished at 19:51:43
