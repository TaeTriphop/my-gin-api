package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	if len(reminder.Month) > 0 && !contains(reminder.Month, int(now.Month())) {
		return false
	}
	if len(reminder.DayOfMonth) > 0 && !contains(reminder.DayOfMonth, now.Day()) {
		return false
	}
	if len(reminder.DayOfWeek) > 0 && !contains(reminder.DayOfWeek, int(now.Weekday())) {
		return false
	}
	return reminder.Hour == now.Hour() && reminder.Minute == now.Minute()
}

func sendToDiscord(webhookURL, title, message string) {
	payload := map[string]string{
		"content": fmt.Sprintf("**%s**\n%s", title, message),
	}
	jsonBody, _ := json.Marshal(payload)

	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("❌ Failed to send to Discord:", err)
	}
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
					sendToDiscord(webhookURL, reminder.Title, reminder.Message)
				}
			}

			time.Sleep(1 * time.Minute)
		}
	}()
}
