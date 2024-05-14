package services

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/Raihanki/go-hotel-reservation-api/models"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"gorm.io/gorm"
)

type RoomNumberServiceInterface interface {
	CreateRoomNumber(roomNumberRequest request.RoomNumberRequest) (resources.RoomNumberResource, error)
	GetRoomNumberByID(roomNumberID int) (resources.RoomNumberWithRoomResource, error)
	DeleteRoomNumberByID(roomNumberID int) (int, error)
	GetAllRoomNumbers(roomID int) ([]resources.RoomNumberResource, error)
}

type RoomNumberService struct {
	DB *gorm.DB
}

func NewRoomNumberService(DB *gorm.DB) RoomNumberServiceInterface {
	return &RoomNumberService{DB: DB}
}

func (service *RoomNumberService) CreateRoomNumber(roomNumberRequest request.RoomNumberRequest) (resources.RoomNumberResource, error) {
	roomNumber := models.RoomNumber{
		RoomID:   roomNumberRequest.RoomID,
		Number:   roomNumberRequest.Number,
		Features: roomNumberRequest.Features,
	}
	errRoomNumber := service.DB.Create(&roomNumber).Error
	if errRoomNumber != nil {
		log.Error("Error create room number : ", errRoomNumber.Error())
		return resources.RoomNumberResource{}, errRoomNumber
	}

	roomResource := resources.RoomNumberResource{
		ID:        roomNumber.ID,
		RoomID:    roomNumber.RoomID,
		Number:    roomNumber.Number,
		Features:  roomNumber.Features,
		CreatedAt: roomNumber.CreatedAt,
	}
	return roomResource, nil
}

func (service *RoomNumberService) GetRoomNumberByID(roomNumberID int) (resources.RoomNumberWithRoomResource, error) {
	roomNumber := models.RoomNumber{}
	errRoomNumber := service.DB.Where("room_numbers.id = ?", roomNumberID).Joins("Room").Joins("Room.Hotel").First(&roomNumber).Error
	if errors.Is(errRoomNumber, gorm.ErrRecordNotFound) {
		return resources.RoomNumberWithRoomResource{}, gorm.ErrRecordNotFound
	}
	if errRoomNumber != nil {
		log.Error("Error get room number by ID : ", errRoomNumber.Error())
		return resources.RoomNumberWithRoomResource{}, errRoomNumber
	}

	roomNumberResource := resources.RoomNumberWithRoomResource{
		ID:     roomNumber.ID,
		RoomID: roomNumber.RoomID,
		Room: resources.RoomResource{
			ID: roomNumber.Room.ID,
			Hotel: resources.HotelMiniResource{
				ID:      roomNumber.Room.Hotel.ID,
				Name:    roomNumber.Room.Hotel.Name,
				Rating:  roomNumber.Room.Hotel.Rating,
				Address: roomNumber.Room.Hotel.Address,
			},
			Name:        roomNumber.Room.Name,
			Price:       roomNumber.Room.Price,
			MaxPeople:   roomNumber.Room.MaxPeople,
			Description: roomNumber.Room.Description,
			Photo:       roomNumber.Room.Photo,
			CreatedAt:   roomNumber.Room.CreatedAt,
		},
		Number:    roomNumber.Number,
		Features:  roomNumber.Features,
		CreatedAt: roomNumber.CreatedAt,
	}
	return roomNumberResource, nil
}

func (service *RoomNumberService) DeleteRoomNumberByID(roomNumberID int) (int, error) {
	roomNumber := models.RoomNumber{}
	deleteRoomNumber := service.DB.Where("id = ?", roomNumberID).Delete(&roomNumber)
	if deleteRoomNumber.RowsAffected == 0 {
		return 0, deleteRoomNumber.Error
	}
	if deleteRoomNumber.Error != nil {
		log.Error("Error delete room number by ID : ", deleteRoomNumber.Error.Error())
		return 0, deleteRoomNumber.Error
	}

	return 1, nil
}

func (service *RoomNumberService) GetAllRoomNumbers(roomID int) ([]resources.RoomNumberResource, error) {
	roomNumbers := []models.RoomNumber{}
	errRoomNumbers := service.DB.Where("room_id = ?", roomID).Find(&roomNumbers).Error
	if errRoomNumbers != nil {
		log.Error("Error get room numbers : ", errRoomNumbers.Error())
		return []resources.RoomNumberResource{}, errRoomNumbers
	}

	roomNumberResources := []resources.RoomNumberResource{}
	for _, roomNumber := range roomNumbers {
		roomNumberResource := resources.RoomNumberResource{
			ID:        roomNumber.ID,
			RoomID:    roomNumber.RoomID,
			Number:    roomNumber.Number,
			Features:  roomNumber.Features,
			CreatedAt: roomNumber.CreatedAt,
		}
		roomNumberResources = append(roomNumberResources, roomNumberResource)
	}
	return roomNumberResources, nil
}
