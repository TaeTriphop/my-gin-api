package controllers

import (
	"context"
	"log"
	"my-gin-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var remindersController *mongo.Collection

func InitRemindersCollection(db *mongo.Database) {
	remindersController = db.Collection("reminders")
}

func GetReminders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := remindersController.Find(ctx, bson.M{})
	log.Println("cursor")
	log.Println(cursor)
	if err != nil {
		log.Println("❌ Find error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	log.Println("✅ Find success")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reminders"})
		return
	}

	var reminders []models.Reminders
	if err = cursor.All(ctx, &reminders); err != nil {
		log.Println("❌ Decode error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding reminders"})
		return
	}
	log.Printf("✅ Found %d reminders\n", len(reminders))
	c.JSON(http.StatusOK, reminders)
}
