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

func main() {
	// Выполнить 10 раз.
	for i := 0; i < 10; i++ {
		// Выполнять один шаг за другим последовательно.
		result := Box(AddToppings(Bake(Mixture(PrepareTray(i)))))
		fmt.Println("Accepting", result)
	}
}
