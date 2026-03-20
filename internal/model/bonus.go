package model

// Bài 6
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Data       []*Asset   `json:"data"`
	Pagination Pagination `json:"pagination"`
}