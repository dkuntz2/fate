package fate

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type RollResult struct {
	DieCount    int
	TotalResult int
	Rolls       []string
}

func Roll(numDice int) *RollResult {
	result := &RollResult{
		DieCount:    numDice,
		TotalResult: 0,
		Rolls:       []string{},
	}

	for i := 0; i < numDice; i++ {
		thisRoll := rand.Intn(3) - 1
		switch thisRoll {
		case -1:
			result.Rolls = append(result.Rolls, "-")
		case 0:
			result.Rolls = append(result.Rolls, " ")
		case 1:
			result.Rolls = append(result.Rolls, "+")
		}

		result.TotalResult += thisRoll
	}

	return result
}
