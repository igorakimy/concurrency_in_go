package main

import (
	"fmt"
	"math/rand"
	"time"
)

func player() chan string {
	output := make(chan string)
	count := rand.Intn(100)
	move := []string{"UP", "DOWN", "LEFT", "RIGHT"}
	go func() {
		defer close(output)
		for i := 0; i < count; i++ {
			output <- move[rand.Intn(4)]
			d := time.Duration(rand.Intn(200))
			time.Sleep(d * time.Millisecond)
		}
	}()
	return output
}

func handlePlayer(id int, inGame bool, direction string, players []chan string, totalPlayers *int) {
	if inGame {
		fmt.Printf("Player %d: %s\n", id, direction)
	} else {
		*totalPlayers--
		players[id] = nil
		fmt.Printf("Player %d left the game. Remaining players: %d\n", id, *totalPlayers)
	}
}

func main() {
	playersAmount := 4
	players := []chan string{
		player(),
		player(),
		player(),
		player(),
	}

	for playersAmount != 1 {
		select {
		case direction, inGame := <-players[0]:
			handlePlayer(0, inGame, direction, players, &playersAmount)
		case direction, inGame := <-players[1]:
			handlePlayer(1, inGame, direction, players, &playersAmount)
		case direction, inGame := <-players[2]:
			handlePlayer(2, inGame, direction, players, &playersAmount)
		case direction, inGame := <-players[3]:
			handlePlayer(3, inGame, direction, players, &playersAmount)
		}
	}

	fmt.Println("Game finished")
}
