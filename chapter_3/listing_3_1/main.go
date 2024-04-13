package main

import (
	"fmt"
	"time"
)

func main() {
	// Выделить пространство в памяти для целочисленной переменной.
	count := 5

	// Запустить горутину и передать ссылку на переменную в памяти.
	go countdown(&count)

	// Горутина main() читает значение общей переменной каждые полсекунды.
	for count > 0 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(count)
	}
}

func countdown(seconds *int) {
	for *seconds > 0 {
		time.Sleep(1 * time.Second)
		// Горутина обновляет значение общей переменной.
		*seconds -= 1
	}
}

// *Итоговый вывод может отличаться.

// 5
// 4
// 4
// 3
// 3
// 2
// 2
// 1
// 1
// 0
