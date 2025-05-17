package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ybazli/auth-service/src/controllers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh-token", controllers.RefreshToken)
	r.POST("/logout", controllers.Logout)
}
