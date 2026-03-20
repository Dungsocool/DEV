package model

import "time"

// Định nghĩa các loại scan hỗ trợ
const (
	ScanTypeDNS  = "dns"
	ScanTypeIP   = "ip"
	ScanTypePort = "port"
)

// Scanner là interface chung mà tất cả các tool scan (IP, Port, DNS...) phải tuân theo
type Scanner interface {
	Scan(target string) (interface{}, error)
}

// Struct lưu thông tin 1 Job Scan
type ScanJob struct {
	ID        string      `json:"id"`
	AssetID   string      `json:"asset_id"`
	ScanType  string      `json:"scan_type"`
	Status    string      `json:"status"` // pending, running, completed, failed
	StartedAt time.Time   `json:"started_at"`
	EndedAt   *time.Time  `json:"ended_at"`
	Error     string      `json:"error"`
	Results   interface{} `json:"results"` // Chứa kết quả scan
	CreatedAt time.Time   `json:"created_at"`
}

// -----------------------------------------
// Struct cho kết quả trả về của IP Scan
// -----------------------------------------
type IPScanResult struct {
	IPAddress   string      `json:"ip_address"`
	Geolocation Geolocation `json:"geolocation"`
	ASN         ASNInfo     `json:"asn"`
	ReverseDNS  string      `json:"reverse_dns"`
	CreatedAt   time.Time   `json:"created_at"`
}

type Geolocation struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
}

type ASNInfo struct {
	Number      int    `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// -----------------------------------------
// Struct cho kết quả trả về của Port Scan
// -----------------------------------------
type PortScanResult struct {
	IPAddress      string     `json:"ip_address"`
	OpenPorts      []PortInfo `json:"open_ports"`
	ClosedPorts    int        `json:"closed_ports"`
	TotalScanned   int        `json:"total_scanned"`
	ScanDurationMs int64      `json:"scan_duration_ms"`
	CreatedAt      time.Time  `json:"created_at"`
}

type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Service  string `json:"service"`
	Version  string `json:"version"`
}
