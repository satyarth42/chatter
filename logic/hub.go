package logic

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID  string
	MsgChan chan []byte
	Conn    *websocket.Conn
}

type Message struct {
	ID          int64
	MessageType int
	Sender      string
	Receiver    string
	Body        []byte
}

type Hub struct {
	Clients     map[string]*Client
	Register    chan *Client
	Deregister  chan string
	MessageChan chan Message
}

func InitHub() *Hub {
	hub := &Hub{
		Clients:     make(map[string]*Client),
		Register:    make(chan *Client, 5),
		Deregister:  make(chan string, 5),
		MessageChan: make(chan Message, 5),
	}

	go hub.startListeners(context.Background())

	return hub
}

func (h *Hub) startListeners(ctx context.Context) {
	select {
	case user := <-h.Register:
		h.Clients[user.UserID] = user
		err := RegisterUserWithServer(ctx, user.UserID)
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("failed to register user: %s", user.UserID))
		}
	case userID := <-h.Deregister:
		delete(h.Clients, userID)
		err := DeregisterUserWithServer(ctx, userID)
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("failed to deregister user: %s", userID))
		}
	case msg := <-h.MessageChan:
		h.sendMessage(ctx, msg)
	}
}

func (h *Hub) sendMessage(ctx context.Context, msg Message) {
	conn, ok := h.Clients[msg.Receiver]
	if !ok {
		slog.WarnContext(ctx, fmt.Sprintf("failed to get connection msg:%d to user: %s", msg.ID, msg.Receiver))
	}
	err := conn.Conn.WriteMessage(msg.MessageType, msg.Body)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to send msg: %d to user: %s", msg.ID, msg.Receiver))
	}
}
