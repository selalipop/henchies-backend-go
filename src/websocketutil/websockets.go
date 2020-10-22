package websocketutil

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func WriteValueToWebsocket(value interface{}, conn *websocket.Conn) error {
	serializedValue, err := json.Marshal(&value)
	if err != nil {
		return fmt.Errorf("failed to serialize value for websocket: %w", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, serializedValue)
	if err != nil {
		return fmt.Errorf("failed to write value to websocket: %w", err)
	}
	return nil
}
