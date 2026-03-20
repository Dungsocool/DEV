package model

import (
	"fmt"
	"time"
)

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

// Validate kiểm tra xem dữ liệu của Asset có hợp lệ không
func (a *Asset) Validate() error {
	// 1. Kiểm tra tên không được để trống
	if a.Name == "" {
		return fmt.Errorf("tên tài sản (Name) không được để trống")
	}

	// 2. Kiểm tra loại (Type) phải nằm trong danh sách cho phép
	switch a.Type {
	case TypeDomain, TypeIP, TypeService:
		// Hợp lệ
		return nil
	default:
		return fmt.Errorf("loại tài sản '%s' không hợp lệ (chỉ chấp nhận: domain, ip, service)", a.Type)
	}
}
