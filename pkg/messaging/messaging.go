package messaging

import (
	"errors"

	"go.uber.org/zap"
)

type Participant interface {
	Username() string
	SetUsername(username string)
	ReceiveMessage() (Message, error)
	SendMessage(m Message) error
}

type Hub struct {
	l *zap.Logger

	participants []Participant
}

func NewHub(l *zap.Logger) *Hub {
	return &Hub{l: l}
}

func (h *Hub) Handle(participant Participant) error {
	for {
		message, err := participant.ReceiveMessage()
		if err != nil {
			h.l.Error("receiving message error", zap.Error(err))
		}

		participant.SetUsername(message.Author)

		switch message.Type {
		case JoinHub:
			if err := h.Join(participant); err != nil {
				h.l.Error("participant joining error", zap.Error(err))
			}
		case SendMessage:
			if err := h.Send(message); err != nil {
				h.l.Error("message sending error", zap.Error(err))
			}
		case LeaveHub:
			if err := h.Leave(participant); err != nil {
				h.l.Error("leaving hub error", zap.Error(err))
			}

			return nil
		}

		h.l.Info("received message", zap.Any("message", message))
	}
}

func (h *Hub) Join(participant Participant) error {
	for _, p := range h.participants {
		if p.Username() == participant.Username() {
			return errors.New("already joined")
		}
	}

	h.participants = append(h.participants, participant)
	h.l.Info("participant joined", zap.String("username", participant.Username()))

	return nil
}

func (h *Hub) Send(message Message) error {
	for _, p := range h.participants {
		if err := p.SendMessage(message); err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) Leave(participant Participant) error {
	for i, p := range h.participants {
		if p.Username() == participant.Username() {
			h.participants = append(h.participants[:i], h.participants[i+1:]...)

			return nil
		}
	}

	return errors.New("participant not found")
}

type MessageType string

const (
	JoinHub     MessageType = "JOIN_HUB"
	SendMessage MessageType = "SEND_MESSAGE"
	LeaveHub    MessageType = "LEAVE_HUB"
)

type Message struct {
	Type    MessageType
	Author  string
	Payload string
}
