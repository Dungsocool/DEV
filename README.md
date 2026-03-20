# 🚀 CMC Intern API Project - Asset Management & EASM

Đây là dự án xây dựng hệ thống RESTful API dùng để quản lý các tài sản IT (Domain, IP, Service) và hệ thống quét tài nguyên (EASM - External Attack Surface Management). Dự án được phát triển bằng Go (Golang), ReactJS (Frontend) và cơ sở dữ liệu PostgreSQL. Dự án tuân thủ nghiêm ngặt mô hình **Clean Architecture** để tối ưu hóa việc bảo trì và phân tách rõ ràng các nghiệp vụ.

## 📑 Mục lục (Bài Tập Về Nhà - Sessions 5-7)
* [✅ Tiến độ Bài Tập (Day 3)](#tinh-nang)
* [🛠️ Công nghệ & Thư viện sử dụng](#cong-nghe)
* [📂 Cấu trúc thư mục (Project Structure)](#cau-truc)
* [📡 Danh sách API Endpoints](#api)
* [💻 Hướng dẫn Cài đặt & Khởi chạy (How to Run)](#cai-dat)
* [🧪 Hướng Dẫn Test Từng Bài (API Testing Guide)](#test)
* [💾 Bài tập cũ (Day 1 & 2)](#day1-2)

---
<a id="cong-nghe"></a>
## 🛠️ Công nghệ & Thư viện sử dụng
* **Backend:** Go (Golang)
* **Frontend:** ReactJS, Vite
* **Cơ sở dữ liệu:** PostgreSQL
* **Bộ định tuyến (Router):** `gorilla/mux`
* **Hạ tầng & Triển khai:** Docker, Docker Compose, GitHub Actions (CI/CD)
* **Kiến trúc:** Clean Architecture (Handler -> Service -> Storage)

<a id="bao-mat"></a>
## 🛡️ Tiêu điểm Bảo mật: Ngăn chặn hoàn toàn SQL Injection (SQLi)

Trong dự án này, rủi ro SQL Injection đã được triệt tiêu 100% tại tầng Storage bằng cách áp dụng triệt để kỹ thuật **Parameterized Queries** (Truy vấn tham số hóa) thông qua driver `database/sql` của Go. 

Nguyên lý hoạt động (Tại sao nó an toàn?): Khi một truy vấn được gửi đi, quá trình diễn ra qua 2 bước bảo mật khép kín:
**Bước 1 (Prepare):** Database nhận cấu trúc câu lệnh SQL với các "chỗ trống" ảo. 
**Bước 2 (Execute):** Dữ liệu người dùng nhập vào mới được gửi xuống để "lấp" vào các chỗ trống đó dưới dạng **Văn bản thuần túy**.

*(Xem thêm ảnh minh họa ở phần Test Day 1 & 2).*

<a id="tinh-nang"></a>
## ✅ Các tính năng đã hoàn thành

### Giai đoạn 1 & 2 (Day 1-2):
1. **Bài 1 - Statistics APIs:** Cung cấp API thống kê tổng tài sản và đếm số lượng tài sản.
2. **Bài 2 - Batch Create:** Thêm mới hàng loạt tài sản (tối đa 100 items/request).
3. **Bài 3 - Batch Delete:** Xóa hàng loạt tài sản thông minh dựa trên danh sách ids.
4. **Bài 4 - Connection Retry:** Thuật toán **Exponential Backoff** tự động thử lại kết nối DB.
5. **Bài 5 - Health Check:** Cung cấp API giám sát tình trạng hệ thống.
6. **Bài 6 - Pagination & Filtering:** Xử lý phân trang và lọc dữ liệu động.
7. **Bài 7 - Search by Name:** Tìm kiếm tài sản theo tên (Partial match).

### Giai đoạn 3 (Day 3) - EASM & Deployment (Hoàn thành 6/9 bài):
1. **Bài 1 - Mở rộng Scan API:** Hỗ trợ tính năng Web/Port Scan qua Giao thức API.
2. **Bài 2 - Viết Unit Tests:** Hệ thống Test nội bộ cho Model dự án và Scanner tool với Coverage cao.
3. **Bài 3 - Tích hợp Frontend:** Tích hợp giao diện UI Web Dashboard bằng ReactJS để thao tác trực quan.
4. **Bài 4 - CI/CD GitHub Actions:** Pipeline kiểm tra bảo mật (gosec, gitleaks, trivy, trufflehog).
5. **Bài 5 - Deploy Docker Compose:** Đóng gói toàn bộ Frontend, Backend, Database bằng Container.
6. **Bài 6 - Tính năng EASM Mới:** Hệ thống Cảnh báo (Alerts API) - Quản lý và thống kê các cảnh báo bảo mật.

### 🌟 Bonus (Đang chờ - Day 3):
- **Bài 7 - Deploy Cloud VM:** Đưa ứng dụng lên Server thực tế.
- **Bài 8 - Domain & TLS/HTTPS:** Cấu hình tên miền và chứng chỉ SSL.
- **Bài 9 - Auto Deploy On Merge:** Tự động đẩy code lên Server khi Merge nhánh Main.

<a id="cau-truc"></a>
## 📂 Cấu trúc thư mục (Project Structure)
Dự án được chia thành các package nhỏ lẻ tuân theo Clean Architecture:

```text
├── .github/workflows/      # Cấu hình CI/CD GitHub Actions
├── cmd/server/             # Entry point (main.go)
├── frontend/               # Mã nguồn UI (ReactJS)
├── homeworks/              # Chứa các báo cáo nộp bài (SUBMISSION.md)
├── internal/
│   ├── handler/            # Tầng giao tiếp HTTP
│   ├── model/              # Định nghĩa các cấu trúc dữ liệu
│   ├── scanner/            # Chứa các logic gọi Tool Scan EASM (Bài 1, Day 3)
│   ├── service/            # Tầng xử lý Logic nghiệp vụ
│   └── storage/            # Tầng giao tiếp Cơ sở dữ liệu (PostgreSQL)
├── docker-compose.yml      # Cấu hình cài đặt multi-container
└── README.md
```

<a id="api"></a>
## 📡 Danh sách API Endpoints

| Method | Endpoint | Tính năng |
| :--- | :--- | :--- |
| `GET` | `/health` | Kiểm tra trạng thái của Server |
| `POST` | `/assets` | Tạo mới 1 tài sản |
| `GET` | `/assets` | Lấy danh sách tài sản (Query: `?page=1&limit=20`) |
| `GET` | `/assets/stats` | Thống kê số lượng tài sản theo loại |
| `GET` | `/assets/search` | Tìm kiếm tài sản theo tên (Query: `?q=keyword`) |
| `POST` | `/assets/batch` | Tạo mới nhiều tài sản cùng lúc |
| `DELETE` | `/assets/batch` | Xóa nhiều tài sản (Query: `?ids=1,2`) |
| **`POST`** | **`/assets/{id}/scan`** | **Khởi tạo chạy công cụ quét Scan (IP/Port)** |
| **`GET`** | **`/scan-jobs/{id}/results`** | **Lấy dữ liệu kết quả sau khi quét xong** |
| `GET` | `/alerts` | Lấy danh sách cảnh báo bảo mật (Bài 6) |
| `GET` | `/assets/{id}/alerts` | Lấy cảnh báo theo một tài sản cụ thể |

<a id="cai-dat"></a>
## 💻 Hướng dẫn Cài đặt & Khởi chạy (How to Run)

### Cách 1: Chạy bằng Docker Compose (Khuyên dùng)
Đảm bảo máy bạn đã cài Docker Desktop. Chạy lệnh:
```bash
docker compose up -d --build
```
Dự án sẽ khởi chạy tại:
- **Frontend Dashboard:** `http://localhost:3000`
- **Backend API:** `http://localhost:8080`

### Cách 2: Chạy độc lập (Local/Dev)
1. **Khởi chạy Database:**
```bash
docker compose up -d db
```
2. **Khởi chạy API Server:**
```bash
go run cmd/server/main.go
```
3. **Khởi chạy Frontend:** Mở terminal mới:
```bash
cd frontend
npm install
npm run dev
```

<a id="test"></a>
## 🧪 Hướng Dẫn Test Từng Bài (API Testing Guide)

*(Dưới đây là phần minh chứng kiểm thử các chức năng của toàn bộ dự án)*

### 🚀 [PHẦN MỚI] BÀI TẬP DAY 3 (Sessions 5-7)

#### Bài 1: Mở rộng Scan API
Gửi một POST request tới hệ thống để quét Port của 1 IP Asset.
```bash
curl.exe -s -X POST http://localhost:8080/assets/<IP-ASSET-ID>/scan -H "Content-Type: application/json" -d '{"scan_type":"port"}'
```
*(Chèn ảnh kết quả Test Scan API bằng Postman/Curl vào đây)*
![Test Scan API](docs/day3-bai1-scan.png)

#### Bài 2: Viết Unit Tests
Chạy lệnh `go test -v -coverprofile=coverage.out ./...` để kiểm tra toàn bộ ứng dụng. 
*(Chèn ảnh Terminal hiện "PASS" xanh và mức độ test coverage vào đây)*
![Unit Tests Result](docs/day3-bai2-unittest.png)

#### Bài 3: Tích hợp Frontend
Dashboard ở port `3000` cho phép dễ dàng tạo Asset và Xem kết quả Port Scan.
*(Chèn ảnh giao diện trình duyệt web đang hiển thị kết quả thao tác vào đây)*
![Frontend UI](docs/day3-bai3-frontend.png)

#### Bài 4: CI/CD với GitHub Actions
Luồng CI/CD được kích hoạt tự động mỗi khi có Push/Merge với các công cụ bảo mật Gosec, Gitleaks, v.v.
*(Chèn ảnh màn hình workflow ticks xanh Github Actions vào đây)*
![Github Actions CI/CD](docs/day3-bai4-cicd.png)

#### Bài 5: Deploy với Docker Compose
Kiểm tra trạng thái các container bằng `docker compose ps`:
*(Chèn ảnh màn hình Docker Desktop hoặc Terminal hiển thị container đang "Up" vào đây)*
![Docker Compose Deployment](docs/day3-bai5-docker.png)

#### Bài 6 (Tính năng EASM mới)
Hệ thống Alerts API cảnh báo bảo mật. Bạn có thể kiểm tra danh sách qua endpoint `/alerts`.
*(Chèn ảnh kết quả test Tính năng EASM mới)*
![Bonus Bài 6](docs/day3-bai6-bonus.png)

#### Bài 7 -> 9: Các tính năng nâng cao (Bonus)
*(Chưa hoàn thiện)*

---
<a id="day1-2"></a>
### 💾 [PHẦN CŨ] BÀI TẬP DAY 1 & 2

#### [Bài 1] Thống kê & Đếm tài sản (Statistics & Count )
```bash
curl.exe -X GET http://localhost:8080/assets/stats
```
<img width="697" height="67" alt="image" src="https://github.com/user-attachments/assets/50f12b73-ebc5-4092-9122-a6e68b4cf119" />

#### [Bài 2] Thêm hàng loạt tài sản (Batch Create)
<img width="1578" height="96" alt="image" src="https://github.com/user-attachments/assets/c5e77937-ee30-4bb3-a455-bc212485a99a" />

#### [Bài 3] Xóa hàng loạt (Batch Delete)
```bash
curl.exe -X DELETE "http://localhost:8080/assets/batch?ids=id-1,id-2,id-3"
```
<img width="1069" height="56" alt="image" src="https://github.com/user-attachments/assets/d874ad13-2cc2-46b4-b765-3a69e54f903e" />

#### [Bài 4 & 5] Thuật toán Retry & Giám sát sức khỏe (Health Check)
```bash
curl.exe -X GET http://localhost:8080/health
```
<img width="799" height="81" alt="image" src="https://github.com/user-attachments/assets/783af0ba-404c-4da2-b732-c45142297ef2" />

Test Retry:
<img width="1529" height="103" alt="image" src="https://github.com/user-attachments/assets/e956bb7e-eb24-4a99-8a63-e6165282b939" />
<img width="1255" height="239" alt="image" src="https://github.com/user-attachments/assets/9293e874-cdd8-48b5-9d19-d81723387dd4" />

#### [Bài 6] Phân trang danh sách (Pagination)
```bash
curl.exe -X GET "http://localhost:8080/assets?page=1&limit=5"
```
<img width="1600" height="130" alt="image" src="https://github.com/user-attachments/assets/acfc191b-572d-4ddd-be39-31c782561090" />

#### [Bài 7] Tìm kiếm gần đúng (Search)
```bash
curl.exe -X GET "http://localhost:8080/assets/search?q=firewall"
```
<img width="1607" height="106" alt="image" src="https://github.com/user-attachments/assets/038183b5-74ec-494d-9214-3a329ba97322" />

<a id="phong-chong-khac"></a>
### 🛡️ Các cách phòng chống SQL Injection KHÁC (Đã phân tích)

```text
🛡️ Cách 1: Áp dụng Nguyên tắc đặc quyền tối thiểu (PoLP)
Thay vì để API dùng tài khoản siêu quản trị postgres, tạo một user riêng chỉ có quyền SELECT, INSERT, UPDATE, DELETE trên đúng bảng assets. 

🛡️ Cách 2: Sử dụng ORM hoặc Query Builder
Sử dụng các thư viện ORM (như GORM) tự động hóa hoàn toàn việc "làm sạch" dữ liệu.

🛡️ Cách 3: Lớp giáp hạ tầng - Tường lửa ứng dụng web (WAF)
WAF chặn lại ngay từ ban đầu các cụm từ độc hại (như UNION SELECT).
```
