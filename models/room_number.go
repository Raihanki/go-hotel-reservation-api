package models

import "time"

type RoomNumber struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoomID    uint      `json:"room_id" gorm:"not null,foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Number    string    `json:"number" gorm:"varchar(199),not null"`
	Features  string    `json:"features" gorm:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
