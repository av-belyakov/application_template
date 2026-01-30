package websocketserver

import (
	"github.com/gorilla/websocket"

	"github.com/av-belyakov/application_template/interfaces"
)

// Client реализованное подключение
type Client struct {
	conn   *websocket.Conn
	logger interfaces.Logger
	send   chan []byte
}

// Hub управление всеми клиентами
type Hub struct {
	clients        map[*Client]bool
	chBroadcast    chan []byte
	chIncomingData chan []byte
	chRegister     chan *Client
	chUnregister   chan *Client
}
