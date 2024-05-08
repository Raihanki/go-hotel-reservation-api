package resources

import (
	"time"

	"github.com/shopspring/decimal"
)

type RoomResource struct {
	ID          uint              `json:"id"`
	Hotel       HotelMiniResource `json:"hotel"`
	Name        string            `json:"name"`
	Price       decimal.Decimal   `json:"price"`
	MaxPeople   int               `json:"max_people"`
	Description string            `json:"description"`
	Photo       string            `json:"photo"`
	CreatedAt   time.Time         `json:"created_at"`
}

type RoomWithRoomNumberResource struct {
	ID          uint               `json:"id"`
	Hotel       HotelMiniResource  `json:"hotel"`
	Name        string             `json:"name"`
	Price       decimal.Decimal    `json:"price"`
	MaxPeople   int                `json:"max_people"`
	Description string             `json:"description"`
	Photo       string             `json:"photo"`
	RoomNumber  RoomNumberResource `json:"room_number"`
	CreatedAt   time.Time          `json:"created_at"`
}
