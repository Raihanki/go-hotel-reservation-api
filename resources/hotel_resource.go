package resources

import "time"

type HotelResource struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	City        string    `json:"city"`
	Address     string    `json:"address"`
	Photo       string    `json:"photo"`
	Description string    `json:"description"`
	Rating      uint      `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

type HotelMiniResource struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Rating  uint   `json:"rating"`
	Address string `json:"address"`
}
