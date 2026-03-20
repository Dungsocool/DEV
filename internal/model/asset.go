package model

import "time"

// Asset đại diện cho một tài sản trong hệ thống 
type Asset struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Định nghĩa sẵn các loại Type hợp lệ
const (
	TypeDomain  = "domain"
	TypeIP      = "ip"
	TypeService = "service"
)

// Định nghĩa sẵn các Trạng thái (Status)
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)