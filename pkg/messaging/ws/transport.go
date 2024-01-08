package ws

import (
	"chat/pkg/messaging"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type participant struct {
	username string

	conn *websocket.Conn
}

func NewParticipant(conn *websocket.Conn) messaging.Participant {
	return &participant{conn: conn}
}

func (p *participant) Username() string {
	return p.username
}

func (p *participant) SetUsername(username string) {
	p.username = username
}

func (p *participant) SendMessage(message messaging.Message) error {
	return p.conn.WriteMessage(
		1,
		[]byte(fmt.Sprintf("%s %s %s", message.Type, message.Author, message.Payload)),
	)
}

func (p *participant) ReceiveMessage() (messaging.Message, error) {
	_, m, err := p.conn.ReadMessage()
	if err != nil {
		return messaging.Message{}, err
	}

	s := strings.Split(string(m), " ")
	if len(s) != 3 {
		return messaging.Message{}, errors.New("invalid message")
	}

	return messaging.Message{
		Type:    messaging.MessageType(s[0]),
		Author:  s[1],
		Payload: s[2],
	}, nil
}
