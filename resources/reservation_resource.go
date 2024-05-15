package resources

import (
	"time"

	"github.com/shopspring/decimal"
)

type ReservationResource struct {
	ID         uint               `json:"id"`
	UserID     uint               `json:"user_id"`
	RoomNumber RoomNumberResource `json:"room_number"`
	CheckIn    time.Time          `json:"check_in"`
	CheckOut   *time.Time         `json:"check_out"`
	Total      decimal.Decimal    `json:"total"`
	Status     string             `json:"status"`
	Days       uint               `json:"days"`
	CreatedAt  time.Time          `json:"created_at"`
}

type ReservationWithPayemtResource struct {
	ID                     uint               `json:"id"`
	UserID                 uint               `json:"user_id"`
	RoomNumber             RoomNumberResource `json:"room_number"`
	CheckIn                time.Time          `json:"check_in"`
	CheckOut               *time.Time         `json:"check_out"`
	Total                  decimal.Decimal    `json:"total"`
	Status                 string             `json:"status"`
	Days                   uint               `json:"days"`
	CreatedAt              time.Time          `json:"created_at"`
	StripePaymentClientKey string             `json:"stripe_payment_client_key"`
}
