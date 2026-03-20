# 🚀 Hệ thống Quản lý Tài sản EASM (External Attack Surface Management)

Dự án này là một hệ thống RESTful API kết hợp Frontend Web Dashboard để quản lý các tài sản IT và thực hiện truy quét lỗ hổng rủi ro bảo mật (EASM - External Attack Surface Management).

Dự án được phát triển và xây dựng trên lõi kiến trúc **Clean Architecture** với Go (Golang) và ReactJS.

## 📑 Mục lục
* [🌟 Các tính năng chính của hệ thống](#tinh-nang)
* [🛠️ Công nghệ & Kĩ thuật áp dụng](#cong-nghe)
* [📂 Cấu trúc dự án](#cau-truc)
* [📡 Danh sách API Endpoints](#api)
* [💻 Hướng dẫn Cài đặt & Khởi chạy](#cai-dat)
* [🧪 Hướng dẫn Kiểm thử & Triển khai (Demo Outputs)](#test)

---
<a id="tinh-nang"></a>
## 🌟 Các tính năng chính của hệ thống

Hệ thống được thiết kế toàn diện, từ Backend, Frontend, Testing đến các tiến trình DevOps:

1. **Quản lý Tài nguyên Cốt lõi (Base System):** Hỗ trợ chức năng CRUD, Thống kê, Phân trang, Tìm kiếm nội dung và tính năng Retry Backoff tự động. Triệt tiêu hoàn toàn rủi ro SQL Injection qua Parameterized Queries.
2. **Bài 1 - Mở rộng Scan API:** Động cơ EASM (EASM Scan Engine) chuyên sâu hỗ trợ quét IP (Geolocation & ASN), Port Scan thực tế (TCP Open Ports), giám sát chứng chỉ SSL và nhận diện công nghệ Tech Stack.
3. **Bài 2 - Tích hợp Unit Tests:** Các modules quét (Scanners) và Models được trang bị Unit Tests để đảm bảo độ tin cậy và đạt độ bao phủ Test (Coverage) tiêu chuẩn.
4. **Bài 3 - Tích hợp UI Frontend Dashboard:** Giao diện Web thiết kế hiện đại trên nền Vite + ReactJS, giao tiếp thời gian thực với API Backend cùng cấu hình CORS an toàn.
5. **Bài 4 - Đầu tư CI/CD Security Pipeline (GitHub Actions):** Tự động quét mã nguồn trực tiếp để rà soát rủi ro bảo mật qua các công cụ lớn (Gosec, Gitleaks, Trivy, TruffleHog) trước mỗi luồng Merge Request.
6. **Bài 5 - Infrastructure as Code (Docker Compose):** Đóng gói toàn vẹn Database (Postgres), Backend (Go) và Frontend (Vite) qua Multi-container Docker giúp khởi chạy bằng một lệnh duy nhất.
7. **Bài 6 - Tính năng EASM Mới (Alerts API):** Ghi nhận tự động các nguy cơ (Vulnerabilities/Alerts) từ kết quả EASM Scans và cung cấp Endpoints duyệt/lọc cảnh báo theo mức độ ưu tiên.
8. **Bài 7 - Cloud Deployment:** Máy chủ được triển khai Live Operation chạy liên tục trên DigitalOcean Cloud Droplet (Ubuntu 22.04 LTS).
9. **Bài 8 - Đăng ký Domain & TLS/HTTPS:** Cấu hình Webserver Nginx với chứng chỉ số Let's Encrypt bảo mật tuyệt đối qua HTTPS tại tên miền: `https://dungsocool-asm.duckdns.org`
10. **Bài 9 - Auto Deploy on Merge (CD Pipeline):** Máy chủ tự động SSH Pull Code và khởi chạy lại các containers sau khi nhánh main được cập nhật commit mới thông qua GitHub Actions.

<a id="cong-nghe"></a>
## 🛠️ Công nghệ & Thư viện sử dụng
* **Backend:** Go (Golang), `gorilla/mux`
* **Frontend:** ReactJS, Vite
* **Database:** PostgreSQL 15
* **Deploy/DevOps:** Docker, Docker Compose, GitHub Actions, Nginx + Certbot
* **Kiến trúc kĩ thuật:** Cấu trúc tầng Clean Architecture chuẩn mực: `Handler -> Service -> Storage`

<a id="cau-truc"></a>
## 📂 Cấu trúc dự án

Dự án phân rã chức năng theo quy mô chuẩn Clean Architecture để tối ưu tính Scale:

```text
├── .github/workflows/      # Cấu hình CI/CD GitHub Actions và Auto Deploy
├── cmd/server/             # Tệp Entry point nơi server khởi chạy (main.go)
├── frontend/               # Mã nguồn Giao diện UI Dashboard (ReactJS)
├── docker-compose.yml      # Cấu hình cài đặt đa môi trường Multi-container
├── internal/
│   ├── handler/            # Tầng HTTP Request/Response (REST API)
│   ├── model/              # Định nghĩa Structs Models và API Schema
│   ├── scanner/            # Modules gọi Tool Scan EASM (IP, Port, SSL, Tech)
│   ├── service/            # Tầng xử lý Logic trung tâm (Business Logic)
│   └── storage/            # Lớp tương tác với Cơ sở dữ liệu (Database layer)
└── README.md
```

<a id="api"></a>
## 📡 Danh sách API Endpoints

| Method | Endpoint | Tính năng |
| :--- | :--- | :--- |
| `GET` | `/health` | Kiểm tra trạng thái của Server DB Connect |
| `POST` | `/assets` | Tạo mới một tài sản |
| `GET` | `/assets` | Lấy danh sách tài sản (Query: `?page=1&limit=20`) |
| `POST` | `/assets/batch` | Tạo mới nhiều tài sản cùng lúc |
| `DELETE` | `/assets/batch` | Xóa nhiều tài sản đồng loạt (Query: `?ids=id1,id2`) |
| `GET` | `/assets/stats` | Thống kê số lượng tổng tài sản hiện có |
| `GET` | `/assets/search` | Tìm kiếm tên tài sản gần đúng (Query: `?q=text`) |
| `POST` | `/assets/{id}/scan` | Khởi tạo chạy công cụ quét Scan EASM theo IP/Port/SSL |
| `GET` | `/scan-jobs/{id}/results` | Truy xuất lấy dữ liệu kết quả phân tích sau khi đã quét xong |
| `GET` | `/alerts` | Lấy danh sách tổng hợp tất cả cảnh báo bảo mật EASM |
| `GET` | `/assets/{id}/alerts` | Lấy chi tiết cảnh báo thuộc về một tài sản nhất định |

<a id="cai-dat"></a>
## 💻 Hướng dẫn Cài đặt & Khởi chạy

### Cách 1: Chạy chuẩn hóa bằng Docker Compose (Khuyên dùng)
Hãy cài đặt Docker Compose vào máy hoặc Server VM của bạn và chạy lệnh dưới đây tại thư mục gốc của project:
```bash
docker compose up -d --build
```
Webserver của kiến trúc sẽ bắt đầu phục vụ ở cổng cục bộ:
- **Frontend Dashboard:** `http://localhost:3000` (Local) hoặc `https://dungsocool-asm.duckdns.org` (Production)
- **Backend API:** `http://localhost:8080` (Local) hoặc `https://dungsocool-asm.duckdns.org/api` (Production)

---

<a id="test"></a>
## 🧪 Hướng dẫn Kiểm thử & Triển khai (Demo Outputs)

### 1. Trạng thái Triển khai Cloud (Bài 7 & 8)
Hệ thống hiện đã được cấu hình chạy Production Ready qua Nginx Reverse Proxy với HTTPS bảo mật:
```text
URL Truy cập: https://dungsocool-asm.duckdns.org
SSL Issuer: Let's Encrypt Authority R3
Server IP: 159.223.60.128
```

### 2. Kiểm thử Docker Container Status (Bài 5)
```text
root@ubuntu:~/project# docker compose ps
NAME           IMAGE                COMMAND                  SERVICE    STATUS             PORTS
cmc_backend    dev-backend          "./main"                 backend    Up 19 minutes      0.0.0.0:8080->8080/tcp
cmc_frontend   dev-frontend         "/docker-entrypoint…"    frontend   Up 19 minutes      0.0.0.0:3000->80/tcp
cmc_postgres   postgres:15-alpine   "docker-entrypoint.s…"   db         Up 4 hrs (healthy) 0.0.0.0:5432->5432/tcp
```

### 3. Kiểm thử Unit Tests Code Coverage (Bài 2)
```text
PS C:\Users\xxx\DEV> go test -cover ./...
ok      mini-asm/internal/model         0.036s  coverage: 100.0% of statements
ok      mini-asm/internal/scanner       7.051s  coverage: 55.3% of statements
```

### 4. Kiểm thử HTTPS API Health Check
```json
> curl -s https://dungsocool-asm.duckdns.org/api/health
{
  "status":"ok",
  "database":"connected",
  "timestamp":"2026-03-20T21:20:00Z"
}
```

### 5. Kiểm thử CI/CD & Auto Deploy Pipeline (Bài 4 & 9)
Khi mã nguồn được đẩy lên nhánh `main`, tiến trình GitHub Actions sẽ tự rà soát mã nguồn (Gosec, Gitleaks) và thực hiện SSH Deploy lên Droplet:
```text
Workflow: Deploy to Production
Event: push to main
Status: Success ✅ (Build -> Test -> Security Check -> Deploy)
```

---
*Dự án được bảo mật bởi cơ chế Gosec và chống SQL Injection.*
