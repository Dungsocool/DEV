package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"mini-asm/internal/model"
	"mini-asm/internal/service"
)

// AlertHandler - HTTP handler cho alerts
type AlertHandler struct {
	alertService service.AlertService
}

// NewAlertHandler - Constructor
func NewAlertHandler(alertService service.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

// RegisterRoutes - Register alert routes
func (h *AlertHandler) RegisterRoutes(router *mux.Router) {
	// GET /assets/{id}/alerts - Get alerts for an asset
	router.HandleFunc("/assets/{id}/alerts", h.GetAssetAlerts).Methods(http.MethodGet)

	// GET /alerts/{id} - Get alert by ID
	router.HandleFunc("/alerts/{id}", h.GetAlert).Methods(http.MethodGet)

	// PUT /alerts/{id}/status - Update alert status
	router.HandleFunc("/alerts/{id}/status", h.UpdateAlertStatus).Methods(http.MethodPut)

	// DELETE /alerts/{id} - Delete alert
	router.HandleFunc("/alerts/{id}", h.DeleteAlert).Methods(http.MethodDelete)

	// GET /assets/{id}/alerts/stats - Get alert statistics
	router.HandleFunc("/assets/{id}/alerts/stats", h.GetAlertStats).Methods(http.MethodGet)

	// GET /alerts - List all alerts with filters
	router.HandleFunc("/alerts", h.ListAlerts).Methods(http.MethodGet)
}

// GetAlert - GET /alerts/{id}
func (h *AlertHandler) GetAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]

	alert, err := h.alertService.GetAlert(ctx, id)
	if err != nil {
		http.Error(w, "Alert not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

// GetAssetAlerts - GET /assets/{id}/alerts
func (h *AlertHandler) GetAssetAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	assetID := params["id"]

	// Parse pagination
	page := 1
	limit := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	alerts, total, err := h.alertService.GetAssetAlerts(ctx, assetID, page, limit)
	if err != nil {
		http.Error(w, "Failed to get alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       alerts,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_page": (total + limit - 1) / limit,
	})
}

// UpdateAlertStatus - PUT /alerts/{id}/status
func (h *AlertHandler) UpdateAlertStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.alertService.UpdateAlertStatus(ctx, id, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated alert
	alert, _ := h.alertService.GetAlert(ctx, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

// DeleteAlert - DELETE /alerts/{id}
func (h *AlertHandler) DeleteAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]

	if err := h.alertService.DeleteAlert(ctx, id); err != nil {
		http.Error(w, "Failed to delete alert", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAlertStats - GET /assets/{id}/alerts/stats
func (h *AlertHandler) GetAlertStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	assetID := params["id"]

	stats, err := h.alertService.GetAlertStats(ctx, assetID)
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ListAlerts - GET /alerts
func (h *AlertHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse filters
	filter := &model.AlertFilter{
		AssetID:   r.URL.Query().Get("asset_id"),
		AlertType: r.URL.Query().Get("alert_type"),
		Severity:  r.URL.Query().Get("severity"),
		Status:    r.URL.Query().Get("status"),
	}

	// Parse pagination
	page := 1
	limit := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	filter.Page = page
	filter.Limit = limit

	alerts, total, err := h.alertService.ListAlerts(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to list alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       alerts,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_page": (total + limit - 1) / limit,
	})
}
