package handler

import (
	"damapp-server/internal/apperror"
	"damapp-server/internal/service"

	"github.com/gofiber/fiber/v2"

	"net/http"
	"strconv"
)

type FriendshipHandler struct {
	service *service.FriendshipService
}

func NewFriendshipHandler(service *service.FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{service: service}
}

func (h *FriendshipHandler) SendFriendRequest(c *fiber.Ctx) error {
	var body struct {
		senderID   uint64
		receiverID uint64
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.SendFriendRequest(body.senderID, body.receiverID); err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *FriendshipHandler) AcceptFriendRequest(c *fiber.Ctx) error {
	friendshipID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid friendship id"})
	}

	if err = h.service.AcceptFriendRequest(friendshipID); err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *FriendshipHandler) RejectFriendRequest(c *fiber.Ctx) error {
	friendshipID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid friendship id"})
	}

	if err = h.service.RejectFriendRequest(friendshipID); err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *FriendshipHandler) GetPendingRequests(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userID"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user id"})
	}

	friendships, err := h.service.GetPendingRequests(userID)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(friendships)
}

func (h *FriendshipHandler) GetSentRequests(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userID"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user id"})
	}

	friendships, err := h.service.GetSentRequests(userID)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(friendships)
}

func (h *FriendshipHandler) AreFriends(c *fiber.Ctx) error {
	userID1, err := strconv.ParseUint(c.Params("userID1"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user id 1"})
	}

	userID2, err := strconv.ParseUint(c.Params("userID2"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user id 2"})
	}

	areFriends, err := h.service.AreFriends(userID1, userID2)
	if err != nil {
		return c.Status(err.(*apperror.AppError).Code).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"are_friends": areFriends})
}
