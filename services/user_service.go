package services

import (
	"errors"

	"github.com/Raihanki/go-hotel-reservation-api/helpers"
	"github.com/Raihanki/go-hotel-reservation-api/models"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	Login(loginRequest request.LoginRequest) (token string, err error)
	Register(registerRequest request.RegisterRequest) (token string, err error)
	GetUserByEmail(email string) (models.User, error)
	GetUserById(userID int) (resources.UserResource, error)
}

type UserService struct {
	DB *gorm.DB
}

func NewUserService(DB *gorm.DB) UserServiceInterface {
	return &UserService{DB: DB}
}

func (service *UserService) Login(loginRequest request.LoginRequest) (token string, err error) {
	user, errGetUser := service.GetUserByEmail(loginRequest.Email)
	if errGetUser != nil {
		return "", errGetUser
	}

	errVerifyHash := helpers.VerifyHash(user.Password, loginRequest.Password)
	if errVerifyHash != nil {
		return "", errVerifyHash
	}

	token, errToken := helpers.GenerateJWT(int(user.ID))
	if errToken != nil {
		return token, errToken
	}

	return token, nil
}

func (service *UserService) Register(registerRequest request.RegisterRequest) (token string, err error) {
	hashedPassword, errHash := helpers.Hash(registerRequest.Password)
	if errHash != nil {
		return "", errHash
	}

	user := models.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: hashedPassword,
		IsAdmin:  false,
	}
	errCreateUser := service.DB.Create(&user).Error
	if errCreateUser != nil {
		return "", errCreateUser
	}

	token, errToken := helpers.GenerateJWT(int(user.ID))
	if errToken != nil {
		return token, errToken
	}

	return token, nil
}

func (service *UserService) GetUserById(userID int) (resources.UserResource, error) {
	user := models.User{}
	getUser := service.DB.Where("id = ?", userID).Take(&user)
	if errors.Is(getUser.Error, gorm.ErrRecordNotFound) {
		return resources.UserResource{}, gorm.ErrRecordNotFound
	}
	if getUser.Error != nil {
		return resources.UserResource{}, getUser.Error
	}

	userResource := resources.UserResource{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	return userResource, nil
}

func (service *UserService) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}
	getUser := service.DB.Where("email = ?", email).Take(&user)
	if errors.Is(getUser.Error, gorm.ErrRecordNotFound) {
		return models.User{}, gorm.ErrRecordNotFound
	}
	if getUser.Error != nil {
		return models.User{}, getUser.Error
	}

	return user, nil
}
