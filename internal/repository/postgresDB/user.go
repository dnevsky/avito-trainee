package postgresDB

import (
	"avito-trainee/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db                *gorm.DB
	segmentRepository *SegmentRepository
}

func NewUserRepository(db *gorm.DB, segmRepo *SegmentRepository) *UserRepository {
	return &UserRepository{
		db:                db,
		segmentRepository: segmRepo,
	}
}

func (r *UserRepository) Create(user *models.User, segments []uint) (*models.User, error) {
	err := r.db.Clauses(clause.Returning{}).Create(&user).Error
	if err != nil {
		return nil, err
	}

	if segments != nil {
		_, err = r.segmentRepository.AttachUser(user.ID, segments, nil)
		if err != nil {
			return nil, err
		}
	}

	err = r.db.
		Preload("Segments").
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *UserRepository) Delete(id uint) error {
	res := r.db.Unscoped().Delete(&models.User{}, "id = ?", id)
	err := res.Error

	if err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepository) Get(id uint) (user models.User, err error) {
	err = r.db.Model(&models.User{}).
		Preload("Segments").
		First(&user, "id = ?", id).Error
	return user, err
}
