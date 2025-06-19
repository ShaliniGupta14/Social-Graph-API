package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string  `json:"name"`
	Email       string  `json:"email" gorm:"unique"`
	Connections []*User `gorm:"many2many:user_connections;" json:"connections"`
}
