package request

import "time"

type ReservationRequest struct {
	RoomNumberID  uint       `json:"room_id" form:"room_number_id" validate:"required"`
	CheckIn       string     `json:"check_in" form:"check_in" validate:"required"`
	ParsedCheckin time.Time  `json:"-" form:"-" validate:"-"`
	CheckOut      *time.Time `json:"check_out" form:"check_out" validate:"-"`
	Days          uint       `json:"days" form:"days" validate:"required"`
}
