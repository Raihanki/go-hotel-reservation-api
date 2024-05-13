package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(199);not null"`
	Email     string    `json:"email" gorm:"type:varchar(199);not null"`
	Password  string    `json:"password" gorm:"type:varchar(199);not null"`
	IsAdmin   bool      `json:"is_admin" gorm:"type:bool;default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}
