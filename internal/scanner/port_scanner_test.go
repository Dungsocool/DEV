package scanner

import (
	"testing"
)

func TestPortScanner_SafetyCheck(t *testing.T) {
	// Khởi tạo scanner (biến s)
	s := NewPortScanner()

	tests := []struct {
		name   string
		target string
		isSafe bool
	}{
		{"Localhost phai an toan", "127.0.0.1", true},
		{"IP Noi bo phai an toan", "192.168.1.1", true},
		{"IP Google phai bi chan", "8.8.8.8", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// CÁCH SỬA: Dùng s để gọi hàm Scan và kiểm tra lỗi
			_, err := s.Scan(tt.target)

			// Nếu là IP không an toàn (isSafe = false) thì PHẢI có lỗi
			if !tt.isSafe && err == nil {
				t.Errorf("%s: Mong doi loi bao mat cho %s nhung không thay", tt.name, tt.target)
			}

			// Nếu là IP an toàn (isSafe = true) thì KHÔNG ĐƯỢC có lỗi bảo mật
			if tt.isSafe && err != nil && err.Error() == "cảnh báo bảo mật: Chỉ được phép port scan trên mạng nội bộ hoặc localhost (127.0.0.1)" {
				t.Errorf("%s: %s la an toan nhung lai bi chan", tt.name, tt.target)
			}
		})
	}
}
