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
	controllers.InitRemindersCollection(config.DB)

	routes.UserRoutes(r)
	routes.RemindersRoutes(r)

	r.Run(":3000") // default port 8080
}
