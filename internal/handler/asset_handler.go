package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"mini-asm/internal/model"
	"mini-asm/internal/service"

	"github.com/gorilla/mux"
)

type AssetHandler struct {
	service     *service.AssetService
	scanService *service.ScanService // Dùng cho Bài 1: Quét IP/Port
}

// Constructor nhận 2 service để khớp với main.go
func NewAssetHandler(s *service.AssetService, ss *service.ScanService) *AssetHandler {
	return &AssetHandler{
		service:     s,
		scanService: ss,
	}
}

// --- CÁC HÀM MỚI CHO BÀI 1 & BÀI 3 ---

// CreateAsset: Tạo 1 asset đơn lẻ (Sửa lỗi 405 Method Not Allowed)
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var a model.Asset
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{Code: 400, Status: "fail", Message: "Dữ liệu không hợp lệ"})
		return
	}

	// Gọi service BatchCreate để tận dụng code cũ, truyền mảng 1 phần tử
	ids, err := h.service.BatchCreate([]*model.Asset{&a})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.APIResponse{Code: 500, Status: "fail", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.APIResponse{
		Code:    201,
		Status:  "success",
		Message: "Tạo Asset thành công",
		Data:    map[string]interface{}{"id": ids[0]},
	})
}

// StartScan: Khởi chạy quá trình quét (Bài 1)
func (h *AssetHandler) StartScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)
	assetID := vars["id"]

	var req struct {
		ScanType string `json:"scan_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Mock target là localhost, thực tế sẽ lấy từ DB dựa trên assetID
	target := "127.0.0.1"

	job, err := h.scanService.StartScan(assetID, target, req.ScanType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{Code: 400, Status: "fail", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(job)
}

// GetScanResults: Lấy kết quả scan (Bài 1)
func (h *AssetHandler) GetScanResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Lấy Job ID từ URL để biết user đang hỏi kết quả của job nào
	vars := mux.Vars(r)
	jobID := vars["id"]

	// TẠO DỮ LIỆU MẪU ĐỂ TEST (Vì quá trình scan thật chạy ngầm có thể chưa kịp lưu DB)
	mockResult := map[string]interface{}{
		"job_id":    jobID,
		"status":    "completed",
		"scan_type": "port",
		"results": []map[string]interface{}{
			{
				"port":    80,
				"service": "http",
				"state":   "open",
				"version": "nginx/1.18.0",
			},
			{
				"port":    443,
				"service": "https",
				"state":   "open",
				"version": "openssl",
			},
		},
	}

	json.NewEncoder(w).Encode(mockResult)
}

func (h *AssetHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var req struct {
		Assets []*model.Asset `json:"assets"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{Code: 400, Status: "fail", Message: "Dữ liệu lỗi"})
		return
	}
	ids, err := h.service.BatchCreate(req.Assets)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.APIResponse{Code: 400, Status: "fail", Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.APIResponse{Code: 201, Status: "success", Message: "Tạo thành công", Data: map[string]interface{}{"created": len(ids), "ids": ids}})
}

func (h *AssetHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	stats, err := h.service.GetStats()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func (h *AssetHandler) CountAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	countResult, err := h.service.CountAssets(r.URL.Query().Get("type"), r.URL.Query().Get("status"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(countResult)
}

func (h *AssetHandler) BatchDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "thiếu tham số ids"})
		return
	}
	result, err := h.service.BatchDelete(strings.Split(idsParam, ","))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *AssetHandler) GetAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	res, err := h.service.GetAssets(page, limit, r.URL.Query().Get("type"), r.URL.Query().Get("status"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *AssetHandler) SearchAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	q := r.URL.Query().Get("q")
	if q == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "thiếu tham số q"})
		return
	}
	res, err := h.service.SearchAssets(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}
