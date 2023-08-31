package postgresDB

import (
	"avito-trainee/internal/dto/segment"
	"avito-trainee/internal/models"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SegmentRepository struct {
	db *gorm.DB
}

func NewSegmentRepository(db *gorm.DB) *SegmentRepository {
	return &SegmentRepository{
		db: db,
	}
}

func (r *SegmentRepository) Create(segment *models.Segment, percent *float64) (*models.Segment, error) {
	db := r.db.Begin()
	err := db.Clauses(clause.Returning{}).Create(&segment).Error
	if err != nil {
		db.Rollback()
		return nil, err
	}

	if percent == nil {
		db.Commit()
		return segment, nil
	}

	var selectedUsers []models.User
	var attachList []models.UserSegment

	query := fmt.Sprintf("SELECT * FROM users ORDER BY random() LIMIT (SELECT COUNT(*) * %f FROM users)", *percent)

	err = db.Raw(query).Scan(&selectedUsers).Error
	if err != nil {
		db.Rollback()
		return nil, err
	}

	for _, v := range selectedUsers {
		attachList = append(attachList, models.UserSegment{
			UserID:    v.ID,
			SegmentID: segment.ID,
		})
	}

	err = db.Create(&attachList).Error
	if err != nil {
		db.Rollback()
		return nil, err
	}

	var createHistory []models.SegmentHistory

	for _, v := range attachList {
		createHistory = append(createHistory, models.SegmentHistory{
			UserID:      v.UserID,
			SegmentName: segment.Name,
			Operation:   models.SegmentAddOperationType,
			DateTime:    time.Now(),
		})
	}

	if createHistory != nil {
		err := db.Create(&createHistory).Error
		if err != nil {
			db.Rollback()
			return nil, err
		}
	}

	db.Commit()

	return segment, err
}

func (r *SegmentRepository) GetListByNames(names []string) ([]models.Segment, error) {
	var segments []models.Segment
	db := r.db.Model(&models.Segment{})
	db = db.Where("name IN (?)", names)

	err := db.Find(&segments).Error

	return segments, err
}

func (r *SegmentRepository) AttachUser(userId uint, segmentIds []uint, ttlArg *time.Time) ([]uint, error) {
	var userSegment []models.UserSegment
	var userSegmentAttached []models.UserSegment

	fmt.Println(ttlArg)
	for _, v := range segmentIds {
		userSegmentItem := models.UserSegment{
			UserID:    userId,
			SegmentID: v,
		}

		if ttlArg != nil {
			userSegmentItem.TTL = *ttlArg
		}
		userSegment = append(userSegment, userSegmentItem)
	}

	db := r.db.Begin()

	var userSegmentAttachedBefore []models.UserSegment

	res := db.Clauses(clause.OnConflict{DoNothing: true}).
		Find(&userSegmentAttachedBefore, "user_id = ? AND segment_id IN (?)", userId, segmentIds).
		Create(&userSegment).
		Find(&userSegmentAttached, "user_id = ? AND segment_id IN (?)", userId, segmentIds)

	if res.Error != nil && !strings.Contains(res.Error.Error(), "duplicate key value violates unique constraint") {
		db.Rollback()
		return nil, res.Error
	}

	existing := make(map[int]bool)
	addedIds := make([]uint, 0)

	for _, num := range userSegmentAttachedBefore {
		existing[int(num.SegmentID)] = true
	}

	for _, num := range userSegmentAttached {
		if _, exists := existing[int(num.SegmentID)]; !exists {
			addedIds = append(addedIds, uint(num.SegmentID))
			existing[int(num.SegmentID)] = true
		}
	}

	var createHistory []models.SegmentHistory

	for _, v := range addedIds {
		var segment models.Segment
		err := db.First(&segment, "id = ?", v).Error
		if err != nil {
			db.Rollback()
			return nil, err
		}

		createHistory = append(createHistory, models.SegmentHistory{
			UserID:      userId,
			SegmentName: segment.Name,
			Operation:   models.SegmentAddOperationType,
			DateTime:    time.Now(),
		})
	}

	if createHistory != nil {
		err := db.Create(&createHistory).Error
		if err != nil {
			db.Rollback()
			return nil, err
		}
	}

	db.Commit()

	return addedIds, nil
}

func (r *SegmentRepository) DeleteUser(userId uint, segmentIds []uint) ([]uint, error) {
	var userSegment []models.UserSegment
	var userSegmentDeleted []models.UserSegment

	for _, v := range segmentIds {
		userSegment = append(userSegment, models.UserSegment{
			UserID:    userId,
			SegmentID: v,
		})
	}

	db := r.db.Begin()

	res := db.Where("user_id = ? AND segment_id IN (?)", userId, segmentIds).
		Find(&userSegmentDeleted).
		Unscoped().Delete(&userSegment)

	if res.Error != nil && !strings.Contains(res.Error.Error(), "duplicate key value violates unique constraint") {
		db.Rollback()
		return nil, res.Error
	}

	deletedIds := make([]uint, 0)

	for _, num := range userSegmentDeleted {
		deletedIds = append(deletedIds, uint(num.SegmentID))
	}

	var createHistory []models.SegmentHistory

	for _, v := range deletedIds {
		var segment models.Segment
		err := db.First(&segment, "id = ?", v).Error
		if err != nil {
			db.Rollback()
			return nil, err
		}

		createHistory = append(createHistory, models.SegmentHistory{
			UserID:      userId,
			SegmentName: segment.Name,
			Operation:   models.SegmentDeleteOperationType,
			DateTime:    time.Now(),
		})
	}

	if createHistory != nil {
		err := db.Create(&createHistory).Error
		if err != nil {
			db.Rollback()
			return nil, err
		}
	}

	db.Commit()

	return deletedIds, nil
}

func (r *SegmentRepository) DeleteByName(name string) error {
	var segment models.Segment

	db := r.db.Begin()
	res := db.Where("name = ?", name).
		Find(&segment).
		Unscoped().Delete(&segment)
	err := res.Error

	if err != nil {
		db.Rollback()
		return err
	}

	if res.RowsAffected == 0 {
		db.Rollback()
		return gorm.ErrRecordNotFound
	}

	var userSegments []models.UserSegment

	err = db.Where("segment_id = ?", segment.ID).
		Find(&userSegments).
		Unscoped().Delete(&userSegments).Error
	if err != nil {
		db.Rollback()
		return err
	}

	var createHistory []models.SegmentHistory

	for _, v := range userSegments {
		createHistory = append(createHistory, models.SegmentHistory{
			UserID:      v.UserID,
			SegmentName: segment.Name,
			Operation:   models.SegmentDeleteOperationType,
			DateTime:    time.Now(),
		})
	}

	if createHistory != nil {
		err := db.Create(&createHistory).Error
		if err != nil {
			db.Rollback()
			return err
		}
	}

	db.Commit()

	return nil
}

func (r *SegmentRepository) GetByName(name string) (segment models.Segment, err error) {
	err = r.db.First(&segment, "name = ?", name).Error
	return segment, err
}

func (r *SegmentRepository) DeleteTTLSegmentExpired() error {
	var segments []models.Segment
	timeNow := time.Now().UTC()
	db := r.db.Begin()

	err := db.Find(&segments, "ttl < ?", timeNow).Error
	if err != nil {
		db.Rollback()
		return err
	}

	var userSegmentsDeleted []models.UserSegment
	segmentsIdsToDelete := make([]uint, 0)

	for _, v := range segments {
		segmentsIdsToDelete = append(segmentsIdsToDelete, v.ID)
	}

	err = db.Find(&userSegmentsDeleted, "segment_id IN (?)", segmentsIdsToDelete).
		Unscoped().Delete(&userSegmentsDeleted, "segment_id IN (?)", segmentsIdsToDelete).Error
	if err != nil {
		db.Rollback()
		return err
	}

	var createHistory []models.SegmentHistory

	for _, v := range userSegmentsDeleted {
		var segment models.Segment
		fmt.Println(1)
		err := db.First(&segment, "id = ?", v.SegmentID).Error
		if err != nil {
			fmt.Println(2)
			db.Rollback()
			return err
		}

		createHistory = append(createHistory, models.SegmentHistory{
			UserID:      v.UserID,
			SegmentName: segment.Name,
			Operation:   models.SegmentDeleteOperationType,
			DateTime:    time.Now(),
		})
	}

	if createHistory != nil {
		err := db.Create(&createHistory).Error
		if err != nil {
			db.Rollback()
			return err
		}
	}

	err = db.Unscoped().Delete(&segments, "ttl < ?", timeNow).Error
	if err != nil {
		db.Rollback()
		return err
	}

	db.Commit()

	return nil
}

func (r *SegmentRepository) DeleteTTLUserSegmentExpired() error {
	var userSegmentsDeleted []models.UserSegment
	timeNow := time.Now().UTC()
	db := r.db.Begin()

	err := db.Find(&userSegmentsDeleted, "ttl < ?", timeNow).
		Unscoped().Delete(&userSegmentsDeleted, "ttl < ?", timeNow).Error
	if err != nil {
		db.Rollback()
		return err
	}

	var createHistory []models.SegmentHistory

	for _, v := range userSegmentsDeleted {
		var segment models.Segment

		err := db.First(&segment, "id = ?", v.SegmentID).Error
		if err != nil {
			db.Rollback()
			return err
		}

		createHistory = append(createHistory, models.SegmentHistory{
			UserID:      v.UserID,
			SegmentName: segment.Name,
			Operation:   models.SegmentDeleteOperationType,
			DateTime:    time.Now(),
		})
	}

	if createHistory != nil {
		err := db.Create(&createHistory).Error
		if err != nil {
			db.Rollback()
			return err
		}
	}

	db.Commit()

	return nil
}

func (r *SegmentRepository) DataRowsHistory(dto segment.DownloadHistoryDTO) ([][]string, error) {
	var data [][]string

	head := []string{
		"ID User",
		"Segment Name",
		"Operation",
		"Date",
	}

	data = append(data, head)

	timeStart := time.Date(dto.Year, time.Month(dto.Month), 1, 0, 0, 0, 0, time.UTC)
	timeEnd := timeStart.AddDate(0, 1, 0)
	timeEnd = timeEnd.Add(-time.Second)

	var segmentsHistory []models.SegmentHistory

	err := r.db.Find(&segmentsHistory, "date_time > ? AND date_time < ? AND user_id = ?", timeStart, timeEnd, dto.UserID).Error
	if err != nil {
		return nil, err
	}

	for _, v := range segmentsHistory {
		data = append(data, []string{
			fmt.Sprintf("%d", v.UserID),
			v.SegmentName,
			string(v.Operation),
			v.DateTime.Local().String(),
		})
	}

	return data, nil
}
