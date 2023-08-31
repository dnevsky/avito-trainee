package service

import (
	"avito-trainee/internal/dto/segment"
	"avito-trainee/internal/models"
	"avito-trainee/internal/repository"
	"avito-trainee/pkg/logger/helpers"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Segment interface {
	Create(dto segment.CreateSegmentDTO) (*models.Segment, error)
	Delete(dto segment.DeleteSegmentDTO) error
	Get(name string) (models.Segment, error)
	AttachSegmentsToUser(dto segment.AttachSegmentsDTO) ([]string, []string, error)
	CheckTTLSegment() error
	CheckTTLUserSegment() error
	DownloadHistory(dto segment.DownloadHistoryDTO) (*os.File, error)
}

type SegmentService struct {
	SegmentRepository repository.SegmentRepository
}

func NewSegmentService(segmentRep repository.SegmentRepository) *SegmentService {
	return &SegmentService{
		SegmentRepository: segmentRep,
	}
}

func (s *SegmentService) CheckTTLSegment() error {
	return s.SegmentRepository.DeleteTTLSegmentExpired()
}

func (s *SegmentService) CheckTTLUserSegment() error {
	return s.SegmentRepository.DeleteTTLUserSegmentExpired()
}

func (s *SegmentService) DownloadHistory(dto segment.DownloadHistoryDTO) (*os.File, error) {
	fileData, err := s.SegmentRepository.DataRowsHistory(dto)
	if err != nil {
		return nil, err
	}

	return helpers.DataToCsv(fileData)
}

func (s *SegmentService) AttachSegmentsToUser(dto segment.AttachSegmentsDTO) ([]string, []string, error) {
	segmentsToDeleteObj, err := s.SegmentRepository.GetListByNames(dto.DeleteSegments)
	if err != nil {
		return nil, nil, err
	}

	segmentsToAddObj, err := s.SegmentRepository.GetListByNames(dto.AddSegments)
	if err != nil {
		return nil, nil, err
	}

	segmentsToAddIds := make([]uint, 0)
	segmentsToDeleteIds := make([]uint, 0)

	for _, v := range segmentsToAddObj {
		segmentsToAddIds = append(segmentsToAddIds, v.ID)
	}

	for _, v := range segmentsToDeleteObj {
		segmentsToDeleteIds = append(segmentsToDeleteIds, v.ID)
	}

	var attachAddUint []uint
	var attachDelUint []uint

	if len(segmentsToDeleteIds) > 0 {
		attachDelUint, err = s.SegmentRepository.DeleteUser(dto.UserID, segmentsToDeleteIds)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(segmentsToAddIds) > 0 {
		var ttl *time.Time
		if dto.TTL != "" {
			ttlParse, err := helpers.ParseStringToTime(dto.TTL)
			if err != nil {
				return nil, nil, err
			}

			ttl = &ttlParse
		}

		attachAddUint, err = s.SegmentRepository.AttachUser(dto.UserID, segmentsToAddIds, ttl)
		if err != nil {
			return nil, nil, err
		}

	}

	attachedAddNames := []string{}
	attachedDelNames := []string{}

	for _, v := range segmentsToAddObj {
		for _, vv := range attachAddUint {
			if v.ID == vv {
				attachedAddNames = append(attachedAddNames, v.Name)
			}
		}
	}

	for _, v := range segmentsToDeleteObj {
		for _, vv := range attachDelUint {
			if v.ID == vv {
				attachedDelNames = append(attachedDelNames, v.Name)
			}
		}
	}

	return attachedAddNames, attachedDelNames, nil
}

func (s *SegmentService) Create(dto segment.CreateSegmentDTO) (*models.Segment, error) {
	segment := &models.Segment{
		Name: dto.Name,
	}

	if dto.TTL != nil {
		ttl, err := helpers.ParseStringToTime(*dto.TTL)
		if err != nil {
			return nil, err
		}
		segment.TTL = &ttl
	}

	var percent *float64
	if dto.AutoAttachPercent != nil {
		*percent = float64(*dto.AutoAttachPercent) / 100
	}

	segment, err := s.SegmentRepository.Create(segment, percent)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, err
	}

	return segment, nil
}

func (s *SegmentService) Delete(dto segment.DeleteSegmentDTO) error {
	err := s.SegmentRepository.DeleteByName(dto.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *SegmentService) Get(name string) (models.Segment, error) {
	segment, err := s.SegmentRepository.GetByName(name)
	if err != nil {
		return segment, err
	}

	return segment, nil
}
