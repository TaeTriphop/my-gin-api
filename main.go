package main

import (
	"my-gin-api/config"
	"my-gin-api/controllers"
	"my-gin-api/routes"
	"my-gin-api/scheduler"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	config.ConnectDB()
	controllers.InitUserCollection(config.DB)
	controllers.InitRemindersCollection(config.DB)

	scheduler.StartReminderScheduler(config.DB, "https://discord.com/api/webhooks/1298847619480813619/lC28xS7fPk_ELEdCAO7-E7usGJib0PHx4Kr7Wazfd9rzDQDgwth9pZ3I8djj-7eGPO5t")
	routes.UserRoutes(r)
	routes.RemindersRoutes(r)

	r.Run(":3000") // default port 8080
}
