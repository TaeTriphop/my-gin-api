package main

import (
	"my-gin-api/config"
	"my-gin-api/controllers"
	"my-gin-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.ConnectDB()
	controllers.InitUserCollection(config.DB)

	routes.UserRoutes(r)

	r.Run(":8080") // default port 8080
}
