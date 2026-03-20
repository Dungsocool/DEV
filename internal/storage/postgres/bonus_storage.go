package postgres

import (
	"fmt"
	"strings"
	"mini-asm/internal/model"
)

// bài 6
func (s *PostgresStorage) GetAssets(page, limit int, assetType, status string) ([]*model.Asset, int, error) {
	conditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if assetType != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, assetType)
		argIdx++
	}
	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Đếm tổng số
	var total int
	s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM assets %s", whereClause), args...).Scan(&total)

	// Lấy dữ liệu
	offset := (page - 1) * limit
	dataQuery := fmt.Sprintf("SELECT id, name, type, status, created_at, updated_at FROM assets %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d", whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(dataQuery, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()

	assets := make([]*model.Asset, 0)
	for rows.Next() {
		var a model.Asset
		rows.Scan(&a.ID, &a.Name, &a.Type, &a.Status, &a.CreatedAt, &a.UpdatedAt)
		assets = append(assets, &a)
	}
	return assets, total, nil
}

//  Bài 7 
func (s *PostgresStorage) SearchAssets(query string) ([]*model.Asset, error) {
	sqlQuery := "SELECT id, name, type, status, created_at, updated_at FROM assets WHERE name ILIKE $1 LIMIT 100"
	rows, err := s.db.Query(sqlQuery, "%"+query+"%")
	if err != nil { return nil, err }
	defer rows.Close()

	assets := make([]*model.Asset, 0)
	for rows.Next() {
		var a model.Asset
		rows.Scan(&a.ID, &a.Name, &a.Type, &a.Status, &a.CreatedAt, &a.UpdatedAt)
		assets = append(assets, &a)
	}
	return assets, nil
}
