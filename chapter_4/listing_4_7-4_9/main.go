package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func matchRecorder(matchEvents *[]string, mutex *sync.Mutex) {
	for i := 0; ; i++ {
		// Защитить доступ к matchEvents с помощью мьютекса.
		mutex.Lock()
		// Добавить моковую строку, которая содержит событие матча
		// каждые 200 миллисекунд.
		*matchEvents = append(*matchEvents, "Match event "+strconv.Itoa(i))
		// Разблокировать мьютекс.
		mutex.Unlock()
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Appended match event")
	}
}

func clientHandler(mEvents *[]string, mutex *sync.Mutex, st time.Time) {
	for i := 0; i < 100; i++ {
		// Защитить доступ к списку событий матча с помощью мьютекса.
		mutex.Lock()
		// Скопировать всё содержимое слайса, симулируя построение
		// ответа клиенту.
		allEvents := copyAllEvents(mEvents)
		// Разблокировать мьютекс.
		mutex.Unlock()

		// Вычислить время, которое прошло с момента старта матча.
		timeTaken := time.Since(st)
		// Вывести в консоль время, которое было потрачено на клиента.
		fmt.Println(len(allEvents), "events copied in", timeTaken)
	}
}

func copyAllEvents(matchEvents *[]string) []string {
	allEvents := make([]string, 0, len(*matchEvents))
	for _, e := range *matchEvents {
		allEvents = append(allEvents, e)
	}
	return allEvents
}

func main() {
	// Инициализировать новый мьютекс.
	mutex := sync.Mutex{}
	var matchEvents = make([]string, 0, 10000)

	for j := 0; j < 10000; j++ {
		// Предварительно заполнить срез событий матча,
		// имитируя продолжающуюся игру.
		matchEvents = append(matchEvents, "Match event")
	}

	// Запустить регистратор матчей в отдельной горутине.
	go matchRecorder(&matchEvents, &mutex)

	// Записать начальное время до начала работы
	// клиенских обработчиков запущеных в отдельных горутинах.
	start := time.Now()
	for j := 0; j < 5000; j++ {
		// Запустить большое количество горутин клиентских обработчиков.
		go clientHandler(&matchEvents, &mutex, start)
	}
	time.Sleep(100 * time.Second)
}
