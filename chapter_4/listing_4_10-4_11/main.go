package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func matchRecorder(matchEvents *[]string, mutex *sync.RWMutex) {
	for i := 0; ; i++ {
		// Защитить "критическую секцию" с помощью write мьютекса.
		mutex.Lock()
		*matchEvents = append(*matchEvents, "Match event "+strconv.Itoa(i))
		mutex.Unlock()
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Appended match event")
	}
}

func clientHandler(mEvents *[]string, mutex *sync.RWMutex, st time.Time) {
	for i := 0; i < 100; i++ {
		// Защитить "критическую секцию" с помощью read мьютекса.
		mutex.RLock()
		allEvents := copyAllEvents(mEvents)
		mutex.RUnlock()
		timeTaken := time.Since(st)
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
	// Иниацилизировать новый reader/writer мьютекс.
	mutex := sync.RWMutex{}
	var matchEvents = make([]string, 0, 10000)

	for j := 0; j < 10000; j++ {
		matchEvents = append(matchEvents, "Match event")
	}

	// Передать reader/writer мьютекс регистратору матчей.
	go matchRecorder(&matchEvents, &mutex)

	start := time.Now()
	for j := 0; j < 5000; j++ {
		// Передать reader/writer мьютекс всем клиентским обработчикам.
		go clientHandler(&matchEvents, &mutex, start)
	}
	time.Sleep(100 * time.Second)
}
