package services

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/Raihanki/go-hotel-reservation-api/models"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ReservationServiceInterface interface {
	CreateReservation(userID int, price decimal.Decimal, roomID int, reservationRequest request.ReservationRequest) (resources.ReservationResource, error)
	GetReservationByID(userID int, reservationID int) (resources.ReservationResource, error)
	GetUserReservation(userID int) ([]resources.ReservationResource, error)
}

type ReservationService struct {
	DB *gorm.DB
}

func NewReservationService(DB *gorm.DB) ReservationServiceInterface {
	return &ReservationService{DB: DB}
}

func (service *ReservationService) CreateReservation(userID int, roomPrice decimal.Decimal, roomID int, reservationRequest request.ReservationRequest) (resources.ReservationResource, error) {
	//check available room
	availableRoom := models.Reservation{}
	errAvailableRoom := service.DB.Joins("RoomNumber").Where("room_number_id = ? AND DATE(check_in) = ? AND room_id = ? AND ISNULL(check_out)", reservationRequest.RoomNumberID, reservationRequest.ParsedCheckin.Format("2006-01-02"), roomID).First(&availableRoom).Error
	if errAvailableRoom == nil {
		return resources.ReservationResource{}, errors.New("unavailable")
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
		log.Error("Error create reservation : ", errReservation.Error())
		return resources.ReservationResource{}, errReservation
	}

	//todo create payment

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
