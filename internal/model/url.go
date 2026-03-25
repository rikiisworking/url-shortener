package model

import "time"

type URL struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	OriginalURL string     `gorm:"not null;size:2048" json:"original_url" validate:"required,url"`
	ShortCode   string     `gorm:"unique;size:10;not null" json:"short_code"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	ClickCount  int64      `gorm:"default:0" json:"click_count"`
}
