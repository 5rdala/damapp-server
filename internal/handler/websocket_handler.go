package handler

import (
	"damapp-server/internal/service"
	"damapp-server/utils"

	"github.com/gofiber/contrib/websocket"

	"strings"
)

type WebSocketHandler struct {
	service *service.WebSocketService
}

func NewWebSocketHandler(service *service.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{service: service}
}

func (h *WebSocketHandler) HandleWebSocket() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		auth := c.Headers("Authorization")
		if auth == "auth" {
			_ = c.Close()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			_ = c.Close()
			return
		}

		h.service.RegisterConnection(claims.UserID, claims.Username, c)

		// read messages (just to keep conn alive)
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}
