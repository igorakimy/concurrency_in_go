package main

import (
	"fmt"
	"math/rand"
)

func findFactors(number int) []int {
	result := make([]int, 0)

	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}

	return result
}

func main() {
	numTasks := 10
	resultChs := make([]chan []int, numTasks)

	for i := 0; i < numTasks; i++ {
		resultChs[i] = make(chan []int)
		go func(n int) {
			resultChs[n] <- findFactors(rand.Intn(1_000_000_000))
		}(i)
	}

	for _, resultCh := range resultChs {
		fmt.Println(<-resultCh)
	}
}
