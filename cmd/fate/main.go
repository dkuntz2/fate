package main

import (
	"fmt"
	"github.com/dkuntz2/fate"
)

func main() {
	for i := 1; i < 6; i++ {
		result := fate.Roll(i)
		fmt.Printf("Rolling %d Fate Dice: %d\n\trolls: %v\n", i, result.TotalResult, result.Rolls)
	}
}
