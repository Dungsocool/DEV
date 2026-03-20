package storage

import "mini-asm/internal/model"

type Storage interface {
	Ping() error
	BatchCreate(assets []*model.Asset) error
	GetStats() (*model.AssetStats, error)
	CountAssets(assetType string, status string) (int, error)
	BatchDelete(ids []string) (int, error)
	
	// Thêm 2 hàm Bonus
	GetAssets(page, limit int, assetType, status string) ([]*model.Asset, int, error)
	SearchAssets(query string) ([]*model.Asset, error)
}