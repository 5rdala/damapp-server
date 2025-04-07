package middleware

import (
	"damapp-server/utils"
	"github.com/gofiber/fiber/v2"

	"fmt"
	"strings"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get the token from the "Authorization" header
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization token is missing",
			})
		}

		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid token: %s", err.Error()),
			})
		}

		// store userID username in context
		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		return c.Next()
	}
}
