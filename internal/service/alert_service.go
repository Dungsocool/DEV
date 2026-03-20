package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"mini-asm/internal/model"
	"mini-asm/internal/storage"
	"mini-asm/internal/storage/postgres"
)

// AlertService - Interface cho xử lý logic alerts
type AlertService interface {
	CreateAlert(ctx context.Context, alert *model.Alert) error
	GetAlert(ctx context.Context, id string) (*model.Alert, error)
	GetAssetAlerts(ctx context.Context, assetID string, page, limit int) ([]model.Alert, int, error)
	UpdateAlertStatus(ctx context.Context, id string, status string) error
	GetAlertStats(ctx context.Context, assetID string) (*model.AlertStats, error)
	DeleteAlert(ctx context.Context, id string) error
	ListAlerts(ctx context.Context, filter *model.AlertFilter) ([]model.Alert, int, error)
	
	// Generate alerts từ kết quả scan
	GenerateAlertsFromScan(ctx context.Context, scanJob *model.ScanJob, assetID string, results interface{}) error
}

// DefaultAlertService - Implementation
type DefaultAlertService struct {
	alertStorage    postgres.AlertStorage
	assetStorage    storage.Storage
}

// NewAlertService - Constructor
func NewAlertService(alertStorage postgres.AlertStorage, assetStorage storage.Storage) AlertService {
	return &DefaultAlertService{
		alertStorage:    alertStorage,
		assetStorage:    assetStorage,
	}
}

// CreateAlert - Tạo alert mới
func (s *DefaultAlertService) CreateAlert(ctx context.Context, alert *model.Alert) error {
	if alert.ID == "" {
		alert.ID = uuid.New().String()
	}

	if alert.Status == "" {
		alert.Status = model.AlertStatusOpen
	}

	return s.alertStorage.Create(ctx, alert)
}

// GetAlert - Lấy alert theo ID
func (s *DefaultAlertService) GetAlert(ctx context.Context, id string) (*model.Alert, error) {
	return s.alertStorage.GetByID(ctx, id)
}

// GetAssetAlerts - Lấy alerts của asset
func (s *DefaultAlertService) GetAssetAlerts(ctx context.Context, assetID string, page, limit int) ([]model.Alert, int, error) {
	return s.alertStorage.GetByAssetID(ctx, assetID, page, limit)
}

