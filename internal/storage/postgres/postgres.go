package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

// ConnectWithRetry 
func ConnectWithRetry(dsn string, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("🔄 Database connection attempt %d/%d...", attempt, maxRetries)

		// Thử mở kết nối
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("✅ Database connected successfully!")
				return db, nil
			}
		}

		if attempt < maxRetries {
			waitTime := time.Duration(1<<uint(attempt-1)) * time.Second
			log.Printf("⚠️ Connection failed: %v. Retrying in %v...", err, waitTime)
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("không thể kết nối database sau %d lần thử", maxRetries)
}

// NewPostgresStorage khởi tạo bộ lưu trữ PostgreSQL
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := ConnectWithRetry(dsn, 5)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

// Ping  (Health Check)
func (s *PostgresStorage) Ping() error {
	return s.db.Ping()
}

// InitTables tự động tạo bảng assets nếu chưa có
func (s *PostgresStorage) InitTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS assets (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		status VARCHAR(50) DEFAULT 'active',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("lỗi tạo bảng assets: %v", err)
	}
	log.Println("✅ Đã kiểm tra/khởi tạo bảng 'assets' thành công!")
	return nil
}