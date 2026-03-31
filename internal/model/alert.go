package model

import "time"

// Alert types - các loại cảnh báo
const (
	AlertTypeHighRisk    = "high_risk"      // Phát hiện high-risk ports (22, 3389)
	AlertTypeExpiredCert = "expired_cert"   // Certificate sắp hết hạn
	AlertTypeNewService  = "new_service"    // Phát hiện service mới
	AlertTypeVulnerable  = "vulnerable"     // Phát hiện điểm yếu bảo mật
	AlertTypeSuspicious  = "suspicious"     // Hoạt động bất thường
)

// Alert severity levels
const (
	SeverityCritical = "critical" // Cần xử lý ngay
	SeverityHigh     = "high"     // Ưu tiên cao
	SeverityMedium   = "medium"   // Ưu tiên trung bình
	SeverityLow      = "low"      // Ưu tiên thấp
	SeverityInfo     = "info"     // Thông tin
)

// Alert status
const (
	AlertStatusOpen      = "open"       // Chưa xử lý
	AlertStatusAcknonwed = "acknowledged" // Đã biết về
	AlertStatusResolved  = "resolved"   // Đã giải quyết
)

// Alert struct - Lưu thông tin một cảnh báo
type Alert struct {
	ID          string    `json:"id"`
	AssetID     string    `json:"asset_id"`
	ScanJobID   *string   `json:"scan_job_id,omitempty"` // Reference đến scan job nếu có
	AlertType   string    `json:"alert_type"`            // high_risk, expired_cert, etc
	Severity    string    `json:"severity"`              // critical, high, medium, low, info
	Title       string    `json:"title"`                 // Tiêu đề cảnh báo
	Description string    `json:"description"`           // Chi tiết mô tả
	Details     string    `json:"details"`               // JSON data về chi tiết alert
	Status      string    `json:"status"`                // open, acknowledged, resolved
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// AlertFilter - Để filter alerts khi query
type AlertFilter struct {
	AssetID   string   // Lọc theo asset
	AlertType string   // Lọc theo loại alert
	Severity  string   // Lọc theo mức độ  
	Status    string   // Lọc theo trạng thái
	Page      int
	Limit     int
}

// AlertStats - Thống kê cảnh báo
type AlertStats struct {
	TotalAlerts     int `json:"total_alerts"`
	CriticalCount   int `json:"critical_count"`
	HighCount       int `json:"high_count"`
	MediumCount     int `json:"medium_count"`
	LowCount        int `json:"low_count"`
	OpenCount       int `json:"open_count"`
	AcknowledgedCount int `json:"acknowledged_count"`
	ResolvedCount   int `json:"resolved_count"`
}
