package service

import (
	"errors"
	"time"
	"mini-asm/internal/model"
	"mini-asm/internal/storage"
	"github.com/google/uuid"
)

type AssetService struct { store storage.Storage }
func NewAssetService(store storage.Storage) *AssetService { return &AssetService{store: store} }

func (s *AssetService) BatchCreate(assets []*model.Asset) ([]string, error) {
	if len(assets) == 0 { return nil, errors.New("danh sách trống") }
	if len(assets) > 100 { return nil, errors.New("tối đa 100 tài sản") }
	var ids []string
	now := time.Now()
	for _, asset := range assets {
		if asset.Name == "" { return nil, errors.New("tên không được trống") }
		if asset.Type != model.TypeDomain && asset.Type != model.TypeIP && asset.Type != model.TypeService { return nil, errors.New("loại không hợp lệ") }
		asset.ID = uuid.New().String()
		asset.Status = model.StatusActive
		asset.CreatedAt = now
		asset.UpdatedAt = now
		ids = append(ids, asset.ID)
	}
	err := s.store.BatchCreate(assets)
	if err != nil { return nil, err }
	return ids, nil
}

func (s *AssetService) GetStats() (*model.AssetStats, error) { return s.store.GetStats() }

func (s *AssetService) CountAssets(assetType, status string) (*model.AssetCount, error) {
	count, err := s.store.CountAssets(assetType, status)
	if err != nil { return nil, err }
	return &model.AssetCount{Count: count, Filters: map[string]string{"type": assetType, "status": status}}, nil
}

func (s *AssetService) BatchDelete(ids []string) (map[string]int, error) {
	if len(ids) == 0 { return nil, errors.New("danh sách trống") }
	deleted, err := s.store.BatchDelete(ids)
	if err != nil { return nil, err }
	notFound := len(ids) - deleted
	if notFound < 0 { notFound = 0 }
	return map[string]int{"deleted": deleted, "not_found": notFound}, nil
}

// Hàm Bonus Bài 6
func (s *AssetService) GetAssets(page, limit int, assetType, status string) (*model.PaginatedResponse, error) {
	if page < 1 { page = 1 }
	if limit < 1 { limit = 20 }
	if limit > 100 { limit = 100 }

	assets, total, err := s.store.GetAssets(page, limit, assetType, status)
	if err != nil { return nil, err }

	totalPages := total / limit
	if total%limit != 0 { totalPages++ }

	return &model.PaginatedResponse{
		Data: assets,
		Pagination: model.Pagination{Page: page, Limit: limit, Total: total, TotalPages: totalPages},
	}, nil
}

// Hàm Bonus Bài 7
func (s *AssetService) SearchAssets(query string) ([]*model.Asset, error) {
	if query == "" { return []*model.Asset{}, nil }
	return s.store.SearchAssets(query)
}