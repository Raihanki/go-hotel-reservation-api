package request

import "time"

type ReservationRequest struct {
	RoomID   uint      `json:"room_id" form:"room_id" validate:"required"`
	CheckIn  time.Time `json:"check_in" form:"check_in" validate:"required,datetime"`
	CheckOut time.Time `json:"check_out" form:"check_out" validate:"datetime"`
}
