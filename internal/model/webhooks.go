package model

type Webhook struct {
	Model
	TenantID string   `json:"tenant_id" gorm:"index"`
	Name     string   `json:"name" gorm:"unique_index"`
	PostUrl  string   `json:"post_url"`
	Events   []string `json:"events" gorm:"type:text[];index:idx_events,using:gin"`
	IsActive bool     `json:"is_active" gorm:"index"`
}
