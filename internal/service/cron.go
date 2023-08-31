package service

import (
	"avito-trainee/pkg/logger"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	Scheduler *cron.Cron
	Logger    *logger.LogManager
	Segment   Segment
}

func NewCronService(
	sh *cron.Cron,
	logger *logger.LogManager,
	segment Segment,
) *CronService {
	return &CronService{
		Scheduler: sh,
		Logger:    logger,
		Segment:   segment,
	}
}

func (s *CronService) Start() {
	s.Scheduler.Start()
}

func (s *CronService) Init() {
	_, err := s.Scheduler.AddFunc("*/1 * * * *", func() {
		err := s.Segment.CheckTTLSegment()
		if err != nil {
			s.Logger.Error(err)
		}
	})
	if err != nil {
		s.Logger.Error(err)
		return
	}

	_, err = s.Scheduler.AddFunc("*/1 * * * *", func() {
		err := s.Segment.CheckTTLUserSegment()
		if err != nil {
			s.Logger.Error(err)
		}
	})
	if err != nil {
		s.Logger.Error(err)
		return
	}

}

func (s *CronService) SetJob(notificationCron string) (cron.EntryID, error) {
	return 0, nil
}

func (s *CronService) RemoveJob() {

}
