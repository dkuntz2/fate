package fate

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type RollResult struct {
	TotalResult int
	Rolls       []string
}

func Roll() *RollResult {
	result := &RollResult{
		TotalResult: 0,
		Rolls:       make([]string, 4),
	}

	for i := 0; i < 4; i++ {
		thisRoll := rand.Intn(3) - 1
		switch thisRoll {
		case -1:
			result.Rolls[i] = "-"
		case 0:
			result.Rolls[i] = " "
		case 1:
			result.Rolls[i] = "+"
		}

		result.TotalResult += thisRoll
	}

	return result
}
