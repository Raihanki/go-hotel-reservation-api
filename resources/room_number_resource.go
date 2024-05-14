package resources

import (
	"time"

	"github.com/shopspring/decimal"
)

type RoomNumberResource struct {
	ID        uint            `json:"id"`
	RoomID    uint            `json:"room_id"`
	Number    string          `json:"number"`
	Features  string          `json:"features"`
	Price     decimal.Decimal `json:"price"`
	CreatedAt time.Time       `json:"created_at"`
}

type RoomNumberWithRoomResource struct {
	ID        uint         `json:"id"`
	RoomID    uint         `json:"room_id"`
	Number    string       `json:"number"`
	Features  string       `json:"features"`
	Room      RoomResource `json:"room"`
	CreatedAt time.Time    `json:"created_at"`
}
