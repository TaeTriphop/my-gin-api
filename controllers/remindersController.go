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
	"go.mongodb.org/mongo-driver/bson/primitive"

)

var remindersCollection *mongo.Collection

func InitRemindersCollection(db *mongo.Database) {
	remindersCollection = db.Collection("reminders")
}

func GetReminders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := remindersCollection.Find(ctx, bson.M{})
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

func CreateReminders(c *gin.Context) {

	var reminders models.Reminders
	if err := c.BindJSON(&reminders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminders.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := remindersCollection.InsertOne(ctx, reminders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminders"})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func GetReminder(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reminders models.Reminders
	err := remindersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&reminders)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reminders not found"})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func UpdateReminders(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var reminder models.Reminders
	if err := c.BindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       reminder.Title,
			"message":     reminder.Message,
			"minute":      reminder.Minute,
			"hour":        reminder.Hour,
			"dayOfMonth":  reminder.DayOfMonth,
			"month":       reminder.Month,
			"dayOfWeek":   reminder.DayOfWeek,
		},
	}

	result, err := remindersCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reminder"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reminder updated"})
}

func DeleteReminder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := remindersCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reminder"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reminder deleted"})
}
