package main

import (
	"fmt"
	"log"
	"mini-asm/internal/handler"
	"mini-asm/internal/service"
	"mini-asm/internal/storage/postgres"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// 1. Lấy cấu hình từ Environment Variables (Để chạy được Docker Bài 5)
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "postgres")

	// Tạo chuỗi kết nối linh hoạt
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	log.Printf("📡 Đang kết nối tới Database tại: %s:%s...", dbHost, dbPort)

	// 2. Kết nối Database
	store, err := postgres.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatalf("❌ Server không thể khởi động: %v", err)
	}

	// Khởi tạo bảng nếu chưa có
	err = store.InitTables()
	if err != nil {
		log.Fatalf("❌ Lỗi tạo bảng: %v", err)
	}

	// 3. Khởi tạo Services & Handlers
	assetService := service.NewAssetService(store)
	scanService := service.NewScanService()
	
	alertStorage := postgres.NewPostgresAlertStorage(store.DB())
	alertService := service.NewAlertService(alertStorage, store)

	assetHandler := handler.NewAssetHandler(assetService, scanService)
	healthHandler := handler.NewHealthHandler(store)
	alertHandler := handler.NewAlertHandler(alertService)

	// 4. Định nghĩa Router
	router := mux.NewRouter()

	// --- Health Check Route ---
	router.HandleFunc("/health", healthHandler.Check).Methods("GET")

	// --- Assets Routes (Bài 1 & Bài 3) ---
	router.HandleFunc("/assets", assetHandler.GetAssets).Methods("GET")
	router.HandleFunc("/assets", assetHandler.CreateAsset).Methods("POST") // Thêm mới 1 asset
	router.HandleFunc("/assets/batch", assetHandler.BatchCreate).Methods("POST")
	router.HandleFunc("/assets/batch", assetHandler.BatchDelete).Methods("DELETE")
	router.HandleFunc("/assets/stats", assetHandler.GetStats).Methods("GET")
	router.HandleFunc("/assets/count", assetHandler.CountAssets).Methods("GET")
	router.HandleFunc("/assets/search", assetHandler.SearchAssets).Methods("GET")

	// --- Scan Routes (Bài 1) ---
	router.HandleFunc("/assets/{id}/scan", assetHandler.StartScan).Methods("POST")
	router.HandleFunc("/scan-jobs/{id}/results", assetHandler.GetScanResults).Methods("GET")

	// --- Alert Routes (Bài 6) ---
	alertHandler.RegisterRoutes(router)

	// 5. Chạy Server với CORS Middleware (Bài 3)
	// Bọc router bằng CORSMiddleware để Frontend cổng 3000 gọi được API cổng 8080
	port := ":8080"
	log.Printf("🚀 Server đang khởi chạy tại http://localhost%s ...", port)

	server := &http.Server{
		Addr:              port,
		Handler:           handler.CORSMiddleware(router),
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Hàm hỗ trợ lấy biến môi trường hoặc dùng giá trị mặc định
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
