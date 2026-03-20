package main

import (
	"log"
	"mini-asm/internal/handler"
	"mini-asm/internal/service"
	"mini-asm/internal/storage/postgres"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// 1. Kết nối Database
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	store, err := postgres.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatalf("❌ Server không thể khởi động: %v", err)
	}

	err = store.InitTables()
	if err != nil {
		log.Fatalf("❌ Lỗi tạo bảng: %v", err)
	}

	// 2. Khởi tạo Services
	assetService := service.NewAssetService(store)
	scanService := service.NewScanService() // Dòng 28: Khai báo scanService

	// 2. Khởi tạo Handler
	// SỬA DÒNG NÀY: Truyền CẢ HAI tham số (assetService VÀ scanService) vào hàm
	assetHandler := handler.NewAssetHandler(assetService, scanService)

	healthHandler := handler.NewHealthHandler(store)

	// 4. Định nghĩa Router
	router := mux.NewRouter()

	// --- Các Route cũ của bạn ---
	router.HandleFunc("/health", healthHandler.Check).Methods("GET")
	router.HandleFunc("/assets/batch", assetHandler.BatchCreate).Methods("POST")
	router.HandleFunc("/assets/stats", assetHandler.GetStats).Methods("GET")
	router.HandleFunc("/assets/count", assetHandler.CountAssets).Methods("GET")
	router.HandleFunc("/assets/batch", assetHandler.BatchDelete).Methods("DELETE")
	router.HandleFunc("/assets/search", assetHandler.SearchAssets).Methods("GET")
	router.HandleFunc("/assets", assetHandler.GetAssets).Methods("GET")

	// --- CẬP NHẬT CÁC ROUTE MỚI CHO BÀI 1 & BÀI 3 ---

	// Cho phép POST /assets để tạo 1 asset (Sửa lỗi 405 của bạn)
	// Lưu ý: Đảm bảo trong asset_handler.go đã có hàm CreateAsset
	router.HandleFunc("/assets", assetHandler.CreateAsset).Methods("POST")

	// Route để bắt đầu Scan (Bài 1)
	router.HandleFunc("/assets/{id}/scan", assetHandler.StartScan).Methods("POST")

	// Route để xem kết quả Scan (Bài 1)
	// Giả sử bạn đặt hàm này trong assetHandler
	router.HandleFunc("/scan-jobs/{id}/results", assetHandler.GetScanResults).Methods("GET")

	// 5. Chạy Server với CORS Middleware (Bài 3)
	// Bọc router bằng CORSMiddleware để Frontend gọi được API
	log.Println("🚀 Server đang chạy tại cổng http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", handler.CORSMiddleware(router)); err != nil {
		log.Fatal(err)
	}
}
