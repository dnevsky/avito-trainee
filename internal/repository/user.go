package repository

import (
	"avito-trainee/internal/models"
)

type UserRepository interface {
	Create(*models.User, []uint) (*models.User, error)
	Delete(id uint) error
	Get(id uint) (models.User, error)
}