// UpdateAlertStatus - Cập nhật status alert
func (s *DefaultAlertService) UpdateAlertStatus(ctx context.Context, id string, status string) error {
	validStatuses := map[string]bool{
		model.AlertStatusOpen:      true,
		model.AlertStatusAcknonwed: true,
		model.AlertStatusResolved:  true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	return s.alertStorage.UpdateStatus(ctx, id, status)
}

// GetAlertStats - Lấy thống kê alerts
func (s *DefaultAlertService) GetAlertStats(ctx context.Context, assetID string) (*model.AlertStats, error) {
	return s.alertStorage.GetStats(ctx, assetID)
}

// DeleteAlert - Xóa alert
func (s *DefaultAlertService) DeleteAlert(ctx context.Context, id string) error {
	return s.alertStorage.Delete(ctx, id)
}

// ListAlerts - List alerts với filter
func (s *DefaultAlertService) ListAlerts(ctx context.Context, filter *model.AlertFilter) ([]model.Alert, int, error) {
	return s.alertStorage.List(ctx, filter)
}

// GenerateAlertsFromScan - Generate alerts từ scan results
func (s *DefaultAlertService) GenerateAlertsFromScan(ctx context.Context, scanJob *model.ScanJob, assetID string, results interface{}) error {
	switch scanJob.ScanType {
	case "port":
		return s.generatePortScanAlerts(ctx, scanJob, assetID, results)
	case "ssl":
		return s.generateSSLScanAlerts(ctx, scanJob, assetID, results)
	case "ip":
		return s.generateIPScanAlerts(ctx, scanJob, assetID, results)
	case "tech":
		return s.generateTechScanAlerts(ctx, scanJob, assetID, results)
	}

	return nil
}

// generatePortScanAlerts - Phát hiện alerts từ port scan
func (s *DefaultAlertService) generatePortScanAlerts(ctx context.Context, scanJob *model.ScanJob, assetID string, results interface{}) error {
	// Parse results
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return err
	}

	var portResults []map[string]interface{}
	err = json.Unmarshal(resultsJSON, &portResults)
	if err != nil {
		return err
	}

	if len(portResults) == 0 {
		return nil
	}

	// Get first result (port scan usually has one result per IP)
	result := portResults[0]
	openPorts, ok := result["open_ports"].([]interface{})
	if !ok {
		return nil
	}

	// High risk ports
	highRiskPorts := map[int]bool{
		22:   true, // SSH
		3389: true, // RDP
		5432: true, // PostgreSQL
		3306: true, // MySQL
		6379: true, // Redis
		1433: true, // MSSQL
	}

	// Check for high-risk ports
	for _, port := range openPorts {
		portData, ok := port.(map[string]interface{})
		if !ok {
			continue
		}

		portNum, ok := portData["port"].(float64)
		if !ok {
			continue
		}

		if highRiskPorts[int(portNum)] {
			details, _ := json.Marshal(portData)
			alert := &model.Alert{
				ID:          uuid.New().String(),
				AssetID:     assetID,
				ScanJobID:   &scanJob.ID,
				AlertType:   model.AlertTypeHighRisk,
				Severity:    model.SeverityCritical,
				Title:       fmt.Sprintf("High-Risk Port Detected: %d", int(portNum)),
				Description: fmt.Sprintf("Port %d is exposed on this asset. This port is commonly targeted by attackers.", int(portNum)),
				Details:     string(details),
				Status:      model.AlertStatusOpen,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			if err := s.alertStorage.Create(ctx, alert); err != nil {
				return fmt.Errorf("failed to create port alert: %w", err)
			}
		}
	}

	return nil
}

// generateSSLScanAlerts - Phát hiện alerts từ SSL scan
func (s *DefaultAlertService) generateSSLScanAlerts(ctx context.Context, scanJob *model.ScanJob, assetID string, results interface{}) error {
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return err
	}

	var sslResults []map[string]interface{}
	err = json.Unmarshal(resultsJSON, &sslResults)
	if err != nil {
		return err
	}

	if len(sslResults) == 0 {
		return nil
	}

	result := sslResults[0]
	cert, ok := result["certificate"].(map[string]interface{})
	if !ok {
		return nil
	}

	// Check certificate expiry
	daysUntilExpiry, ok := cert["days_until_expiry"].(float64)
	if ok && daysUntilExpiry < 30 {
		severity := model.SeverityMedium
		if daysUntilExpiry < 7 {
			severity = model.SeverityCritical
		}

		details, _ := json.Marshal(cert)
		alert := &model.Alert{
			ID:          uuid.New().String(),
			AssetID:     assetID,
			ScanJobID:   &scanJob.ID,
			AlertType:   model.AlertTypeExpiredCert,
			Severity:    severity,
			Title:       "Certificate Expiring Soon",
			Description: fmt.Sprintf("SSL certificate expires in %.0f days", daysUntilExpiry),
			Details:     string(details),
			Status:      model.AlertStatusOpen,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.alertStorage.Create(ctx, alert); err != nil {
			return fmt.Errorf("failed to create SSL alert: %w", err)
		}
	}

	return nil
}

// generateIPScanAlerts - Generate alerts từ IP scan
func (s *DefaultAlertService) generateIPScanAlerts(ctx context.Context, scanJob *model.ScanJob, assetID string, results interface{}) error {
	// IP scan thường không có alerts
	return nil
}

// generateTechScanAlerts - Generate alerts từ tech scan
func (s *DefaultAlertService) generateTechScanAlerts(ctx context.Context, scanJob *model.ScanJob, assetID string, results interface{}) error {
	// Có thể phát hiện outdated technologies, vulnerable versions, etc
	return nil
}
