package model

// APIResponse là khuôn mẫu chuẩn để trả kết quả về cho Postman
type APIResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Tự động ẩn đi nếu không có dữ liệu
}