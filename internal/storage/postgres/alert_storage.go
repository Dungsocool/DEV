package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"mini-asm/internal/model"
)

// AlertStorage - Interface để tạo/lấy/cập nhật alerts
type AlertStorage interface {
	Create(ctx context.Context, alert *model.Alert) error
	GetByID(ctx context.Context, id string) (*model.Alert, error)
	GetByAssetID(ctx context.Context, assetID string, page, limit int) ([]model.Alert, int, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	GetStats(ctx context.Context, assetID string) (*model.AlertStats, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *model.AlertFilter) ([]model.Alert, int, error)
}

// PostgresAlertStorage - Implementation của AlertStorage
type PostgresAlertStorage struct {
	db *sql.DB
}

// NewPostgresAlertStorage - Constructor
func NewPostgresAlertStorage(db *sql.DB) AlertStorage {
	return &PostgresAlertStorage{db: db}
}

// Create - Tạo alert mới
func (s *PostgresAlertStorage) Create(ctx context.Context, alert *model.Alert) error {
	query := `
		INSERT INTO alerts (id, asset_id, scan_job_id, alert_type, severity, title, description, details, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	// Convert details to JSON
	var detailsJSON *string
	if alert.Details != "" {
		detailsJSON = &alert.Details
	}

	err := s.db.QueryRowContext(ctx, query,
		alert.ID,
		alert.AssetID,
		alert.ScanJobID,
		alert.AlertType,
		alert.Severity,
		alert.Title,
		alert.Description,
		detailsJSON,
		alert.Status,
		time.Now(),
		time.Now(),
	).Scan(&alert.ID, &alert.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}

	return nil
}

// GetByID - Lấy alert theo ID
func (s *PostgresAlertStorage) GetByID(ctx context.Context, id string) (*model.Alert, error) {
	query := `
		SELECT id, asset_id, scan_job_id, alert_type, severity, title, description, details, status, created_at, updated_at, resolved_at
		FROM alerts
		WHERE id = $1
	`

	alert := &model.Alert{}
	var details sql.NullString
	var scanJobID sql.NullString
	var resolvedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&alert.ID,
		&alert.AssetID,
		&scanJobID,
		&alert.AlertType,
		&alert.Severity,
		&alert.Title,
		&alert.Description,
		&details,
		&alert.Status,
		&alert.CreatedAt,
		&alert.UpdatedAt,
		&resolvedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("alert not found")
		}
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	if scanJobID.Valid {
		alert.ScanJobID = &scanJobID.String
	}
	if details.Valid {
		alert.Details = details.String
	}
	if resolvedAt.Valid {
		alert.ResolvedAt = &resolvedAt.Time
	}

	return alert, nil
}

// GetByAssetID - Lấy tất cả alerts của một asset
func (s *PostgresAlertStorage) GetByAssetID(ctx context.Context, assetID string, page, limit int) ([]model.Alert, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Count total
	countQuery := "SELECT COUNT(*) FROM alerts WHERE asset_id = $1"
	var total int
	err := s.db.QueryRowContext(ctx, countQuery, assetID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count alerts: %w", err)
	}

	// Get alerts
	query := `
		SELECT id, asset_id, scan_job_id, alert_type, severity, title, description, details, status, created_at, updated_at, resolved_at
		FROM alerts
		WHERE asset_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, assetID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get alerts: %w", err)
	}
	defer rows.Close()

	alerts := []model.Alert{}
	for rows.Next() {
		alert := model.Alert{}
		var details sql.NullString
		var scanJobID sql.NullString
		var resolvedAt sql.NullTime

		err := rows.Scan(
			&alert.ID,
			&alert.AssetID,
			&scanJobID,
			&alert.AlertType,
			&alert.Severity,
			&alert.Title,
			&alert.Description,
			&details,
			&alert.Status,
			&alert.CreatedAt,
			&alert.UpdatedAt,
			&resolvedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan alert row: %w", err)
		}

		if scanJobID.Valid {
			alert.ScanJobID = &scanJobID.String
		}
		if details.Valid {
			alert.Details = details.String
		}
		if resolvedAt.Valid {
			alert.ResolvedAt = &resolvedAt.Time
		}

		alerts = append(alerts, alert)
	}

	return alerts, total, nil
}

