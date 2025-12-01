package logic

import (
	"Goonker/common"
	"math/rand"
	"time"
)

// SimpleBot Chooses a random empty cell for its move
func GetBotMove(logic *GameLogic) (int, int) {
	// Small delay to simulate "thinking" and make the game feel more natural
	time.Sleep(500 * time.Millisecond)

	emptyCells := [][2]int{}

	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if logic.Board[x][y] == common.Empty {
				emptyCells = append(emptyCells, [2]int{x, y})
			}
		}
	}

	if len(emptyCells) == 0 {
		return -1, -1
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	choice := emptyCells[r.Intn(len(emptyCells))]
	return choice[0], choice[1]
}