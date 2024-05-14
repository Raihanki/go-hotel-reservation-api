package services

import (
	"errors"
	log "github.com/sirupsen/logrus"

	"github.com/Raihanki/go-hotel-reservation-api/models"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"gorm.io/gorm"
)

type RoomServiceInterface interface {
	CreateRoom(roomRequest request.RoomRequest) (resources.RoomResource, error)
	GetRoomByID(roomID int) (resources.RoomResource, error)
	UpdateRoomByID(roomRequest request.RoomRequest, roomID int) error
	DeleteRoomByID(roomID int) (int, error)
	GetAllRooms(hotelID int) ([]resources.RoomResource, error)
}

type RoomService struct {
	DB *gorm.DB
}

func NewRoomService(DB *gorm.DB) RoomServiceInterface {
	return &RoomService{DB: DB}
}

func (service *RoomService) CreateRoom(roomRequest request.RoomRequest) (resources.RoomResource, error) {
	room := models.Room{
		HotelID:     roomRequest.HotelID,
		Name:        roomRequest.Name,
		Price:       roomRequest.Price,
		MaxPeople:   roomRequest.MaxPeople,
		Description: roomRequest.Description,
	}
	errRoom := service.DB.Create(&room).Error
	if errRoom != nil {
		log.Error("Error create room : ", errRoom.Error())
		return resources.RoomResource{}, errRoom
	}

	getRoomErr := service.DB.Where("rooms.id = ?", room.ID).Joins("Hotel").Take(&room).Error
	if getRoomErr != nil {
		return resources.RoomResource{}, getRoomErr
	}

	roomResource := resources.RoomResource{
		ID: room.ID,
		Hotel: resources.HotelMiniResource{
			ID:      room.Hotel.ID,
			Name:    room.Hotel.Name,
			Rating:  room.Hotel.Rating,
			Address: room.Hotel.Address,
		},
		Name:        room.Name,
		Price:       room.Price,
		MaxPeople:   room.MaxPeople,
		Description: room.Description,
		CreatedAt:   room.CreatedAt,
	}
	return roomResource, nil
}

func (service *RoomService) GetRoomByID(roomID int) (resources.RoomResource, error) {
	room := models.Room{}
	errRoom := service.DB.Where("rooms.id = ?", roomID).Joins("Hotel").First(&room).Error
	if errors.Is(errRoom, gorm.ErrRecordNotFound) {
		return resources.RoomResource{}, gorm.ErrRecordNotFound
	}
	if errRoom != nil {
		log.Error("Error get room by ID : ", errRoom.Error())
		return resources.RoomResource{}, errRoom
	}

	roomResource := resources.RoomResource{
		ID: room.ID,
		Hotel: resources.HotelMiniResource{
			ID:      room.Hotel.ID,
			Name:    room.Hotel.Name,
			Rating:  room.Hotel.Rating,
			Address: room.Hotel.Address,
		},
		Name:        room.Name,
		Price:       room.Price,
		MaxPeople:   room.MaxPeople,
		Description: room.Description,
		CreatedAt:   room.CreatedAt,
	}
	return roomResource, nil
}

func (service *RoomService) UpdateRoomByID(roomRequest request.RoomRequest, roomID int) error {
	room := models.Room{
		HotelID:     roomRequest.HotelID,
		Name:        roomRequest.Name,
		Price:       roomRequest.Price,
		MaxPeople:   roomRequest.MaxPeople,
		Description: roomRequest.Description,
	}
	updateRoom := service.DB.Where("id = ?", roomID).Updates(&room).Error
	if updateRoom != nil {
		return updateRoom
	}

	return nil
}

func (service *RoomService) DeleteRoomByID(roomID int) (int, error) {
	room := models.Room{}
	deleteRoom := service.DB.Where("id = ?", roomID).Delete(&room)
	if deleteRoom.Error != nil {
		log.Error("Error get room by ID : ", deleteRoom.Error.Error())
		return 0, deleteRoom.Error
	}

	return int(deleteRoom.RowsAffected), nil
}

func (service *RoomService) GetAllRooms(hotelID int) ([]resources.RoomResource, error) {
	var rooms []models.Room
	errRooms := service.DB.Where("hotel_id = ?", hotelID).Preload("Hotel").Find(&rooms).Error
	if errRooms != nil {
		log.Error("Error get all rooms : ", errRooms.Error())
		return []resources.RoomResource{}, errRooms
	}

	var roomResources []resources.RoomResource
	for _, room := range rooms {
		roomResources = append(roomResources, resources.RoomResource{
			ID: room.ID,
			Hotel: resources.HotelMiniResource{
				ID:      room.Hotel.ID,
				Name:    room.Hotel.Name,
				Rating:  room.Hotel.Rating,
				Address: room.Hotel.Address,
			},
			Name:        room.Name,
			Price:       room.Price,
			MaxPeople:   room.MaxPeople,
			Description: room.Description,
			CreatedAt:   room.CreatedAt,
		})
	}

	return roomResources, nil
}
