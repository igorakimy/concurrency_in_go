package main

import (
	"fmt"
	"time"
)

const (
	ovenTime           = 5
	everyThingElseTime = 2
)

func PrepareTray(trayNumber int) string {
	fmt.Println("Preparing empty tray", trayNumber)
	// На каждом шагу изготовления кекса засыпать на 2 секунды,
	// симулируя реальную работу.
	time.Sleep(everyThingElseTime * time.Second)
	// Каждая функция возвращает описание того, что было сделано.
	return fmt.Sprintf("tray number %d", trayNumber)
}

func Mixture(tray string) string {
	fmt.Println("Pouring cupcake Mixture in", tray)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("cupcake in %s", tray)
}

func Bake(mixture string) string {
	fmt.Println("Baking", mixture)
	// На шаге выпекания заснуть на 5 секунд вместо 2.
	time.Sleep(ovenTime * time.Second)
	return fmt.Sprintf("baked %s", mixture)
}

func AddToppings(bakedCupCake string) string {
	fmt.Println("Adding topping to", bakedCupCake)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("topping on %s", bakedCupCake)
}

func Box(finishedCupCake string) string {
	fmt.Println("Boxing", finishedCupCake)
	time.Sleep(everyThingElseTime * time.Second)
	return fmt.Sprintf("%s boxed", finishedCupCake)
}

func AddOnPipe[X, Y any](q <-chan int, f func(X) Y, in <-chan X) chan Y {
	// Создать выходящий канал с типом Y.
	output := make(chan Y)
	// Запустить горутину.
	go func() {
		defer close(output)
		// Вызывать select в бесконечном цикле.
		for {
			select {
			// Когда канал q закроется, выйти из цикла
			// и прекратить выполнение горутины.
			case <-q:
				return
			// Получает сообщение из входящего канала,
			// если он доступен.
			case input := <-in:
				// Вызывает функцию f и отправляет результат,
				// который вернула функция в канал output.
				output <- f(input)
			}
		}
	}()
	return output
}

func main() {
	// Создать входящий канал, который будет использован
	// чтобы подключиться к первому шагу.
	input := make(chan int)
	// Создать канал для выхода.
	quit := make(chan int)

	// Подключить каждый шаг конвейера, передавая выходное значение
	// каждого шага на вход к другому шагу.
	output := AddOnPipe(quit, Box,
		AddOnPipe(quit, AddToppings,
			AddOnPipe(quit, Bake,
				AddOnPipe(quit, Mixture,
					AddOnPipe(quit, PrepareTray, input)))))

	// Создать горутину, которая отправляет 10 целых чисел
	// в конвейер, чтобы изготовить 10 коробок кексов.
	go func() {
		for i := 0; i < 10; i++ {
			input <- i
		}
	}()

	// Прочитать 10 коробок с кексами из последнего шага конвейера.
	for i := 0; i < 10; i++ {
		fmt.Println(<-output, "received")
	}
}
