package main

import (
	"log"
	"net/http"

	"chat/pkg/messaging"
	"chat/pkg/messaging/ws"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func main() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	h := messaging.NewHub(l)

	l.Info("starting WebSockets Hub")

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}

	http.Handle(
		"/chat",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					l.Error("Websocket connection upgrade error", zap.Error(err))

					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer conn.Close()

				participant := ws.NewParticipant(conn)

				if err = h.Handle(participant); err != nil {
					l.Error("handling error", zap.Error(err))
				}

				w.WriteHeader(http.StatusOK)
			},
		),
	)

	err = http.ListenAndServe(":8080", nil)
	if err != http.ErrServerClosed {
		l.Fatal("HTTP server running error", zap.Error(err))
	}

	l.Info("WebSocket Hub shut down")
}
