package fate

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type wsHub struct {
	connections map[*wsClient]bool
	broadcast   chan interface{}
	register    chan *wsClient
	unregister  chan *wsClient
}

type wsClient struct {
	fate   *Fate
	conn   *websocket.Conn
	hub    *wsHub
	sender chan interface{}
}

type wsMessage struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

type wsBroadcast struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type rollMessage struct {
	Who string `json:"who"`
}

type rollResult struct {
	Who    string   `json:"who"`
	Rolls  []string `json:"rolls"`
	Result int      `json:"result"`
}

type wsErrorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var (
	errBadMessageType = &wsErrorMessage{Type: "error", Message: "Bad message type"}
)

func newWsHub() *wsHub {
	return &wsHub{
		connections: make(map[*wsClient]bool),
		broadcast:   make(chan interface{}),
		register:    make(chan *wsClient),
		unregister:  make(chan *wsClient),
	}
}

func (hub *wsHub) run() {
	// log.Info("starting websocket hub")
	for {
		select {
		case client := <-hub.register:
			// log.WithFields(log.Fields{
			// 	"address": client.conn.LocalAddr(),
			// }).Debug("Client Connect")
			hub.connections[client] = true

		case client := <-hub.unregister:
			// log.WithFields(log.Fields{
			// 	"address": client.conn.LocalAddr(),
			// }).Debug("Client Disconnect")
			_, ok := hub.connections[client]
			if ok {
				delete(hub.connections, client)
			}

		case message := <-hub.broadcast:
			// log.Info("sending message to all clients")
			for client := range hub.connections {
				client.sender <- message
			}
		}
	}
}

func (hub *wsHub) NewConn(fate *Fate, conn *websocket.Conn) {
	client := &wsClient{
		fate:   fate,
		conn:   conn,
		hub:    hub,
		sender: make(chan interface{}),
	}

	go client.run()
}

func (client *wsClient) run() {
	// log.Info("new websocket connection")
	client.hub.register <- client
	go client.writer()
	go client.reader()
}

func (client *wsClient) writer() {
	defer func() {
		client.conn.Close()
	}()

	for {
		message := <-client.sender
		// log.Info("picked up message to send to websocket")
		err := client.conn.WriteJSON(message)
		if err != nil {
			log.WithError(err).Info("couldn't write message")
			break
		}
	}
}

func (client *wsClient) reader() {
	defer func() {
		client.hub.unregister <- client
		client.conn.Close()
	}()

	for {
		// log.Info("waiting for message")
		var message wsMessage
		err := client.conn.ReadJSON(&message)
		// log.Info("read message")
		if err != nil {
			log.WithError(err).Error("Bad message from client")
			break
		}

		// log.Info("message type ", message.Type)

		switch message.Type {
		case "roll":
			var rollMsg rollMessage
			err := json.Unmarshal(message.Message, &rollMsg)
			if err != nil {
				log.WithError(err).Error("Bad roll message from client")
				client.sender <- errBadMessageType
				continue
			}

			roll := Roll()
			brdMsg := &wsBroadcast{
				Type: "roll",
				Message: &rollResult{
					Who:    rollMsg.Who,
					Rolls:  roll.Rolls,
					Result: roll.TotalResult,
				},
			}

			client.hub.broadcast <- brdMsg

		default:
			client.sender <- errBadMessageType
		}
	}
}
