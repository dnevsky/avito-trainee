package repository

import (
	"avito-trainee/internal/repository/postgresDB"

	"gorm.io/gorm"
)

type Repository struct {
	Segment SegmentRepository
	User    UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Segment: postgresDB.NewSegmentRepository(db),
		User:    postgresDB.NewUserRepository(db, postgresDB.NewSegmentRepository(db)),
	}
}
