package scanner

import (
	"errors"
	"fmt"
	"net"
	"time"

	"mini-asm/internal/model" // Thay "your_project" bằng tên module của bạn
)

type PortScanner struct {
	timeout time.Duration
}

func NewPortScanner() *PortScanner {
	return &PortScanner{
		timeout: 1 * time.Second, // Chờ tối đa 1 giây cho mỗi port
	}
}

// Kiểm tra bảo mật: Chỉ cho quét localhost hoặc mạng nội bộ (10.x, 192.x)
func isSafeToScan(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate()
}

func (s *PortScanner) Scan(targetIP string) (interface{}, error) {
	// BẮT BUỘC: Chặn quét IP ngoài Internet
	if !isSafeToScan(targetIP) {
		return nil, errors.New("cảnh báo bảo mật: Chỉ được phép port scan trên mạng nội bộ hoặc localhost (127.0.0.1)")
	}

	startTime := time.Now()
	var openPorts []model.PortInfo
	closedCount := 0

	// Quét thử các port phổ biến
	portsToScan := []int{21, 22, 23, 80, 443, 3306, 5432, 6379, 8080}

	for _, port := range portsToScan {
		target := fmt.Sprintf("%s:%d", targetIP, port)
		conn, err := net.DialTimeout("tcp", target, s.timeout)

		if err != nil {
			closedCount++ // Bị lỗi kết nối -> Port đang đóng
			continue
		}
		conn.Close()

		// Kết nối thành công -> Port đang mở
		openPorts = append(openPorts, model.PortInfo{
			Port:     port,
			Protocol: "tcp",
			State:    "open",
			Service:  "unknown",
		})
	}

	duration := time.Since(startTime).Milliseconds()

	result := model.PortScanResult{
		IPAddress:      targetIP,
		OpenPorts:      openPorts,
		ClosedPorts:    closedCount,
		TotalScanned:   len(portsToScan),
		ScanDurationMs: duration,
		CreatedAt:      time.Now(),
	}

	return []model.PortScanResult{result}, nil
}
