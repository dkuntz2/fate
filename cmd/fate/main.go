package main

import (
	"fmt"
	"github.com/dkuntz2/fate"
)

func main() {
	db := fate.ProvideDb()
	char := &fate.Character{
		Player:      "Don",
		Name:        "Ford Prefect",
		Refresh:     1,
		FatePoints:  1,
		HighConcept: "Just this guy, you know?",
		Trouble:     "Part of my brain was removed so i could be president",
		Aspects: []string{
			"Stealer of The Heart of Gold",
		},
		Approaches: &fate.Approaches{
			Careful:  0,
			Clever:   1,
			Flashy:   3,
			Forceful: 1,
			Quick:    2,
			Sneaky:   2,
		},
		Stress:       0,
		Consequences: []string{},
	}

	err := db.SaveCharacter(char)
	if err != nil {
		panic(err)
	}

	chars, err := db.AllCharacters()
	if err != nil {
		panic(err)
	}

	for _, char := range chars {
		fmt.Printf("%+v\n\n", char)
	}
}
