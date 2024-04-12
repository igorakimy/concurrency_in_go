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
		doWork(i)
	}
}

// Work 0 started at 19:48:44
// Work 0 finished at 19:48:45
// Work 1 started at 19:48:45
// Work 1 finished at 19:48:46
// Work 2 started at 19:48:46
// Work 2 finished at 19:48:47
// Work 3 started at 19:48:47
// Work 3 finished at 19:48:48
// Work 4 started at 19:48:48
// Work 4 finished at 19:48:49
