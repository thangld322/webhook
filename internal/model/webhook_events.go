package model

import "time"

type WebhookEvent struct {
	TenantID   string      `json:"tenant_id"`
	EventName  string      `json:"event_name"`
	EventTime  time.Time   `json:"event_time"`
	Subscriber *Subscriber `json:"subscriber"`
	WebhookID  string      `json:"webhook_id"`
}
