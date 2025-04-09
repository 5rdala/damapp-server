package service

import (
	"github.com/gofiber/contrib/websocket"

	"fmt"
	"log"
	"sync"
)

type WebSocketService struct {
	clients map[uint64]*websocket.Conn
	mu      sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[uint64]*websocket.Conn),
	}
}

func (s *WebSocketService) RegisterConnection(userID uint64, username string, conn *websocket.Conn) {
	log.Println("[WS] registring conn:", userID)

	s.mu.Lock()
	s.clients[userID] = conn
	s.mu.Unlock()

	log.Printf("[WS] %s connected", username)

	conn.SetCloseHandler(func(code int, text string) error {
		s.mu.Lock()
		delete(s.clients, userID)
		s.mu.Unlock()
		log.Printf("[WS] %d disconnected", userID)
		return nil
	})
}

func (s *WebSocketService) SendTo(userID uint64, username string, message []byte) error {
	s.mu.RLock()
	conn, ok := s.clients[userID]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("user %s not connected", username)
	}

	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		s.mu.Lock()
		delete(s.clients, userID)
		s.mu.Unlock()
	}

	return nil
}
