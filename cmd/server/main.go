package main

import (
	"log"
	"net/http"
	"mini-asm/internal/handler"
	"mini-asm/internal/service"
	"mini-asm/internal/storage/postgres"
	"github.com/gorilla/mux"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	store, err := postgres.NewPostgresStorage(dsn)
	if err != nil { log.Fatalf("❌ Server không thể khởi động: %v", err) }
	err = store.InitTables()
	if err != nil { log.Fatalf("❌ Lỗi tạo bảng: %v", err) }

	assetService := service.NewAssetService(store)
	assetHandler := handler.NewAssetHandler(assetService)
	healthHandler := handler.NewHealthHandler(store)

	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler.Check).Methods("GET")
	router.HandleFunc("/assets/batch", assetHandler.BatchCreate).Methods("POST")
	router.HandleFunc("/assets/stats", assetHandler.GetStats).Methods("GET")
	router.HandleFunc("/assets/count", assetHandler.CountAssets).Methods("GET")
	router.HandleFunc("/assets/batch", assetHandler.BatchDelete).Methods("DELETE")
	
	// Khai báo 2 API
	router.HandleFunc("/assets/search", assetHandler.SearchAssets).Methods("GET") 
	router.HandleFunc("/assets", assetHandler.GetAssets).Methods("GET")

	log.Println("🚀 Server đang chạy tại cổng http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", router); err != nil { log.Fatal(err) }
}