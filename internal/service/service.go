package service

import (
	"avito-trainee/internal/repository"
	"avito-trainee/pkg/logger"

	"github.com/robfig/cron/v3"
)

type Service struct {
	Repository *repository.Repository
	Logger     *logger.LogManager

	CronService *CronService

	Segment Segment
	User    User
}

type Deps struct {
	Repository   *repository.Repository
	Logger       *logger.LogManager
	CronSheduler *cron.Cron
}

func NewService(deps *Deps) *Service {
	segment := NewSegmentService(deps.Repository.Segment)
	user := NewUserService(deps.Repository.User, deps.Repository.Segment)

	cronService := NewCronService(deps.CronSheduler, deps.Logger, segment)

	return &Service{
		Repository:  deps.Repository,
		Logger:      deps.Logger,
		Segment:     segment,
		User:        user,
		CronService: cronService,
	}
}
