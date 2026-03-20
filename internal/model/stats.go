package model

type AssetStats struct {
	Total    int            `json:"total"`
	ByType   map[string]int `json:"by_type"`
	ByStatus map[string]int `json:"by_status"`
}

type AssetCount struct {
	Count   int               `json:"count"`
	Filters map[string]string `json:"filters"`
}