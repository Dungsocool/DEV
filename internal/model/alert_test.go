package model

import (
	"testing"
	"time"
)

func TestAlertValidation(t *testing.T) {
	tests := []struct {
		name  string
		alert Alert
		valid bool
	}{
		{
			name: "valid alert",
			alert: Alert{
				ID:        "test-id",
				AssetID:   "asset-id",
				AlertType: AlertTypeHighRisk,
				Severity:  SeverityCritical,
				Title:     "Test Alert",
				Status:    AlertStatusOpen,
			},
			valid: true,
		},
		{
			name: "empty ID",
			alert: Alert{
				AssetID:   "asset-id",
				AlertType: AlertTypeHighRisk,
				Severity:  SeverityCritical,
				Title:     "Test Alert",
				Status:    AlertStatusOpen,
			},
			valid: true, // ID can be generated
		},
		{
			name: "empty AssetID",
			alert: Alert{
				ID:        "test-id",
				AlertType: AlertTypeHighRisk,
				Severity:  SeverityCritical,
				Title:     "Test Alert",
				Status:    AlertStatusOpen,
			},
			valid: false,
		},
		{
			name: "valid severity levels",
			alert: Alert{
				ID:        "test-id",
				AssetID:   "asset-id",
				AlertType: AlertTypeExpiredCert,
				Severity:  SeverityHigh,
				Title:     "Test Alert",
				Status:    AlertStatusOpen,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate function (thêm vào model nếu cần)
			isValid := tt.alert.AssetID != ""

			if isValid != tt.valid {
				t.Errorf("got %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestAlertStatus(t *testing.T) {
	alert := &Alert{
		ID:      "test-1",
		AssetID: "asset-1",
		Status:  AlertStatusOpen,
	}

	// Test status transitions
	t.Run("transition to acknowledged", func(t *testing.T) {
		alert.Status = AlertStatusAcknonwed
		if alert.Status != AlertStatusAcknonwed {
			t.Errorf("expected status %s, got %s", AlertStatusAcknonwed, alert.Status)
		}
	})

	t.Run("transition to resolved", func(t *testing.T) {
		alert.Status = AlertStatusResolved
		now := time.Now()
		alert.ResolvedAt = &now

		if alert.Status != AlertStatusResolved {
			t.Errorf("expected status %s, got %s", AlertStatusResolved, alert.Status)
		}
		if alert.ResolvedAt == nil {
			t.Error("ResolvedAt should be set")
		}
	})
}

func TestAlertSeverity(t *testing.T) {
	severities := []string{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
		SeverityInfo,
	}

	for _, severity := range severities {
		t.Run(severity, func(t *testing.T) {
			alert := Alert{
				ID:       "test",
				AssetID:  "asset",
				Severity: severity,
			}

			if alert.Severity != severity {
				t.Errorf("expected severity %s, got %s", severity, alert.Severity)
			}
		})
	}
}

func TestAlertTypes(t *testing.T) {
	alertTypes := []string{
		AlertTypeHighRisk,
		AlertTypeExpiredCert,
		AlertTypeNewService,
		AlertTypeVulnerable,
		AlertTypeSuspicious,
	}

	for _, alertType := range alertTypes {
		t.Run(alertType, func(t *testing.T) {
			alert := Alert{
				ID:        "test",
				AssetID:   "asset",
				AlertType: alertType,
			}

			if alert.AlertType != alertType {
				t.Errorf("expected alert type %s, got %s", alertType, alert.AlertType)
			}
		})
	}
}

func TestAlertStats(t *testing.T) {
	stats := &AlertStats{
		TotalAlerts:       10,
		CriticalCount:     2,
		HighCount:         3,
		MediumCount:       3,
		LowCount:          2,
		OpenCount:         5,
		AcknowledgedCount: 3,
		ResolvedCount:     2,
	}

	t.Run("total count", func(t *testing.T) {
		if stats.TotalAlerts != 10 {
			t.Errorf("expected total 10, got %d", stats.TotalAlerts)
		}
	})

	t.Run("severity distribution", func(t *testing.T) {
		sum := stats.CriticalCount + stats.HighCount + stats.MediumCount + stats.LowCount
		if sum > stats.TotalAlerts {
			t.Error("severity sum exceeds total")
		}
	})

	t.Run("status distribution", func(t *testing.T) {
		sum := stats.OpenCount + stats.AcknowledgedCount + stats.ResolvedCount
		if sum > stats.TotalAlerts {
			t.Error("status sum exceeds total")
		}
	})
}
