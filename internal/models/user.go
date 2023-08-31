package models

import "time"

type User struct {
	BaseModel
	Segments []*Segment `json:"segments" gorm:"many2many:user_segments"`
}

type UserSegment struct {
	UserID    uint      `json:"user_id" gorm:"primaryKey"`
	SegmentID uint      `json:"segment_id" gorm:"primaryKey"`
	TTL       time.Time `json:"ttl,omitempty" gorm:"default:null;index"`
}

func (UserSegment) TableName() string {
	return "user_segments"
}
