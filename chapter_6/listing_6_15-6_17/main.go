package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const matrixSize = 3

type Barrier struct {
	size      int
	waitCount int
	cond      *sync.Cond
}

func NewBarrier(size int) *Barrier {
	condVar := sync.NewCond(&sync.Mutex{})
	return &Barrier{size, 0, condVar}
}

func (b *Barrier) Wait() {
	b.cond.L.Lock()
	b.waitCount += 1
	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
	b.cond.L.Unlock()
}

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			// Заполнить каждую ячейку рандомным числом от -5 до 4.
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

func rowMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int,
	row int, barrier *Barrier) {
	// Запустить бесконечный цикл.
	for {
		// Ожидает на барьере, пока main() горутина загрузить матрицы.
		barrier.Wait()
		for col := 0; col < matrixSize; col++ {
			sum := 0
			for i := 0; i < matrixSize; i++ {
				// Вычислить результат ряда в этой горутине.
				sum += matrixA[row][i] * matrixB[i][col]
			}
			// Присвоить результат правильному ряду и столбцу.
			result[row][col] = sum
		}
		// Ожидает на барьере пока каждый другой ряд будет вычислен.
		barrier.Wait()
	}
}

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int
	// Создать новый барьер с размерностью, равной кол-ву горутин,
	// которые вычисляют ряд + главная main() горутина, итого 4.
	barrier := NewBarrier(matrixSize + 1)

	for row := 0; row < matrixSize; row++ {
		// Создать горутину на каждый ряд, передавая корректный номер ряда.
		go rowMultiply(&matrixA, &matrixB, &result, row, barrier)
	}

	for i := 0; i < 4; i++ {
		// Загрузить обе матрицы, случайным образом генерируя их.
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)

		// Снимает барьер, поэтому горутины могут начать свои вычисления.
		barrier.Wait()
		// Ожидает, пока горутины завершат свои вычисления.
		barrier.Wait()

		// Вывести результат в консоль.
		for j := 0; j < matrixSize; j++ {
			fmt.Println(matrixA[j], matrixB[j], result[j])
		}
		fmt.Println()
	}
}
