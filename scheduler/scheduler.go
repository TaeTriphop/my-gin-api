package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"my-gin-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

func contains(list []int, value int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func isTimeToNotify(reminder models.Reminders, now time.Time) bool {
	log.Printf("⏰ Checking reminder: %s at %02d:%02d | Now: %02d:%02d\n",
		reminder.Title, reminder.Hour, reminder.Minute, now.Hour(), now.Minute())

	if len(reminder.Month) > 0 && !contains(reminder.Month, int(now.Month())) {
		return false
	}
	if len(reminder.DayOfMonth) > 0 && !contains(reminder.DayOfMonth, now.Day()) {
		return false
	}
	if len(reminder.DayOfWeek) > 0 && !contains(reminder.DayOfWeek, int(now.Weekday())) {
		return false
	}

	match := reminder.Hour == now.Hour() && reminder.Minute == now.Minute()
	if match {
		log.Println("✅ Time matched! Sending notification")
	}
	return match
}

func sendToDiscord(webhookURL, title, message, imageURL string) {
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       title,
				"description": message,
				"color":       16626879,
				"thumbnail": map[string]string{
					"url": imageURL,
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("❌ Failed to send to Discord:", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("📤 Sent reminder '%s' to Discord with status %s\n", title, resp.Status)
}

func StartReminderScheduler(db *mongo.Database, webhookURL string) {
	go func() {
		for {
			now := time.Now()
			cursor, err := db.Collection("reminders").Find(context.TODO(), map[string]any{})
			if err != nil {
				log.Println("❌ MongoDB query error:", err)
				continue
			}

			var reminders []models.Reminders
			if err := cursor.All(context.TODO(), &reminders); err != nil {
				log.Println("❌ Cursor decode error:", err)
				continue
			}

			for _, reminder := range reminders {
				if isTimeToNotify(reminder, now) {
					sendToDiscord(webhookURL, reminder.Title, reminder.Message, reminder.ImageURL)
				}
			}

			time.Sleep(1 * time.Minute)
		}
	}()
}
