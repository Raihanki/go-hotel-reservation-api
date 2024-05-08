package models

import "time"

type Reservation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null,foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RoomID    uint      `json:"room_id" gorm:"not null,foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CheckIn   time.Time `json:"check_in" gorm:"timestamp;not null"`
	CheckOut  time.Time `json:"check_out" gorm:"timestamp;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
