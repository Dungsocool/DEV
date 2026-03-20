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

| Thành phần | Công nghệ |
| :--- | :--- |
| Backend | Go (Golang), gorilla/mux |
| Frontend | ReactJS, Vite |
| Database | PostgreSQL 15 |
| Containerization | Docker, Docker Compose |
| CI/CD | GitHub Actions (ci.yml + deploy.yml) |
| Security Scanning | Gosec, Gitleaks, Trivy, TruffleHog |
| Web Server | Nginx + Certbot (Let's Encrypt) |
| Cloud | DigitalOcean Droplet (Ubuntu 22.04 LTS) |
| DNS | DuckDNS (Free Dynamic DNS) |
| Kiến trúc | Clean Architecture: Handler → Service → Storage |

<a id="cau-truc"></a>
## 📂 Cấu trúc dự án

Dự án phân rã chức năng theo quy mô chuẩn Clean Architecture để tối ưu tính Scale:

```text
├── .github/
│   └── workflows/
│       ├── ci.yml              # Pipeline CI bảo mật (Gosec, Gitleaks, Trivy, TruffleHog)
│       └── deploy.yml          # Pipeline CD tự động deploy lên Droplet khi merge
├── cmd/
│   └── server/
│       └── main.go             # Entry point: khởi tạo DB, Router, Services, Server
├── frontend/
│   ├── index.html              # Trang chủ Dashboard ReactJS
│   └── Dockerfile              # Build image Nginx phục vụ static files
├── internal/
│   ├── handler/
│   │   ├── asset_handler.go    # REST API xử lý CRUD Assets + Scan
│   │   ├── alert_handler.go    # REST API xử lý Alerts (Bài 6)
│   │   ├── health_handler.go   # Health check endpoint
│   │   └── middleware.go       # CORS Middleware cho Frontend
│   ├── model/
│   │   ├── asset.go            # Struct Asset, APIResponse, ScanJob
│   │   ├── alert.go            # Struct Alert, AlertStats
│   │   └── alert_test.go       # Unit tests cho Alert model
│   ├── scanner/
│   │   ├── ip_scanner.go       # Module quét IP (Geolocation, ASN)
│   │   ├── port_scanner.go     # Module quét Port (TCP Open Ports)
│   │   ├── ssl_scanner.go      # Module kiểm tra chứng chỉ SSL
│   │   ├── tech_scanner.go     # Module nhận diện công nghệ
│   │   └── *_test.go           # Unit tests cho các scanner
│   ├── service/
│   │   ├── asset_service.go    # Business logic quản lý Assets
│   │   ├── scan_service.go     # Business logic điều phối Scan
│   │   └── alert_service.go    # Business logic quản lý Alerts
│   └── storage/
│       ├── storage.go          # Interface Storage (Clean Architecture)
│       └── postgres/
│           ├── postgres.go     # Kết nối DB với Retry Backoff
│           ├── alert_storage.go # CRUD Alerts trong PostgreSQL
│           ├── delete_storage.go # Batch Delete với Parameterized Query
│           └── migrations/     # SQL migration files
├── docs/images/                # Ảnh minh chứng các bài tập
├── docker-compose.yml          # Cấu hình 3 containers: backend, frontend, db
├── Dockerfile                  # Multi-stage build cho Go backend
└── README.md
```

<a id="api"></a>
## 📡 Danh sách API Endpoints

### Assets API (Quản lý tài sản)

| Method | Endpoint | Tính năng | Ví dụ |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Kiểm tra trạng thái Server + DB | `curl localhost:8080/health` |
| `POST` | `/assets` | Tạo mới một tài sản | Body: `{"name":"Server","type":"ip","value":"1.1.1.1"}` |
| `GET` | `/assets` | Lấy danh sách tài sản (phân trang) | `?page=1&limit=20` |
| `POST` | `/assets/batch` | Tạo mới nhiều tài sản cùng lúc | Body: `[{...},{...}]` |
| `DELETE` | `/assets/batch` | Xóa nhiều tài sản đồng loạt | `?ids=id1,id2` |
| `GET` | `/assets/stats` | Thống kê số lượng tổng tài sản | - |
| `GET` | `/assets/search` | Tìm kiếm tài sản gần đúng | `?q=server` |

### Scan API (Quét bảo mật EASM)

| Method | Endpoint | Tính năng | Ví dụ |
| :--- | :--- | :--- | :--- |
| `POST` | `/assets/{id}/scan` | Khởi tạo quét Scan EASM | Body: `{"scan_type":"port"}` |
| `GET` | `/scan-jobs/{id}/results` | Lấy kết quả phân tích sau khi quét | - |

### Alerts API (Cảnh báo bảo mật - Bài 6)

| Method | Endpoint | Tính năng |
| :--- | :--- | :--- |
| `GET` | `/alerts` | Lấy danh sách tất cả cảnh báo (phân trang, lọc theo severity) |
| `GET` | `/assets/{id}/alerts` | Lấy cảnh báo thuộc về một tài sản cụ thể |
| `GET` | `/alerts/stats` | Thống kê số lượng cảnh báo theo loại và mức độ |
| `PUT` | `/alerts/{id}/status` | Cập nhật trạng thái cảnh báo (open/resolved/dismissed) |

<a id="cai-dat"></a>
## 💻 Hướng dẫn Cài đặt & Khởi chạy

### Yêu cầu hệ thống
- Docker & Docker Compose (v2+)
- Go 1.21+ (nếu chạy Dev Mode)
- Node.js 18+ (nếu chạy Dev Mode)

### Cách 1: Chạy bằng Docker Compose (Khuyên dùng)
Chỉ cần 1 lệnh duy nhất tại thư mục gốc của project:
```bash
docker compose up -d --build
```
Hệ thống sẽ tự động:
- Khởi tạo PostgreSQL 15 với health check
- Build backend Go binary từ multi-stage Dockerfile
- Build frontend và phục vụ qua Nginx
- Kết nối các service với nhau qua Docker network

Truy cập:
- Frontend Dashboard: `http://localhost:3000` (Local) hoặc `https://dungsocool-asm.duckdns.org` (Production)
- Backend API: `http://localhost:8080` (Local) hoặc `https://dungsocool-asm.duckdns.org/api` (Production)

### Cách 2: Chạy chế độ Phát triển (Dev Mode)
```bash
# 1. Bật Database
docker compose up -d db

# 2. Khởi chạy Backend (terminal 1)
go run cmd/server/main.go

# 3. Khởi chạy Frontend (terminal 2)
cd frontend
npm install
npm run dev
```

### Biến môi trường hỗ trợ (Environment Variables)
| Biến | Mặc định | Mô tả |
| :--- | :--- | :--- |
| `DB_HOST` | `localhost` | Host của PostgreSQL |
| `DB_PORT` | `5432` | Port của PostgreSQL |
| `DB_USER` | `postgres` | Username kết nối DB |
| `DB_PASSWORD` | `postgres` | Password kết nối DB |
| `DB_NAME` | `postgres` | Tên database |

---

<a id="test"></a>
## 🧪 Hướng dẫn Kiểm thử & Triển khai (Demo Outputs)

### 1. Trạng thái Triển khai Cloud & HTTPS (Bài 7 & 8)
Hệ thống hiện đã được cấu hình chạy Production Ready qua Nginx Reverse Proxy với HTTPS bảo mật:
```text
URL Truy cập:  https://dungsocool-asm.duckdns.org
SSL Issuer:    Let's Encrypt Authority R3
Server IP:     159.223.60.128
Cloud:         DigitalOcean Droplet (Ubuntu 22.04 LTS, 512MB RAM, SGP1)
```
![Droplet DigitalOcean đang chạy](docs/images/image.png)
![HTTPS ổ khóa xanh trên trình duyệt](docs/images/image-1.png)

### 2. Kiểm thử Docker Container Status (Bài 5)
```text
NAME           IMAGE                COMMAND     SERVICE    STATUS             PORTS
cmc_backend    dev-backend          "./main"    backend    Up 19 minutes      0.0.0.0:8080->8080/tcp
cmc_frontend   dev-frontend         "nginx..."  frontend   Up 19 minutes      0.0.0.0:3000->80/tcp
cmc_postgres   postgres:15-alpine   "docker…"   db         Up 4 hrs (healthy) 0.0.0.0:5432->5432/tcp
```
![Docker Compose containers đang Up](docs/images/image-3.png)

### 3. Kiểm thử Unit Tests Code Coverage (Bài 2)
```text
$ go test -cover ./...
ok      mini-asm/internal/model         0.036s  coverage: 100.0% of statements
ok      mini-asm/internal/scanner       7.051s  coverage: 55.3% of statements
```
![Unit Tests Coverage](docs/images/image-2.png)

### 4. Kiểm thử HTTPS API Health Check (Bài 7 & 8)
```json
$ curl.exe -s https://dungsocool-asm.duckdns.org/api/health
{"database":{"status":"connected"},"status":"ok","timestamp":"2026-03-20T14:46:24Z"}
```
![API Health Check qua HTTPS](docs/images/image-4.png)

### 5. Kiểm thử CI/CD & Auto Deploy Pipeline (Bài 4 & 9)
Khi mã nguồn được đẩy lên nhánh `main`, tiến trình GitHub Actions sẽ tự rà soát mã nguồn (Gosec, Gitleaks, Trivy, TruffleHog) và thực hiện SSH Deploy lên Droplet:
```text
Workflow:  Deploy to Production / mini-asm Security CI
Pipeline:  Build → Test → Security Check → Deploy → Post-Deploy Verification
Status:    Success ✅
```
![GitHub Actions CI/CD Pipeline](docs/images/image-6.png)
![GitHub Actions Workflow Steps](docs/images/image-7.png)

### 6. Kiểm thử Giao diện Frontend UI Dashboard (Bài 3)
Mở trình duyệt và truy cập `http://localhost:3000` hoặc `https://dungsocool-asm.duckdns.org`. Dashboard hiển thị danh sách Assets, cho phép tạo mới, chạy Quick Scan (IP/Port/SSL/Tech) và xem kết quả trực tiếp trên giao diện.
![Frontend UI Dashboard](docs/images/image-5.png)

---
*Dự án được bảo mật bởi cơ chế quét Gosec, Gitleaks, Trivy, TruffleHog và chống SQL Injection qua Parameterized Queries.*
