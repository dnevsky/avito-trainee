package repository

import (
	"avito-trainee/internal/dto/segment"
	"avito-trainee/internal/models"
	"time"
)

type SegmentRepository interface {
	Create(*models.Segment, *float64) (*models.Segment, error)
	DeleteByName(name string) error
	GetByName(name string) (models.Segment, error)

	GetListByNames(names []string) ([]models.Segment, error)

	AttachUser(userId uint, segmentIds []uint, ttl *time.Time) (attachedIds []uint, err error)
	DeleteUser(userId uint, segmentIds []uint) (attachedIds []uint, err error)

	DeleteTTLSegmentExpired() error
	DeleteTTLUserSegmentExpired() error

	DataRowsHistory(dto segment.DownloadHistoryDTO) ([][]string, error)
}
