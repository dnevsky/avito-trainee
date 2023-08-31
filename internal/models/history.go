package models

import "time"

type OperationType string

var (
	SegmentAddOperationType    OperationType = "add"
	SegmentDeleteOperationType OperationType = "delete"
)

type SegmentHistory struct {
	BaseModel
	UserID      uint          `json:"user_id"`
	SegmentName string        `json:"segment_name"`
	Operation   OperationType `json:"operation"`
	DateTime    time.Time     `json:"date_time"`
}
