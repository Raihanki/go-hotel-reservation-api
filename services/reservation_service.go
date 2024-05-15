package services

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Raihanki/go-hotel-reservation-api/configs"
	"github.com/Raihanki/go-hotel-reservation-api/models"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"gorm.io/gorm"
)

type ReservationServiceInterface interface {
	CreateReservation(userID int, price decimal.Decimal, roomID int, reservationRequest request.ReservationRequest) (resources.ReservationWithPayemtResource, error)
	GetReservationByID(userID int, reservationID int) (resources.ReservationResource, error)
	GetUserReservation(userID int) ([]resources.ReservationResource, error)
	Webhook(event *stripe.Event) error
}

type ReservationService struct {
	DB *gorm.DB
}

func NewReservationService(DB *gorm.DB) ReservationServiceInterface {
	return &ReservationService{DB: DB}
}

func (service *ReservationService) CreateReservation(userID int, roomPrice decimal.Decimal, roomID int, reservationRequest request.ReservationRequest) (resources.ReservationWithPayemtResource, error) {
	service.DB.Begin()

	//check available room
	availableRoom := models.Reservation{}
	errAvailableRoom := service.DB.Joins("RoomNumber").Where("room_number_id = ? AND DATE(check_in) = ? AND room_id = ? AND ISNULL(check_out)", reservationRequest.RoomNumberID, reservationRequest.ParsedCheckin.Format("2006-01-02"), roomID).First(&availableRoom).Error
	if errAvailableRoom == nil {
		service.DB.Rollback()
		return resources.ReservationWithPayemtResource{}, errors.New("unavailable")
	}

	// caalculate total roomPrice * days
	total := roomPrice.Mul(decimal.NewFromInt(int64(reservationRequest.Days)))

	reservation := models.Reservation{
		UserID:       uint(userID),
		RoomNumberID: reservationRequest.RoomNumberID,
		CheckIn:      reservationRequest.ParsedCheckin,
		Days:         reservationRequest.Days,
		Total:        total,
		Status:       "unpaid",
	}
	errReservation := service.DB.Create(&reservation).Error
	if errReservation != nil {
		service.DB.Rollback()
		log.Error("Error create reservation : ", errReservation.Error())
		return resources.ReservationWithPayemtResource{}, errReservation
	}

	//todo create payment
	stripe.Key = configs.ENV.STRIPE_KEY
	paymentParams := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(total.Mul(decimal.NewFromInt(100)).IntPart())),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"reservation_id": string(reservation.ID),
		},
	}
	payment, errPayment := paymentintent.New(paymentParams)
	if errPayment != nil {
		service.DB.Rollback()
		log.Error("Error create payment : ", errPayment.Error())
		return resources.ReservationWithPayemtResource{}, errPayment
	}

	errGetReservation := service.DB.Where("reservations.id = ?", reservation.ID).Joins("RoomNumber").Joins("RoomNumber.Room").Take(&reservation).Error
	if errGetReservation != nil {
		service.DB.Rollback()
		log.Error("Error get reservation : ", errGetReservation.Error())
		return resources.ReservationWithPayemtResource{}, errGetReservation
	}

	reservationResource := resources.ReservationWithPayemtResource{
		ID:     reservation.ID,
		UserID: reservation.UserID,
		RoomNumber: resources.RoomNumberResource{
			ID:       reservation.RoomNumber.ID,
			RoomID:   reservation.RoomNumber.RoomID,
			Number:   reservation.RoomNumber.Number,
			Features: reservation.RoomNumber.Features,
			Price:    reservation.RoomNumber.Room.Price,
		},
		CheckIn:                reservation.CheckIn,
		CheckOut:               reservation.CheckOut,
		Total:                  reservation.Total,
		Status:                 reservation.Status,
		Days:                   reservation.Days,
		CreatedAt:              reservation.CreatedAt,
		StripePaymentClientKey: payment.ClientSecret,
	}

	service.DB.Commit()
	return reservationResource, nil
}

