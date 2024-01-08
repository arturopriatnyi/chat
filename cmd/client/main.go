package main

import (
	"bufio"
	"chat/pkg/messaging"
	"chat/pkg/messaging/ws"
	"fmt"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

func main() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/chat", nil)
	if err != nil {
		l.Fatal("dial error", zap.Error(err))
	}
	defer conn.Close()

	participant := ws.NewParticipant(conn)

	go func() {
		for {
			message, err := participant.ReceiveMessage()
			if err != nil {
				l.Error("message reading error", zap.Error(err))
				return
			}
			fmt.Printf("%s: %s", message.Author, message.Payload)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')

		s := strings.Split(text, " ")
		if len(s) != 3 {
			l.Error("invalid message")
		}

		err := participant.SendMessage(messaging.Message{
			Type:    messaging.MessageType(s[0]),
			Author:  s[1],
			Payload: s[2],
		})
		if err != nil {
			l.Info("write message error", zap.Error(err))
			return
		}
	}
}
