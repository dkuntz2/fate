package fate

import ()

type Character struct {
	Name         string
	Player       string
	Refresh      int8
	FatePoints   int8
	HighConcept  string
	Trouble      string
	Aspects      []string
	Approaches   *Approaches
	Stress       int8
	Consequences []string
}

type Approaches struct {
	Careful  int8
	Clever   int8
	Flashy   int8
	Forceful int8
	Quick    int8
	Sneaky   int8
}
