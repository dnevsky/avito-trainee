package segment

type DownloadHistoryDTO struct {
	Year   int  `json:"year" form:"year"`
	Month  int  `json:"month" form:"month"`
	UserID uint `json:"user_id" form:"user_id"`
}
