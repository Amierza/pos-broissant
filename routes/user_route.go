package routes

import (
	"github.com/Amierza/pos-broissant/controller"
	"github.com/Amierza/pos-broissant/service"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	routes := route.Group("api/user")
	{
		routes.POST("/register", userController.Register)
		routes.POST("/login", userController.Login)
	}
}
