package request

type RoomRequest struct {
	HotelID     uint   `json:"hotel_id" form:"hotel_id" validate:"required"`
	Name        string `json:"name" form:"name" validate:"required,min=3,max=199"`
	Price       string `json:"price" form:"price" validate:"required,numeric"`
	MaxPeople   int    `json:"max_people" form:"max_people" validate:"required,numeric"`
	Description string `json:"description" form:"description" validate:"required,min=3,max=1000"`
	Photo       string `json:"photo" form:"photo"`
}
