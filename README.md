# 🚀 CMC Intern API Project - Asset Management

Đây là dự án xây dựng hệ thống RESTful API dùng để quản lý các tài sản IT (Domain, IP, Service), được phát triển bằng Go (Golang) và cơ sở dữ liệu PostgreSQL. Dự án tuân thủ nghiêm ngặt mô hình **Clean Architecture** để tối ưu hóa việc bảo trì và phân tách rõ ràng các nghiệp vụ.
## 📑 Mục lục (Table of Contents)
* [🛠️ Công nghệ & Thư viện sử dụng](#cong-nghe)
* [🛡️ Tiêu điểm Bảo mật: Ngăn chặn hoàn toàn SQL Injection](#bao-mat)
* [✅ Các tính năng đã hoàn thành](#tinh-nang)
* [📂 Cấu trúc thư mục (Project Structure)](#cau-truc)
* [📡 Danh sách API Endpoints](#api)
* [💻 Hướng dẫn Cài đặt & Khởi chạy (How to Run)](#cai-dat)
* [🧪 Hướng Dẫn Test Từng Bài (API Testing Guide)](#test)
* [🛡️ Các cách phòng chống SQL Injection KHÁC](#phong-chong-khac)

---
<a id="cong-nghe"></a>
## 🛠️ Công nghệ & Thư viện sử dụng
* **Ngôn ngữ:** Go (Golang)
* **Cơ sở dữ liệu:** PostgreSQL
* **Bộ định tuyến (Router):** `gorilla/mux` (Kết hợp thư viện chuẩn `net/http`)
* **Hạ tầng:** Docker & Docker Compose
* **Kiến trúc:** Clean Architecture (Handler -> Service -> Storage)

<a id="bao-mat"></a>
## 🛡️ Tiêu điểm Bảo mật: Ngăn chặn hoàn toàn SQL Injection (SQLi)

Trong dự án này, rủi ro SQL Injection đã được triệt tiêu 100% tại tầng Storage bằng cách áp dụng triệt để kỹ thuật **Parameterized Queries** (Truy vấn tham số hóa) thông qua driver `database/sql` của Go.

Nguyên lý hoạt động (Tại sao nó an toàn?)
Thay vì sử dụng cách nối chuỗi (String Concatenation) nguy hiểm để tạo câu lệnh SQL, hệ thống tách bạch hoàn toàn giữa **"Cấu trúc lệnh" (Code)** và **"Dữ liệu" (Data)**. 

Khi một truy vấn được gửi đi, quá trình diễn ra qua 2 bước bảo mật khép kín:


**Bước 1 (Prepare)** : Database nhận cấu trúc câu lệnh SQL với các "chỗ trống" ảo (Placeholders như `$1, $2`). Nó tiến hành biên dịch (compile) cấu trúc lệnh này trước.


**Bước 2 (Execute)**: Dữ liệu người dùng nhập vào mới được gửi xuống để "lấp" vào các chỗ trống đó. Lúc này, Database xử lý chúng hoàn toàn dưới dạng **Văn bản thuần túy (Literal Values)**, tuyệt đối không coi đó là một phần của câu lệnh thực thi.

#### Minh họa thực tế 
📍 Bài 3 (Batch Delete)



<img width="786" height="523" alt="image" src="https://github.com/user-attachments/assets/20bbf63d-7889-4736-99ff-152d3204a86b" />


Sử dụng động Parameterized Queries với toán tử IN để bảo mật tính năng xóa hàng loạt, không cộng chuỗi ID trực tiếp.

📍  Bài 7 (Search by Name)
<img width="1074" height="337" alt="image" src="https://github.com/user-attachments/assets/cea5d9fd-febf-4d38-9c47-79840fde173a" />

Sử dụng tham số $1 kết hợp toán tử ILIKE. Ký tự '%' được nối vào biến ở tầng Go, đảm bảo cấu trúc lệnh SQL không bị phá vỡ bởi input của người dùng.

#### Kết luận: Nhờ cơ chế bọc tham số $1, $2 của PostgreSQL và thư viện Go, mọi ký tự đặc biệt như dấu nháy đơn ('), dấu chấm phẩy (;) hay các lệnh SQL lồng ghép từ Hacker đều bị vô hiệu hóa hoàn toàn trước khi chạm vào dữ liệu thực tế.

<a id="tinh-nang"></a>
## ✅ Các tính năng đã hoàn thành
Dự án đã hoàn thiện toàn bộ các yêu cầu cốt lõi và nâng cao:
1. **Bài 1 - Statistics APIs:** Cung cấp API thống kê tổng tài sản và đếm số lượng tài sản theo các bộ lọc (loại, trạng thái).
2. **Bài 2 - Batch Create:** Thêm mới hàng loạt tài sản (tối đa 100 items/request) an toàn tuyệt đối với cơ chế **Database Transaction** (All-or-nothing).
3. **Bài 3 - Batch Delete:** Xóa hàng loạt tài sản thông minh dựa trên danh sách ids.
4. **Bài 4 - Connection Retry:** Xây dựng thuật toán **Exponential Backoff** giúp Server tự động thử lại kết nối (tối đa 5 lần) nếu Database bị sập hoặc khởi động chậm.
5. **Bài 5 - Health Check:** Cung cấp API giám sát tình trạng hệ thống, trả về trạng thái của Database và các thông số của Connection Pool.
6. **Bài 6 - Pagination & Filtering:** Xử lý phân trang và lọc dữ liệu động bằng SQL để tối ưu hóa hiệu năng truy vấn danh sách tài sản.
7. **Bài 7 - Search by Name:** Hỗ trợ tìm kiếm tài sản theo tên với cơ chế khớp một phần (Partial match / Case-insensitive).
<a id="cau-truc"></a>
## 📂 Cấu trúc thư mục (Project Structure)
Dự án được chia thành các package nhỏ lẻ tuân theo Clean Architecture:

```text
├── cmd/
│   └── server/
│       └── main.go                 # Entry point: Điểm khởi chạy server
├── homeworks/
│   └── submissions/
│       ├── SUBMISSION.md           # File checklist nộp bài
│       └── (Các file ảnh minh chứng)
├── internal/
│   ├── handler/                    # Tầng giao tiếp HTTP
│   │   └── asset_handler.go        # Xử lý Request/Response cho Assets
│   ├── model/                      # Định nghĩa các cấu trúc dữ liệu
│   │   ├── bonus.go
│   │   └── stats.go
│   ├── service/                    # Tầng xử lý Logic nghiệp vụ
│   │   └── asset_service.go
│   └── storage/                    # Tầng giao tiếp Cơ sở dữ liệu
│       ├── storage.go              # Interface của Storage
│       └── postgres/               # Triển khai cụ thể với PostgreSQL
│           ├── asset_storage.go    # Query CRUD cơ bản
│           ├── bonus_storage.go    # Query cho tính năng phân trang, tìm kiếm
│           ├── delete_storage.go   # Query cho tính năng xóa hàng loạt
│           ├── postgres.go         # Cấu hình kết nối và Retry DB
│           └── stats_storage.go    # Query cho tính năng thống kê
├── docker-compose.yml              # File khởi tạo database PostgreSQL
└── README.md

```

<a id="api"></a>
## 📡 Danh sách API Endpoints

| Method | Endpoint | Tính năng |
| :--- | :--- | :--- |
| `GET` | `/health` | Kiểm tra trạng thái của Server và Database |
| `GET` | `/assets/stats` | Thống kê số lượng tài sản theo loại và trạng thái |
| `GET` | `/assets/count` | Đếm tài sản (hỗ trợ Query: `?type=...&status=...`) |
| `POST` | `/assets/batch` | Tạo mới nhiều tài sản cùng lúc |
| `DELETE` | `/assets/batch` | Xóa nhiều tài sản (Query: `?ids=id1,id2...`) |
| `GET` | `/assets` | Lấy danh sách tài sản (Query: `?page=1&limit=20`) |

| `GET` | `/assets/search` | Tìm kiếm tài sản theo tên (Query: `?q=keyword`) |
<a id="cai-dat"></a>
## 💻 Hướng dẫn Cài đặt & Khởi chạy (How to Run)
1. Khởi chạy Database:
Đảm bảo máy bạn đã cài Docker. Mở terminal và chạy lệnh sau để dựng database:





docker-compose up -d



2. Khởi chạy API Server:




Tiếp tục chạy lệnh sau để khởi động ứng dụng Go:

go run cmd/server/main.go




Ghi chú: Nhờ tính năng Connection Retry (Bài 4), Server sẽ tự động chờ và kết nối lại nếu PostgreSQL đang trong quá trình khởi động.

Server sẽ lắng nghe các request tại địa chỉ: http://localhost:8080

<a id="test"></a>
## 🧪 Hướng Dẫn Test Từng Bài (API Testing Guide)

*(Lưu ý: Đảm bảo Server đang chạy ở `http://localhost:8080` trước khi chạy lệnh. Nếu bạn dùng Windows PowerShell, hãy gõ `curl.exe` thay cho `curl` hoặc sử dụng Git Bash).*

### [Bài 1] Thống kê & Đếm tài sản (Statistics & Count )

1. Lấy báo cáo tổng quan (Đếm tổng, nhóm theo Type và Status)

curl.exe -X GET http://localhost:8080/assets/stats
<img width="697" height="67" alt="image" src="https://github.com/user-attachments/assets/50f12b73-ebc5-4092-9122-a6e68b4cf119" />

2. Đếm tài sản kết hợp bộ lọc động (Ví dụ: Đếm số IP đang active)



curl.exe -X GET "http://localhost:8080/assets/count?type=ip&status=active"
<img width="894" height="71" alt="image" src="https://github.com/user-attachments/assets/0bbee3b1-f2a5-4b55-b928-2b613bd3ae07" />


### [Bài 2] Thêm hàng loạt tài sản (Batch Create)

Sử dụng Transaction để đảm bảo tính toàn vẹn dữ liệu (Tối đa 100 tài sản/lần).

curl.exe -X POST http://localhost:8080/assets/batch \
-H "Content-Type: application/json" \
-d '{
  "assets": [
    {"name": "Firewall-Core", "type": "ip"},
    {"name": "cmc.com.vn", "type": "domain"},
    {"name": "DB-Server-01", "type": "service"}
  ]
}'



<img width="1578" height="96" alt="image" src="https://github.com/user-attachments/assets/c5e77937-ee30-4bb3-a455-bc212485a99a" />

### [Bài 3] Xóa hàng loạt (Batch Delete)

Xóa nhiều tài sản cùng lúc bằng toán tử IN. 



(Vui lòng thay thế chuỗi ID bên dưới bằng các ID thực tế sinh ra từ Bài 2)


curl.exe -X DELETE "http://localhost:8080/assets/batch?ids=id-1,id-2,id-3"
<img width="1069" height="56" alt="image" src="https://github.com/user-attachments/assets/d874ad13-2cc2-46b4-b765-3a69e54f903e" />

### [Bài 4 & 5] Thuật toán Retry & Giám sát sức khỏe (Health Check)

Kiểm tra trạng thái của Server và Database (Bài 5)

curl.exe -X GET http://localhost:8080/health
<img width="799" height="81" alt="image" src="https://github.com/user-attachments/assets/783af0ba-404c-4da2-b732-c45142297ef2" />



Test bài 4 (Retry): Bạn hãy thử tắt container db trong Docker đi, sau đó khởi chạy lại server bằng go run cmd/server/main.go. Bạn sẽ thấy log hệ thống kích hoạt Exponential Backoff, tự động lùi thời gian chờ và thử kết nối lại tối đa 5 lần thay vì sập (panic) ngay lập tức

docker-compose stop db
<img width="1529" height="103" alt="image" src="https://github.com/user-attachments/assets/e956bb7e-eb24-4a99-8a63-e6165282b939" />
<img width="1255" height="239" alt="image" src="https://github.com/user-attachments/assets/9293e874-cdd8-48b5-9d19-d81723387dd4" />

### [Bài 6] Phân trang danh sách (Pagination)

Lấy danh sách tài sản ở trang 1, mỗi trang lấy tối đa 5 bản ghi


curl.exe -X GET "http://localhost:8080/assets?page=1&limit=5"
<img width="1600" height="130" alt="image" src="https://github.com/user-attachments/assets/acfc191b-572d-4ddd-be39-31c782561090" />

### [Bài 7] Tìm kiếm gần đúng (Search)

Tìm kiếm "gần đúng" (ILIKE) không phân biệt hoa thường. 


Ví dụ: Tìm các tài sản có chứa chữ "firewall"



curl.exe -X GET "http://localhost:8080/assets/search?q=firewall"

<img width="1607" height="106" alt="image" src="https://github.com/user-attachments/assets/038183b5-74ec-494d-9214-3a329ba97322" />

<a id="phong-chong-khac"></a>
### 🛡️  Các cách phòng chống SQL Injection KHÁC (Có thể triển khai thêm)

```text
🛡️ Cách 1: Áp dụng Nguyên tắc đặc quyền tối thiểu (Principle of Least Privilege - PoLP)
Triển khai: Thay vì để API Go kết nối vào Database bằng tài khoản siêu quản trị postgres (có quyền xóa toàn bộ DB), ta sẽ tạo một user riêng (ví dụ: api_user).

Tác dụng: Cấp cho api_user này chỉ có quyền SELECT, INSERT, UPDATE, DELETE trên đúng bảng assets. Nếu hacker có tìm ra được một lỗ hổng SQLi đi chăng nữa và gửi lệnh DROP TABLE assets, PostgreSQL sẽ chặn lại ngay lập tức và báo lỗi: "Permission denied".

🛡️ Cách 2: Sử dụng ORM hoặc Query Builder
Triển khai: Trong các dự án Go thực tế lớn hơn, thay vì tự viết SQL thuần (database/sql) như hiện tại, người ta thường dùng các thư viện ORM (như GORM) hoặc Query Builder (như Squirrel).

Tác dụng: Các thư viện này tự động hóa hoàn toàn việc "làm sạch" (sanitize) dữ liệu và bọc tham số. Bạn chỉ cần viết code dạng db.Where("name = ?", userInput), thư viện sẽ tự động lo phần chống SQL Injection bên dưới.

🛡️ Cách 3: Lớp giáp hạ tầng - Tường lửa ứng dụng web (WAF)
Triển khai: Triển khai một hệ thống WAF (như Cloudflare, AWS WAF, hoặc ModSecurity) đứng chắn phía trước API Server (Cổng 8080) của bạn.

Tác dụng: WAF có các tập luật (ruleset) nhận diện các dấu hiệu độc hại. Nếu một request gửi lên có chứa các từ khóa nhạy cảm ghép cùng nhau như UNION SELECT, 1=1, hay DROP TABLE, WAF sẽ tự động khóa IP của hacker lại và trả về lỗi 403 Forbidden trước khi request đó kịp chạm tới code Go của bạn.
```







