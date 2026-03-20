package postgres

import (
	"fmt"
	"strings"
	"mini-asm/internal/model"
)

// GetStats đếm tổng số, đếm theo loại và trạng thái
func (s *PostgresStorage) GetStats() (*model.AssetStats, error) {
	stats := &model.AssetStats{
		ByType:   make(map[string]int),
		ByStatus: make(map[string]int),
	}

	s.db.QueryRow("SELECT COUNT(*) FROM assets").Scan(&stats.Total)

	rowsType, _ := s.db.Query("SELECT type, COUNT(*) FROM assets GROUP BY type")
	defer rowsType.Close()
	for rowsType.Next() {
		var t string
		var c int
		rowsType.Scan(&t, &c)
		stats.ByType[t] = c
	}

	rowsStatus, _ := s.db.Query("SELECT status, COUNT(*) FROM assets GROUP BY status")
	defer rowsStatus.Close()
	for rowsStatus.Next() {
		var st string
		var c int
		rowsStatus.Scan(&st, &c)
		stats.ByStatus[st] = c
	}

	return stats, nil
}

// CountAssets đếm theo bộ lọc
func (s *PostgresStorage) CountAssets(assetType string, status string) (int, error) {
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if assetType != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, assetType)
		argIndex++
	}
	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf("SELECT COUNT(*) FROM assets %s", whereClause)
	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	return count, err
}