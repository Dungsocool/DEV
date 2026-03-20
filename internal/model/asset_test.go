package model

import (
	"testing"
)

func TestAssetValidation(t *testing.T) {
	// Khai báo các trường hợp cần test (Table-driven tests)
	tests := []struct {
		name    string
		asset   Asset
		wantErr bool
	}{
		{
			name:    "Domain hợp lệ",
			asset:   Asset{Name: "google.com", Type: "domain"},
			wantErr: false,
		},
		{
			name:    "IP hợp lệ",
			asset:   Asset{Name: "127.0.0.1", Type: "ip"},
			wantErr: false,
		},
		{
			name:    "Loại asset không hợp lệ",
			asset:   Asset{Name: "test", Type: "vps"}, // Chỉ cho phép domain/ip
			wantErr: true,
		},
		{
			name:    "Tên trống",
			asset:   Asset{Name: "", Type: "domain"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Giả sử bạn có hàm Validate() trong struct Asset
			// Nếu chưa có, hãy thêm logic kiểm tra ở đây
			err := tt.asset.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Asset.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
