package fate

import ()

type Fate struct {
	MessageWriter chan interface{}
	ws            *wsHub
}

func New() *Fate {
	return &Fate{
		MessageWriter: make(chan interface{}),
		ws:            newWsHub(),
	}
}

func (fate *Fate) Run() {
	go fate.ws.run()
}
