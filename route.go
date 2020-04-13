package fate

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func (fate *Fate) Router() http.Handler {
	router := chi.NewRouter()

	router.Get("/websocket", fate.Websocket)

	router.Mount("/", http.FileServer(http.Dir("static")))

	return router
}

func (fate *Fate) Websocket(rw http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(rw, req, nil)

	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"remote_addr": req.RemoteAddr,
		}).Info("Couldn't upgrade request to websocket")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Bad Request"))
		return
	}

	fate.ws.NewConn(fate, conn)
}
