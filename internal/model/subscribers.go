package model

type Subscriber struct {
	Model
	TenantID  string `json:"tenant_id" gorm:"index"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
