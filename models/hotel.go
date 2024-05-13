package models

import "time"

type Hotel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(199);not null"`
	Type        string    `json:"type" gorm:"type:varchar(199);not null"`
	City        string    `json:"city" gorm:"type:varchar(199);not null"`
	Address     string    `json:"address" gorm:"type:text;not null"`
	Photo       string    `json:"photo" gorm:"type:text"`
	Description string    `json:"description" gorm:"type:text"`
	Rating      uint      `json:"rating" gorm:"type:uint"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp"`
}
