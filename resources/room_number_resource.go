package resources

import "time"

type RoomNumberResource struct {
	ID        uint      `json:"id"`
	RoomID    uint      `json:"room_id"`
	Number    int       `json:"number"`
	Features  string    `json:"features"`
	CreatedAt time.Time `json:"created_at"`
}

type RoomNumberWithRoomResource struct {
	ID        uint         `json:"id"`
	RoomID    uint         `json:"room_id"`
	Number    int          `json:"number"`
	Features  string       `json:"features"`
	Room      RoomResource `json:"room"`
	CreatedAt time.Time    `json:"created_at"`
}
