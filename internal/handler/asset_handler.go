package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"mini-asm/internal/model"
	"mini-asm/internal/service"
)

type AssetHandler struct { service *service.AssetService }
func NewAssetHandler(s *service.AssetService) *AssetHandler { return &AssetHandler{service: s} }

func (h *AssetHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var req struct { Assets []*model.Asset `json:"assets"` }
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
	if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
	json.NewEncoder(w).Encode(stats)
}

func (h *AssetHandler) CountAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	countResult, err := h.service.CountAssets(r.URL.Query().Get("type"), r.URL.Query().Get("status"))
	if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
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
	if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Bài 6
func (h *AssetHandler) GetAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	
	res, err := h.service.GetAssets(page, limit, r.URL.Query().Get("type"), r.URL.Query().Get("status"))
	if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
	json.NewEncoder(w).Encode(res)
}

// Bài 7
func (h *AssetHandler) SearchAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	q := r.URL.Query().Get("q")
	if q == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "thiếu tham số q"})
		return
	}
	res, err := h.service.SearchAssets(q)
	if err != nil { w.WriteHeader(http.StatusInternalServerError); return }
	json.NewEncoder(w).Encode(res)
}