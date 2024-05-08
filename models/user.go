package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"varchar(199),not null"`
	Email     string    `json:"email" gorm:"varchar(199),not null"`
	Password  string    `json:"password" gorm:"varchar(199),not null"`
	IsAdmin   bool      `json:"is_admin" gorm:"bool;default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
