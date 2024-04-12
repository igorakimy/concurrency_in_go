package main

import (
	"fmt"
	"runtime"
)

func main() {
	// По умолчанию Go устанавливает значение GOMAXPROCS
	// в качестве возвращаемого значения методом NumCPU().
	fmt.Println("Number of CPUs:", runtime.NumCPU())

	// Вызов GOMAXPROCS(n) с n<1 возвращает текущее
	// значение логических потоков процессоров без изменений.
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
}

// Number of CPUs: 24
// GOMAXPROCS: 24
