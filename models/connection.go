package models

import "time"

type Connection struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	TargetID  uint `gorm:"not null"`
	CreatedAt time.Time
}
