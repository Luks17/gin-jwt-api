package routes

import (
	"golang-jwt/controllers"
	"golang-jwt/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incoming_routes *gin.Engine) {
	incoming_routes.Use(middleware.Authenticate())
	incoming_routes.GET("/users", controllers.GetUsers())
	incoming_routes.GET("/users/:user_id", controllers.GetUser())
}
