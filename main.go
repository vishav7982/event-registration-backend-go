package main

import (
	"restapp/db"
	"restapp/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDb()
	server := gin.Default() // It gives us basic http server with logger and recovery middleware
	routes.RegisterRoutes(server)
	server.Run(":8080") // localhost:8080
}
