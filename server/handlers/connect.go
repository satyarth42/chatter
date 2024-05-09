package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/satyarth42/chatter/auth"
	"github.com/satyarth42/chatter/logic"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // For demonstration purposes only. Set restrictions in production.
}
var hub = logic.InitHub()

func Connect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.Header.Get(auth.HeaderUserID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("could not upgrade connection, err: %+v", err))
		return
	}
	defer conn.Close()
	defer func() {
		hub.Deregister <- userID
	}()

	hub.Register <- &logic.Client{
		UserID:  userID,
		MsgChan: make(chan []byte, 10),
		Conn:    conn,
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("error in reading message. err: %+v", err))
		}
		msg := logic.Message{}
		jsonErr := json.Unmarshal(message, &msg)
		if jsonErr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("error in json unmarshal of message: %+v", message))
		}
		msg.MessageType = messageType
		hub.MessageChan <- msg
	}
}
