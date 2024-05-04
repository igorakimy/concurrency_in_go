package main

import (
	"fmt"
	"math"
	"math/rand"
)

// primesOnly принимает числовой канал и возвращат канал,
// который содержит только простые числа.
func primesOnly(inputs <-chan int) <-chan int {
	results := make(chan int)
	// Создается анонимная горутина, которая будет фильтровать числа,
	// отбирая только простые.
	go func() {
		for c := range inputs {
			// Проверяет, не является ли число 1, т.к. 1 не простое число.
			isPrime := c != 1
			// Проверяет, делится ли нацело, число на делитель в диапазоне
			// от 2 до квадратного корня из "c".
			for i := 2; i <= int(math.Sqrt(float64(c))); i++ {
				if c%i == 0 {
					isPrime = false
					break
				}
			}
			// Если число "c" простое, записать его в результирующий канал.
			if isPrime {
				results <- c
			}
		}
	}()

	return results
}

func main() {
	numbersChannel := make(chan int)
	primes := primesOnly(numbersChannel)

	// Повторять, пока не соберем 100 простых чисел.
	for i := 0; i < 100; {
		select {
		// Подает случайное число в диапазоне от 1 до 1000000000
		// в канал numbersChannel.
		case numbersChannel <- rand.Intn(1_000_000_000) + 1:
		case p := <-primes:
			// Читает и выводит простое число.
			fmt.Println("Found prime:", p)
			i++
		}
	}
}
