package main

import (
	"damapp-server/config"
	"damapp-server/internal/database"
	"damapp-server/internal/handler"
	"damapp-server/internal/middleware"
	"damapp-server/internal/repository"
	"damapp-server/internal/service"
	"damapp-server/utils"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()

	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	utils.NewSnowFlake(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(), 1)

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	friendshipRepo := repository.NewFriendshipRepo(db)
	friendshipService := service.NewFriendshipService(friendshipRepo)
	friendshipHandler := handler.NewFriendshipHandler(friendshipService)

	// create a new Fiber app
	app := fiber.New()

	// define routes for User endpoints
	app.Post("/users", userHandler.CreateUser)
	app.Post("/users/login", userHandler.Login)
	app.Get("/users/:id", userHandler.GetUserByID)

	// define routes for Friendship endpoints
	auth := app.Group("/friendships", middleware.JWTMiddleware())

	auth.Post("/", friendshipHandler.SendFriendRequest)
	auth.Put("/accept/:id", friendshipHandler.AcceptFriendRequest)
	auth.Put("/reject/:id", friendshipHandler.RejectFriendRequest)
	auth.Get("/pending/:userID", friendshipHandler.GetPendingRequests)
	auth.Get("/sent/:userID", friendshipHandler.GetSentRequests)

	app.Get("/friendships/are-friends/:userID1/:userID2", friendshipHandler.AreFriends)

	// Start the server and listen on port 30300
	log.Fatal(app.Listen(":30300"))
}
