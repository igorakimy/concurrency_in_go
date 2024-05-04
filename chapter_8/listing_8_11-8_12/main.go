package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateAmounts(n int) <-chan int {
	// Создается выходной канал.
	amounts := make(chan int)
	go func() {
		// Закрыть выходной канал, когда работа будет завершена.
		defer close(amounts)
		// Записать n случайных сумм в диапазоне [1, 100]
		// в выходной канал каждые 100 миллисекунд.
		for i := 0; i < n; i++ {
			amounts <- rand.Intn(100) + 1
			time.Sleep(100 * time.Millisecond)
		}
	}()
	// Вернуть выходной канал.
	return amounts
}

func main() {
	// Сгенерировать 50 сумм в канале продаж.
	sales := generateAmounts(50)
	// Сгенерировать 40 сумм в канале расходов.
	expenses := generateAmounts(40)
	endOfDayAmount := 0

	// Продолжать цикл, пока есть не нулевой канал.
	for sales != nil || expenses != nil {
		select {
		// Получать следующую сумму и флаг открытости
		// канала из канала продаж.
		case sale, moreData := <-sales:
			if moreData {
				fmt.Println("Sale of:", sale)
				// Добавить сумму продажи к общему балансу.
				endOfDayAmount += sale
			} else {
				// Если канал был закрыт, отметить канала как нулевой,
				// отключив тем самым этот случай(case) в select'е.
				sales = nil
			}
		// Получать следующую сумму и флаг открытости канала
		// из канала расходов.
		case expense, moreData := <-expenses:
			if moreData {
				fmt.Println("Expense of:", expense)
				// Отнять значение суммы расходов от общего баланса.
				endOfDayAmount -= expense
			} else {
				// Если канал был закрыт, отметить канала как нулевой,
				// отключив тем самым этот случай(case) в select'е.
				expenses = nil
			}
		}
	}

	fmt.Println("End of day profit and loss:", endOfDayAmount)
}
