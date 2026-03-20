package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"mini-asm/internal/storage/postgres"
)

type HealthHandler struct {
	storage *postgres.PostgresStorage
}

func NewHealthHandler(s *postgres.PostgresStorage) *HealthHandler {
	return &HealthHandler{storage: s}
}

// Trả về trạng thái của Database
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Gọi hàm Ping từ lớp Storage để kiểm tra kết nối
	err := h.storage.Ping()
	
	// Nếu rớt kết nối DB -> Trả về lỗi 503
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "degraded",
			"database":  map[string]string{"status": "disconnected"},
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	// Nếu kết nối tốt -> Trả về 200 OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"database": map[string]string{"status": "connected"},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}