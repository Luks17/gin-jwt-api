package routes

import (
	"golang-jwt/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incoming_routes *gin.Engine) {
	incoming_routes.POST("/users/signup", controllers.Signup())
	incoming_routes.POST("/users/login", controllers.Login())
}
