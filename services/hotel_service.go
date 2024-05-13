package services

import (
	"errors"

	"github.com/Raihanki/go-hotel-reservation-api/models"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"gorm.io/gorm"
)

type HotelServiceInterface interface {
	CreateHotel(hotelRequest request.HotelRequest) (resources.HotelResource, error)
	GetHotelByID(hotelID int) (resources.HotelResource, error)
	UpdateHotelByID(hotelRequest request.HotelRequest, hotelID int) error
	DeleleHotelByID(hotelID int) (int, error)
	GetAllHotels() ([]resources.HotelResource, error)
}

type HotelService struct {
	DB *gorm.DB
}

func NewHotelService(DB *gorm.DB) HotelServiceInterface {
	return &HotelService{DB: DB}
}

func (service *HotelService) CreateHotel(hotelRequest request.HotelRequest) (resources.HotelResource, error) {
	hotel := models.Hotel{
		Name:        hotelRequest.Name,
		Type:        hotelRequest.Type,
		City:        hotelRequest.City,
		Address:     hotelRequest.Address,
		Photo:       hotelRequest.Photo,
		Description: hotelRequest.Description,
	}
	errHotel := service.DB.Create(&hotel).Error
	if errHotel != nil {
		return resources.HotelResource{}, errHotel
	}

	hotelResource := resources.HotelResource{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Type:        hotel.Type,
		City:        hotel.City,
		Address:     hotel.Address,
		Photo:       hotel.Photo,
		Description: hotel.Description,
		CreatedAt:   hotel.CreatedAt,
	}
	return hotelResource, nil
}

func (service *HotelService) GetHotelByID(hotelID int) (resources.HotelResource, error) {
	hotel := models.Hotel{}
	errHotel := service.DB.Where("id = ?", hotelID).First(&hotel).Error
	if errors.Is(errHotel, gorm.ErrRecordNotFound) {
		return resources.HotelResource{}, gorm.ErrRecordNotFound
	}
	if errHotel != nil {
		return resources.HotelResource{}, errHotel
	}

	hotelResource := resources.HotelResource{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Type:        hotel.Type,
		City:        hotel.City,
		Address:     hotel.Address,
		Photo:       hotel.Photo,
		Description: hotel.Description,
		Rating:      hotel.Rating,
		CreatedAt:   hotel.CreatedAt,
	}

	return hotelResource, nil
}

func (service *HotelService) UpdateHotelByID(hotelRequest request.HotelRequest, hotelID int) error {
	hotel := models.Hotel{
		Name:        hotelRequest.Name,
		Type:        hotelRequest.Type,
		City:        hotelRequest.City,
		Address:     hotelRequest.Address,
		Photo:       hotelRequest.Photo,
		Description: hotelRequest.Description,
		Rating:      hotelRequest.Rating,
	}
	errHotel := service.DB.Where("id = ?", hotelID).Updates(&hotel).Error
	if errHotel != nil {
		return errHotel
	}

	return nil
}

func (service *HotelService) DeleleHotelByID(hotelID int) (int, error) {
	hotel := models.Hotel{}
	deleteHotel := service.DB.Where("id = ?", hotelID).Delete(&hotel)
	if deleteHotel.Error != nil {
		return 0, deleteHotel.Error
	}

	return int(deleteHotel.RowsAffected), nil
}

func (service *HotelService) GetAllHotels() ([]resources.HotelResource, error) {
	var hotels []models.Hotel
	errHotels := service.DB.Find(&hotels).Error
	if errHotels != nil {
		return nil, errHotels
	}

	var hotelResources []resources.HotelResource
	for _, hotel := range hotels {
		hotelResource := resources.HotelResource{
			ID:          hotel.ID,
			Name:        hotel.Name,
			Type:        hotel.Type,
			City:        hotel.City,
			Address:     hotel.Address,
			Photo:       hotel.Photo,
			Description: hotel.Description,
			Rating:      hotel.Rating,
			CreatedAt:   hotel.CreatedAt,
		}
		hotelResources = append(hotelResources, hotelResource)
	}

	return hotelResources, nil
}
