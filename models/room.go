package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Room struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	HotelID     uint            `json:"hotel_id" gorm:"not null,foreignKey:HotelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name        string          `json:"name" gorm:"type:varchar(199);not null"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not null"`
	MaxPeople   int             `json:"max_people" gorm:"type:int;not null"`
	Description string          `json:"description" gorm:"type:text;not null"`
	Photo       string          `json:"photo" gorm:"type:text"`
	CreatedAt   time.Time       `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"type:timestamp"`
}
