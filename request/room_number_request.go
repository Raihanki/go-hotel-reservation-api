package request

type RoomNumberRequest struct {
	RoomID   uint   `json:"room_id" form:"room_id" validate:"required"`
	Number   string `json:"number" form:"number" validate:"required"`
	Features string `json:"features" form:"features" validate:"max=1000"`
}
