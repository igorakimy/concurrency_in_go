package main

import (
	"fmt"
	"runtime"
)

func sayHello() {
	fmt.Println("Hello")
}

func main() {
	go sayHello()
	// Вызов планировщика Go дает другим горутинам шанс для запуска.
	runtime.Gosched()
	fmt.Println("Finished")
}

// Hello
// Finished
