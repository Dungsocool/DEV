package scanner

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"mini-asm/internal/model"
)

type IPScanner struct {
	client *http.Client
}

func NewIPScanner() *IPScanner {
	return &IPScanner{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *IPScanner) Scan(targetIP string) (interface{}, error) {
	// 1. Gọi API miễn phí để lấy vị trí địa lý của IP
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,country,countryCode,regionName,city,lat,lon,isp,org,as,asname", targetIP)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi gọi API: %v", err)
	}
	defer resp.Body.Close()

	var apiRes map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&apiRes)

	if apiRes["status"] == "fail" {
		return nil, fmt.Errorf("ip-api error: %v", apiRes["message"])
	}

	// 2. Tìm Reverse DNS
	ptr, _ := net.LookupAddr(targetIP)
	reverseDNS := ""
	if len(ptr) > 0 {
		reverseDNS = ptr[0]
	}

	// 3. Đóng gói dữ liệu trả về theo Model
	result := model.IPScanResult{
		IPAddress: targetIP,
		Geolocation: model.Geolocation{
			Country:     fmt.Sprint(apiRes["country"]),
			CountryCode: fmt.Sprint(apiRes["countryCode"]),
			City:        fmt.Sprint(apiRes["city"]),
			Region:      fmt.Sprint(apiRes["regionName"]),
			ISP:         fmt.Sprint(apiRes["isp"]),
			Org:         fmt.Sprint(apiRes["org"]),
		},
		ASN: model.ASNInfo{
			Name:        fmt.Sprint(apiRes["asname"]),
			Description: fmt.Sprint(apiRes["as"]),
		},
		ReverseDNS: reverseDNS,
		CreatedAt:  time.Now(),
	}

	return []model.IPScanResult{result}, nil
}
