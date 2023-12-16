package services

import (
	"crud-gin/models"
)

// API контракты
type UserService interface {
	CreateUser(user *models.User) error
	GetUser(name string) (*models.User, error)
	GetAll() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(name string) error
}
