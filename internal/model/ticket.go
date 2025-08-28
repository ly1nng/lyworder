package model

import "time"

type Ticket struct {
	ID           string     `json:"id" gorm:"primary_key"`
	TicketType   string     `json:"ticket_type"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Status       string     `json:"status"`
	UserID       string     `json:"user_id"`
	UserName     string     `json:"user_name"`
	OperatorName string     `json:"operator_name"`
	Remark       string     `json:"remark"`
	Solution     string     `json:"solution"`
	CreatedAt    time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt    *time.Time `gorm:"type:timestamp" json:"updated_at,omitempty"`
	ScreenShots  []string   `json:"screenshots,omitempty" gorm:"type:json;serializer:json"`
}
