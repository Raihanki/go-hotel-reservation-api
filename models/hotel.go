package models

import "time"

type Hotel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"varchar(199),not null"`
	Type        string    `json:"type" gorm:"varchar(199),not null"`
	City        string    `json:"city" gorm:"varchar(199),not null"`
	Address     string    `json:"address" gorm:"text,not null"`
	Photo       string    `json:"photo" gorm:"text"`
	Description string    `json:"description" gorm:"text"`
	Rating      uint      `json:"rating" gorm:"uint"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