// UpdateStatus - Cập nhật trạng thái alert
func (s *PostgresAlertStorage) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE alerts
		SET status = $1, updated_at = $2
	`

	var args []interface{}
	args = append(args, status)
	args = append(args, time.Now())

	resolvedAt := ""
	if status == model.AlertStatusResolved {
		query += `, resolved_at = $3`
		resolvedAt = time.Now().Format(time.RFC3339)
		args = append(args, resolvedAt)
	}

	query += ` WHERE id = $%d`
	query = fmt.Sprintf(query, len(args)+1)

	result, err := s.db.ExecContext(ctx, query, append(args, id)...)
	if err != nil {
		return fmt.Errorf("failed to update alert status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("alert not found")
	}

	return nil
}

// GetStats - Lấy thống kê alerts của asset
func (s *PostgresAlertStorage) GetStats(ctx context.Context, assetID string) (*model.AlertStats, error) {
	query := `
		SELECT
			COUNT(*) as total_alerts,
			SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END) as critical_count,
			SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END) as high_count,
			SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) as medium_count,
			SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END) as low_count,
			SUM(CASE WHEN status = 'open' THEN 1 ELSE 0 END) as open_count,
			SUM(CASE WHEN status = 'acknowledged' THEN 1 ELSE 0 END) as acknowledged_count,
			SUM(CASE WHEN status = 'resolved' THEN 1 ELSE 0 END) as resolved_count
		FROM alerts
		WHERE asset_id = $1
	`

	stats := &model.AlertStats{}
	err := s.db.QueryRowContext(ctx, query, assetID).Scan(
		&stats.TotalAlerts,
		&stats.CriticalCount,
		&stats.HighCount,
		&stats.MediumCount,
		&stats.LowCount,
		&stats.OpenCount,
		&stats.AcknowledgedCount,
		&stats.ResolvedCount,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get alert stats: %w", err)
	}

	return stats, nil
}

// Delete - Xóa alert
func (s *PostgresAlertStorage) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM alerts WHERE id = $1"

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("alert not found")
	}

	return nil
}

// List - List alerts với filter
func (s *PostgresAlertStorage) List(ctx context.Context, filter *model.AlertFilter) ([]model.Alert, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	offset := (filter.Page - 1) * filter.Limit

	// Build query
	queryWhere := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.AssetID != "" {
		queryWhere = append(queryWhere, fmt.Sprintf("asset_id = $%d", argIndex))
		args = append(args, filter.AssetID)
		argIndex++
	}

	if filter.AlertType != "" {
		queryWhere = append(queryWhere, fmt.Sprintf("alert_type = $%d", argIndex))
		args = append(args, filter.AlertType)
		argIndex++
	}

	if filter.Severity != "" {
		queryWhere = append(queryWhere, fmt.Sprintf("severity = $%d", argIndex))
		args = append(args, filter.Severity)
		argIndex++
	}

	if filter.Status != "" {
		queryWhere = append(queryWhere, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	whereClause := ""
	if len(queryWhere) > 0 {
		whereClause = "WHERE " + strings.Join(queryWhere, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM alerts %s", whereClause)
	var total int
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count alerts: %w", err)
	}

	// Get alerts
	query := fmt.Sprintf(`
		SELECT id, asset_id, scan_job_id, alert_type, severity, title, description, details, status, created_at, updated_at, resolved_at
		FROM alerts
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex+1, argIndex+2)

	args = append(args, filter.Limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get alerts: %w", err)
	}
	defer rows.Close()

	alerts := []model.Alert{}
	for rows.Next() {
		alert := model.Alert{}
		var details sql.NullString
		var scanJobID sql.NullString
		var resolvedAt sql.NullTime

		err := rows.Scan(
			&alert.ID,
			&alert.AssetID,
			&scanJobID,
			&alert.AlertType,
			&alert.Severity,
			&alert.Title,
			&alert.Description,
			&details,
			&alert.Status,
			&alert.CreatedAt,
			&alert.UpdatedAt,
			&resolvedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan alert row: %w", err)
		}

		if scanJobID.Valid {
			alert.ScanJobID = &scanJobID.String
		}
		if details.Valid {
			alert.Details = details.String
		}
		if resolvedAt.Valid {
			alert.ResolvedAt = &resolvedAt.Time
		}

		alerts = append(alerts, alert)
	}

	return alerts, total, nil
}
