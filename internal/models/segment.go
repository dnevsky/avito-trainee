package models

import "time"

type Segment struct {
	BaseModel
	Name string     `json:"name" gorm:"uniqueIndex"`
	TTL  *time.Time `json:"ttl"`
}
