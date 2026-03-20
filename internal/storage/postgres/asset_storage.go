package postgres

import (
	"fmt"
	"mini-asm/internal/model"
)

// BatchCreate thực hiện lưu nhiều tài sản cùng lúc dùng Transaction
func (s *PostgresStorage) BatchCreate(assets []*model.Asset) error {
	// Bắt đầu Transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Nếu hàm kết thúc mà chưa Commit, nó sẽ tự động Rollback
	defer tx.Rollback()

	// Câu lệnh SQL để chèn dữ liệu
	query := `INSERT INTO assets (id, name, type, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	// Lặp qua từng tài sản để đưa vào DB
	for _, asset := range assets {
		_, err := tx.Exec(query, asset.ID, asset.Name, asset.Type, asset.Status, asset.CreatedAt, asset.UpdatedAt)
		if err != nil {
			// Nếu có MỘT lỗi bất kỳ, return lỗi ngay, defer phía trên sẽ tự động Rollback TOÀN BỘ
			return fmt.Errorf("lỗi khi insert asset %s: %v", asset.Name, err)
		}
	}

	// Nếu chạy hết vòng lặp mà không có lỗi -> Xác nhận lưu toàn bộ vào DB
	return tx.Commit()
}