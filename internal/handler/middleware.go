package handler

import "net/http"

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cho phép tất cả các nguồn (Origin) truy cập
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Cho phép các phương thức HTTP
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Cho phép các Header cần thiết
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Xử lý yêu cầu kiểm tra (Pre-flight request) của trình duyệt
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
