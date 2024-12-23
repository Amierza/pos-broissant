package main

import (
	"log"
	"os"

	"github.com/Amierza/pos-broissant/cmd"
	"github.com/Amierza/pos-broissant/config"
	"github.com/Amierza/pos-broissant/controller"
	"github.com/Amierza/pos-broissant/middleware"
	"github.com/Amierza/pos-broissant/repository"
	"github.com/Amierza/pos-broissant/routes"
	"github.com/Amierza/pos-broissant/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	var (
		jwtService     service.JWTService        = service.NewJWTService()
		userRepository repository.UserRepository = repository.NewUserRespository(db)
		userService    service.UserService       = service.NewUserService(userRepository, jwtService)
		userController controller.UserController = controller.NewUserController(userService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.User(server, userController, jwtService)

	server.Static("/assets", "./assets")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
