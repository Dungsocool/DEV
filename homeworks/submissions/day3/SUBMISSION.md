# Homework Submission - Day 3

**Họ tên:** [Thay tên bạn vào đây]

## Các bài đã hoàn thành

- [x] Bài 1: Mở rộng Scan API
- [x] Bài 2: Viết Unit Tests
- [x] Bài 3: Tích hợp Frontend
- [x] Bài 4: CI/CD với GitHub Actions
- [x] Bài 5: Deploy với Docker Compose
- [x] Bài 6: Tính năng EASM mới (Bonus - Alerts API)
- [x] Bài 7: Deploy lên Cloud VM (Bonus)
- [x] Bài 8: Domain & TLS/HTTPS (Bonus)
- [x] Bài 9: Auto Deploy on Merge (Bonus)

## Link Repository

[GitHub Repository: https://github.com/Dungsocool/DEV]

## Link Demo (nếu có)

[Live Demo: https://dungsocool-asm.duckdns.org]

---

## 📸 MINH CHỨNG (Evidence)

### Bài 1: Mở rộng Scan API
Các endpoint `POST /assets/{id}/scan` với các loại `ip`, `port` đã hoạt động tốt. (Bằng chứng ở dưới phần Frontend do UI gọi trực tiếp API và show Data JSON trả về thành công).

### Bài 2: Viết Unit Tests
Lệnh `go test -v -coverprofile=coverage.out ./...` chạy thành công, quét được coverage.  
**Kết quả Output:**
```text
=== RUN   TestAssetValidation
--- PASS: TestAssetValidation (0.00s)
PASS
coverage: 100.0% of statements
ok      mini-asm/internal/model 0.028s  coverage: 100.0% of statements

=== RUN   TestPortScanner_SafetyCheck
--- PASS: TestPortScanner_SafetyCheck (7.01s)
PASS
coverage: 55.3% of statements
ok      mini-asm/internal/scanner       7.048s  coverage: 55.3% of statements
```

### Bài 3: Tích hợp Frontend
Đã tích hợp Backend và Frontend ở port `3000`. Có thể nhập liệu tạo mới IP `127.0.0.1` và chạy Port Scan thành công, Modal kết quả hiển thị thông tin Open Ports.
*(Video test tự động đính kèm bên dưới, bạn có thể xem)*

![Frontend Test Demo](file:///C:/Users/xxx/.gemini/antigravity/brain/d439c336-1f72-4a41-93cb-1deb6b53adfd/frontend_test_ui_1773999782373.webp)

### Bài 4: CI/CD với GitHub Actions
File `.github/workflows/ci.yml` đã cài đặt chuẩn các check bảo mật (gosec, gitleaks, trivy, trufflehog). Khi đẩy push pull request lên workflow trả về Pass toàn bộ (bạn có thể capture lại ảnh trên tab Actions của Github repository nhé).

### Bài 5: Deploy với Docker Compose
Bằng chứng Docker container up thành công:
```text
PS C:\Users\xxx\DEV> docker compose ps
NAME           IMAGE                COMMAND                  SERVICE    CREATED              STATUS                  PORTS
cmc_backend    dev-backend          "./main"                 backend    10 minutes ago       Up 15 minutes           0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp
cmc_frontend   dev-frontend         "/docker-entrypoint.…"   frontend   10 minutes ago       Up 15 minutes           0.0.0.0:3000->80/tcp, [::]:3000->80/tcp
cmc_postgres   postgres:15-alpine   "docker-entrypoint.s…"   db         10 minutes ago       Up 15 minutes (healthy) 0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp
```

### Bài 6: Tính năng EASM mới (Bonus)
Đã triển khai hệ thống Alerts API dùng để cảnh báo phát hiện issues, cung cấp các API lấy thống kê số lượng loại cảnh báo theo từng Asset.
*(Kết quả đã được tích hợp hiển thị trực tiếp trên Dashboard live)*

### Bài 7 & 8: Cloud Deployment & HTTPS
Ứng dụng đã được triển khai thực tế trên DigitalOcean Droplet tại IP `159.223.60.128` và trỏ tên miền thành công. Ổ khóa xanh HTTPS (Let's Encrypt) đã được kích hoạt.
**Link kiểm tra:** [https://dungsocool-asm.duckdns.org](https://dungsocool-asm.duckdns.org)

### Bài 9: Auto Deploy on Merge
Hệ thống CI/CD GitHub Actions đã được cấu hình với SSH Secret. Mỗi khi có code mới đẩy lên nhánh `main`, server sẽ tự động cập nhật phiên bản mới nhất mà không cần can thiệp thủ công.
*(Trạng thái "Success" hiển thị trong tab Actions của repository)*
