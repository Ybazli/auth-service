package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ybazli/auth-service/src/controllers"
	"github.com/ybazli/auth-service/src/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh-token", controllers.RefreshToken)
	r.POST("/logout", controllers.Logout)

	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())

	auth.GET("/me", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		email := c.GetString("email")
		role := c.GetString("role")

		c.JSON(200, gin.H{
			"user_id": userID,
			"email":   email,
			"role":    role,
		})
	})
}
