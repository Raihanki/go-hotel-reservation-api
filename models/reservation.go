package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Reservation struct {
	ID               uint            `json:"id" gorm:"primaryKey"`
	UserID           uint            `json:"user_id" gorm:"not null;foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RoomNumberID     uint            `json:"room_number_id" gorm:"not null;foreignKey:RoomNumberID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RoomNumber       RoomNumber      `json:"room_number"`
	CheckIn          time.Time       `json:"check_in" gorm:"type:timestamp;not null"`
	CheckOut         *time.Time      `json:"check_out" gorm:"type:timestamp"`
	Total            decimal.Decimal `json:"total" gorm:"type:decimal(10,2);not null"`
	Status           string          `json:"status" gorm:"not null"`
	PaymentSuccessAt *time.Time      `json:"payment_success_at" gorm:"type:timestamp"`
	PaymentFailedAt  *time.Time      `json:"payment_failed_at" gorm:"type:timestamp"`
	Days             uint            `json:"days" gorm:"not null"`
	CreatedAt        time.Time       `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt        time.Time       `json:"updated_at" gorm:"type:timestamp"`
}
