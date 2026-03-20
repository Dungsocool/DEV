package postgres

import (
	"fmt"
	"strings"
)

// BatchDelete xóa nhiều ID cùng lúc
func (s *PostgresStorage) BatchDelete(ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// Tạo danh sách biến $1, $2, $3... cho câu lệnh SQL
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM assets WHERE id IN (%s)", strings.Join(placeholders, ","))
	
	result, err := s.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	deletedCount, err := result.RowsAffected()
	return int(deletedCount), err
}