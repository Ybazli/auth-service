package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ybazli/auth-service/src/config"
	"github.com/ybazli/auth-service/src/routes"
)

func main() {
	//init db
	config.InitDB()

	//init redis
	config.InitRedis()

	//run gin
	r := gin.Default()

	//register all routes
	routes.RegisterRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		panic("Run server error: " + err.Error())
	}
}
