package main

import (
	"fmt"
	"time"
)

const (
	// Установить пароль, который мы будем угадывать.
	passwordToGuess = "go far"
	// Определить все возможные символы, которые могут содержаться в пароле.
	alphabet = " abcdefghijklmnopqrstuvwxyz"
)

func toBase27(n int) string {
	// Алгоритм преобразует десятичное целое число в строку
	// с основанием 27, используя алфавитную константу.
	result := ""
	for n > 0 {
		result = string(alphabet[n%27]) + result
	}
	return result
}

func guessPassword(from int, upto int, stop chan int, result chan string) {
	// Итерировать по всем комбинациям пароля, используя интервальные
	// переменные from и upto как начальную и конечную точку.
	for guessN := from; guessN < upto; guessN++ {
		select {
		// Если в стоп-канал было отправлено сообщение,
		// вывести сообщение и прекратить обработку.
		case <-stop:
			fmt.Printf("Stoped at %d [%d, %d] \n", guessN, from, upto)
			return
		default:
			match := toBase27(guessN)
			// Проверить, совпрадает ли пароль с подбираемым(в реальной
			// системе мы попытаемся получить доступ к защищенному ресурсу).
			if match == passwordToGuess {
				// Отправить сопададающий пароль в результирующий канал.
				result <- match
				// Закрыть канал, что послужит сигналом для других горутин
				// о том, что нужно прекратить проверять пароль.
				close(stop)
				return
			}
		}
	}
	fmt.Printf("Not found between [%d, %d] \n", from, upto)
}

func main() {
	// Создать общий канал, используемый горутинами для того,
	// чтобы отправить сигнал о том, что пароль найден.
	finished := make(chan int)
	// Создать канал, который будет содержать найденный пароль,
	// после окончания поиска.
	passwordFound := make(chan string)

	for i := 1; i <= 387_420_488; i += 10_000_000 {
		// Создать горутины с входящими промежутками
		// [1, 10M), [10M, 20M), ... [380M, 390M].
		go guessPassword(i, i+10_000_000, finished, passwordFound)
	}

	// Симулировать реальную программу, которая будет использовать
	// пароль для доступа к ресурсу.
	fmt.Println("password found:", <-passwordFound)
	close(passwordFound)
	// Подождать, пока пароль будет найден.
	time.Sleep(5 * time.Second)
}
