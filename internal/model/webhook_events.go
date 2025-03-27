package model

import "time"

type SubscriberCreated struct {
	EventName  string     `json:"event_name"`
	EventTime  time.Time  `json:"event_time"`
	Subscriber Subscriber `json:"subscriber"`
	WebhookID  string     `json:"webhook_id"`
}
