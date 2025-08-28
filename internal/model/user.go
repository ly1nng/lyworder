package model

import "time"

type User struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	UserName  string     `json:"user_name" gorm:"unique;not null"`
	OpsType   string     `json:"ops_type" gorm:"default:''"`
	Role      string     `json:"role" gorm:"default:'user'"`
	Status    int        `json:"status" gorm:"not null;default:1"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"type:timestamp" json:"updated_at,omitempty"`
}
