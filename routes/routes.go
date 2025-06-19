package routes

import (
	"my-gin-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/users", controllers.CreateUser)
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUser)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)
}

func RemindersRoutes(r *gin.Engine) {
	r.GET("/reminders", controllers.GetReminders)
	r.POST("/reminders", controllers.CreateReminders)
	r.GET("/reminders/:id", controllers.GetReminder)
	r.PUT("/reminders/:id", controllers.UpdateReminders)
	r.DELETE("/reminders/:id", controllers.DeleteReminder)

}
