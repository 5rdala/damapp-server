package handler

import (
	"damapp-server/internal/service"
	"damapp-server/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"errors"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userFromDB, err := h.service.GetUserByUserName(user.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Incorrect password"})
	}

	token, err := utils.GenerateJWT(userFromDB.ID, userFromDB.Username)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create token"})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	if err := h.service.CreateUser(user.Username, string(hashedPass)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User created"})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	ID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID: userID must be a number")
	}
	user, err := h.service.GetUserByID(uint64(ID))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"id": user.ID, "username": user.Username})
}
