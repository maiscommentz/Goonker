package logic

import (
	"Goonker/server/assets"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/bits-and-blooms/bitset"
)

// ChallengeManager handles the challenges
type ChallengeManager struct {
	challenges      []Challenge
	askedChallenges bitset.BitSet
}

// Challenge represents a challenge
type Challenge struct {
	Question  string   `json:"question"`
	Answers   []string `json:"answers"`
	AnswerKey int      `json:"answer_key"`
}

// NewChallengeManager creates a new challenge manager
func NewChallengeManager() *ChallengeManager {
	// Load challenges from json file
	challengesByte, err := assets.AssetsFS.ReadFile("challenges.json")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize challenge manager
	challengeManager := &ChallengeManager{}

	if err := json.Unmarshal(challengesByte, &challengeManager.challenges); err != nil {
		log.Fatal(err)
	}

	// Initialize asked challenges
	challengeManager.askedChallenges = *bitset.New(uint(len(challengeManager.challenges)))

	return challengeManager
}

// PickChallenge returns a challenge from the challenges list
func (m *ChallengeManager) PickChallenge() (*Challenge, error) {
	if m.challenges == nil {
		return nil, fmt.Errorf("no challenges loaded")
	}

	// To avoid picking the same challenge multiple times in the same game,
	// we use a bitset to store the challenges that have already been picked
	randIndex := 0
	for {
		randIndex = rand.Intn(len(m.challenges))

		// In case every challenges have been picked once, the bitset is cleared
		if m.askedChallenges.All() {
			m.askedChallenges.ClearAll()
		}

		if !m.askedChallenges.Test(uint(randIndex)) {
			m.askedChallenges.Set(uint(randIndex))
			break
		}
	}

	challenge := &m.challenges[randIndex]
	return challenge, nil
}

// Shuffle the order of the answers
func (c *Challenge) Shuffle() {

	for i := range c.Answers {
		j := rand.Intn(i + 1)

		// Keep track of the new position of the answer key
		if i == c.AnswerKey {
			c.AnswerKey = j
		} else if j == c.AnswerKey {
			c.AnswerKey = i
		}

		c.Answers[i], c.Answers[j] = c.Answers[j], c.Answers[i]
	}
}
