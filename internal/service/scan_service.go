package service

import (
	"fmt"
	"time"

	"mini-asm/internal/model"
	"mini-asm/internal/scanner"
)

type ScanService struct {
	// Sau này bạn sẽ truyền Storage (Database) vào đây để lưu Job
	// db storage.ScanStorage
}

func NewScanService() *ScanService {
	return &ScanService{}
}

// Hàm này trả về 1 Job ID và chạy ngầm (Goroutine) quá trình scan
func (s *ScanService) StartScan(assetID string, assetTarget string, scanType string) (*model.ScanJob, error) {
	jobID := fmt.Sprintf("job-%d", time.Now().Unix()) // Tạo ID tạm

	job := &model.ScanJob{
		ID:        jobID,
		AssetID:   assetID,
		ScanType:  scanType,
		Status:    "pending", // Mới tạo thì pending
		StartedAt: time.Now(),
	}

	// Tùy theo loại scan mà chọn công cụ tương ứng
	var tool model.Scanner
	switch scanType {
	case model.ScanTypeIP:
		tool = scanner.NewIPScanner()
	case model.ScanTypePort:
		tool = scanner.NewPortScanner()
	case "tech": // Thêm cái này để hết lỗi 400
		// Nếu chưa viết tech_scanner thì tạm thời mượn IPScanner để test
		tool = scanner.NewIPScanner()
	default:
		return nil, fmt.Errorf("loại scan không được hỗ trợ: %s", scanType)
	}

	// CHẠY NGẦM (Async) - Để API có thể trả về ngay lập tức (202 Accepted)
	go func() {
		// Cập nhật trạng thái thành running (Thực tế bạn sẽ lưu xuống DB ở bước này)
		job.Status = "running"

		results, err := tool.Scan(assetTarget) // Bắt đầu quét thật

		now := time.Now()
		job.EndedAt = &now

		if err != nil {
			job.Status = "failed"
			job.Error = err.Error()
		} else {
			job.Status = "completed"
			job.Results = results
		}

		// (Thực tế bạn sẽ Update Job trong DB ở bước này)
		fmt.Printf("Scan %s hoàn tất! Trạng thái: %s\n", job.ID, job.Status)
	}()

	return job, nil
}
