package models
import "go.mongodb.org/mongo-driver/bson/primitive"


type Reminders struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title      string `json:"title" bson:"title"`           // หัวข้อแจ้งเตือน
	Message    string `json:"message" bson:"message"`       // ข้อความแจ้งเตือน
	Minute     int    `json:"minute" bson:"minute"`         // นาที (0-59)
	Hour       int    `json:"hour" bson:"hour"`             // ชั่วโมง (0-23)
	DayOfMonth []int  `json:"dayOfMonth" bson:"dayOfMonth"` // วันที่ในเดือน (1-31)
	Month      []int  `json:"month" bson:"month"`           // เดือน (1-12)
	DayOfWeek  []int  `json:"dayOfWeek" bson:"dayOfWeek"`   // วันในสัปดาห์ (0-6)
}