func (service *ReservationService) GetReservationByID(userID int, reservationID int) (resources.ReservationResource, error) {
	reservation := models.Reservation{}
	errReservation := service.DB.Where("reservations.id = ? AND user_id = ?", reservationID, userID).Joins("RoomNumber").First(&reservation).Error
	if errors.Is(errReservation, gorm.ErrRecordNotFound) {
		return resources.ReservationResource{}, gorm.ErrRecordNotFound
	}
	if errReservation != nil {
		log.Error("Error get reservation by ID : ", errReservation.Error())
		return resources.ReservationResource{}, errReservation
	}

	errGetReservation := service.DB.Where("reservations.id = ?", reservation.ID).Joins("RoomNumber").Joins("RoomNumber.Room").Take(&reservation).Error
	if errGetReservation != nil {
		log.Error("Error get reservation : ", errGetReservation.Error())
		return resources.ReservationResource{}, errGetReservation
	}

	reservationResource := resources.ReservationResource{
		ID:     reservation.ID,
		UserID: reservation.UserID,
		RoomNumber: resources.RoomNumberResource{
			ID:       reservation.RoomNumber.ID,
			RoomID:   reservation.RoomNumber.RoomID,
			Number:   reservation.RoomNumber.Number,
			Features: reservation.RoomNumber.Features,
			Price:    reservation.RoomNumber.Room.Price,
		},
		CheckIn:   reservation.CheckIn,
		CheckOut:  reservation.CheckOut,
		Total:     reservation.Total,
		Status:    reservation.Status,
		Days:      reservation.Days,
		CreatedAt: reservation.CreatedAt,
	}

	return reservationResource, nil
}

func (service *ReservationService) GetUserReservation(userID int) ([]resources.ReservationResource, error) {
	var reservations []models.Reservation
	errReservations := service.DB.Where("user_id = ?", userID).Joins("RoomNumber").Joins("RoomNumber.Room").Find(&reservations).Error
	if errReservations != nil {
		log.Error("Error get user reservations : ", errReservations.Error())
		return []resources.ReservationResource{}, errReservations
	}

	var reservationResources []resources.ReservationResource
	for _, reservation := range reservations {
		reservationResource := resources.ReservationResource{
			ID:     reservation.ID,
			UserID: reservation.UserID,
			RoomNumber: resources.RoomNumberResource{
				ID:       reservation.RoomNumber.ID,
				RoomID:   reservation.RoomNumber.RoomID,
				Number:   reservation.RoomNumber.Number,
				Features: reservation.RoomNumber.Features,
				Price:    reservation.RoomNumber.Room.Price,
			},
			CheckIn:   reservation.CheckIn,
			CheckOut:  reservation.CheckOut,
			Total:     reservation.Total,
			Status:    reservation.Status,
			Days:      reservation.Days,
			CreatedAt: reservation.CreatedAt,
		}
		reservationResources = append(reservationResources, reservationResource)
	}

	return reservationResources, nil
}

func (service *ReservationService) Webhook(event *stripe.Event) error {
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		errUnmarshal := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if errUnmarshal != nil {
			log.Error("Error unmarshal payment intent : ", errUnmarshal.Error())
			return errUnmarshal
		}

		metadata := paymentIntent.Metadata
		reservationID := metadata["reservation_id"]
		timeNow := time.Now()
		reservation := models.Reservation{
			Status:           "paid",
			PaymentSuccessAt: &timeNow,
		}
		errReservation := service.DB.Where("id = ?", reservationID).Updates(&reservation).Error
		if errReservation != nil {
			log.Error("Error update reservation status : ", errReservation.Error())
			return errReservation
		}
	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		errUnmarshal := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if errUnmarshal != nil {
			log.Error("Error unmarshal payment intent : ", errUnmarshal.Error())
			return errUnmarshal
		}

		metadata := paymentIntent.Metadata
		reservationID := metadata["reservation_id"]
		timeNow := time.Now()
		reservation := models.Reservation{
			Status:          "failed",
			PaymentFailedAt: &timeNow,
		}
		errReservation := service.DB.Where("id = ?", reservationID).Updates(&reservation).Error
		if errReservation != nil {
			log.Error("Error update reservation status : ", errReservation.Error())
			return errReservation
		}
	default:
		log.Error("Unhandled event type : ", event.Type)
	}

	return nil
}
