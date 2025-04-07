package handler

import (
	"damapp-server/internal/apperror"
	"damapp-server/internal/service"

	"github.com/gofiber/fiber/v2"

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

	token, err := h.service.Authenticate(user.Username, user.Password)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"token": err.Error()})
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

	token, err := h.service.CreateUser(user.Username, user.Password)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input, user id must be a number"})
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": user.ID, "username": user.Username})
}
