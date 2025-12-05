package logic

import (
	"Goonker/common"
	"math/rand"
	"time"
)

// Constants for bot behavior
const (
	BotThinkDelay = 500 * time.Millisecond
	InvalidCoord = -1
)

// GetBotMove scans the board for available moves and selects one at random.
func GetBotMove(logic *GameLogic) (int, int) {
	// Simulate "thinking" time for natural gameplay flow
	time.Sleep(BotThinkDelay)

	// Identify all available moves
	var availableMoves [][2]int

	for x := 0; x < common.BoardSize; x++ {
		for y := 0; y < common.BoardSize; y++ {
			if logic.Board[x][y] == common.Empty {
				availableMoves = append(availableMoves, [2]int{x, y})
			}
		}
	}

	// If no moves are available, return invalid coordinates (Draw)
	if len(availableMoves) == 0 {
		return InvalidCoord, InvalidCoord
	}

	// Randomly select one of the available moves
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	choice := availableMoves[r.Intn(len(availableMoves))]
	return choice[0], choice[1]
}