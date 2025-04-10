package handler

import (
	"damapp-server/internal/apperror"
	"damapp-server/internal/service"

	"github.com/gofiber/fiber/v2"

	"net/http"
	"strconv"
)

type MatchHandler struct {
	service *service.MatchService
}

func NewMatchHandler(service *service.MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) CreateMatch(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)

	match, err := h.service.CreateMatch(userID)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(match)
}

func (h *MatchHandler) GetMatchByID(c *fiber.Ctx) error {
	matchID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input, match id must be a number"})
	}

	match, err := h.service.GetByID(matchID)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(match)
}

func (h *MatchHandler) JoinMatch(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)

	matchCode, err := strconv.Atoi(c.Params("code"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input, match code must be a number"})
	}

	match, err := h.service.JoinMatch(matchCode, userID)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(match)
}
