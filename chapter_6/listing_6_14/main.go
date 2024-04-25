package main

import "fmt"

const matrixSize = 3

func matrixMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
	// Итерировать по каждому ряду.
	for row := 0; row < matrixSize; row++ {
		// Итерировать по каждому столбцу.
		for col := 0; col < matrixSize; col++ {
			sum := 0
			for i := 0; i < matrixSize; i++ {
				// Суммировать каждое значение ряда из A, умножая
				// на каждое значение колонки из B.
				sum += matrixA[row][i] * matrixB[i][col]
			}
			// Обновить результат матрицы при помощи суммы.
			result[row][col] = sum
		}
	}
}

func main() {
	var result [matrixSize][matrixSize]int

	matrixA := [matrixSize][matrixSize]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	matrixB := [matrixSize][matrixSize]int{
		{9, 8, 7},
		{6, 5, 4},
		{3, 2, 1},
	}

	matrixMultiply(&matrixA, &matrixB, &result)

	fmt.Println(result)

	// [
	// 		[30 24 18]
	// 		[84 69 54]
	// 		[138 114 90]
	// ]
}
