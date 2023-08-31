package segment

type AttachSegmentsDTO struct {
	UserID         uint     `json:"user_id" form:"user_id"`
	AddSegments    []string `json:"add_segments" form:"add_segments"`
	DeleteSegments []string `json:"delete_segments" form:"delete_segments"`
	TTL            string   `json:"ttl" form:"ttl"`
}
