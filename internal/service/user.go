package service

import (
	"avito-trainee/internal/dto/user"
	"avito-trainee/internal/models"
	"avito-trainee/internal/repository"
)

type User interface {
	Create(dto user.CreateUserDTO) (*models.User, error)
	Delete(id uint) error
	Get(id uint) (models.User, error)
}

type UserService struct {
	UserRepository    repository.UserRepository
	SegmentRepository repository.SegmentRepository
}

func NewUserService(userRep repository.UserRepository, segmRep repository.SegmentRepository) *UserService {
	return &UserService{
		UserRepository:    userRep,
		SegmentRepository: segmRep,
	}
}

func (s *UserService) Create(dto user.CreateUserDTO) (*models.User, error) {
	segments, err := s.SegmentRepository.GetListByNames(dto.Segments)
	if err != nil {
		return nil, err
	}

	var segmentsIds []uint

	for _, v := range segments {
		segmentsIds = append(segmentsIds, v.ID)
	}

	user, err := s.UserRepository.Create(&models.User{}, segmentsIds)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(id uint) error {
	err := s.UserRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Get(id uint) (models.User, error) {
	user, err := s.UserRepository.Get(id)
	if err != nil {
		return user, err
	}

	return user, nil
}
