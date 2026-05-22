# Homework Submission - Day 3

**Student Name:** Dong Tien Dung

## Completed Exercises

- [x] Exercise 1: EASM Scan Engine API
- [x] Exercise 2: Integrated Unit Tests
- [x] Exercise 3: Integrated UI Frontend
- [x] Exercise 4: CI/CD Security Pipeline with GitHub Actions
- [x] Exercise 5: Infrastructure as Code with Docker Compose
- [x] Exercise 6: EASM Alerts API (Bonus)
- [x] Exercise 7: Deploy to Cloud VM (Bonus)
- [x] Exercise 8: Dynamic Domain & TLS/HTTPS (Bonus)
- [x] Exercise 9: Auto Deploy on Merge (Bonus)

## Repository Link

[GitHub Repository: https://github.com/Dungsocool/easm-platform]

## Demo Link (if applicable)

[Live Demo: https://dungsocool-asm.duckdns.org]

---

## 📷 Evidence (Verification Details)

### Exercise 1: EASM Scan Engine API
The `POST /assets/{id}/scan` endpoints with scan types `ip` and `port` are working perfectly. (Verification can be seen on the frontend UI where the dashboard calls the API directly and renders returned JSON data).

### Exercise 2: Integrated Unit Tests
Successfully ran `go test -v -coverprofile=coverage.out ./...` to calculate coverage.  
**Output Results:**
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

### Exercise 3: Integrated UI Frontend
Integrated backend and frontend at port `3000`. We can add a target IP like `127.0.0.1` and execute a Port Scan successfully. The modal pops up and renders all discovered Open Ports.
*(A screenshot of the frontend interface is shown below)*

![Frontend UI](docs/images/frontend-ui.png)

### Exercise 4: CI/CD Security Pipeline with GitHub Actions
The `.github/workflows/ci.yml` is configured with security checkers (gosec, gitleaks, trivy, trufflehog). Pull Requests triggers tests and all CI jobs successfully pass.

### Exercise 5: Infrastructure as Code with Docker Compose
Successfully launched all services via Docker Compose:
```text
PS C:\Users\xxx\easm-platform> docker compose ps
NAME           IMAGE                COMMAND                  SERVICE    CREATED              STATUS                  PORTS
cmc_backend    dev-backend          "./main"                 backend    10 minutes ago       Up 15 minutes           0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp
cmc_frontend   dev-frontend         "/docker-entrypoint..."   frontend   10 minutes ago       Up 15 minutes           0.0.0.0:3000->80/tcp, [::]:3000->80/tcp
cmc_postgres   postgres:15-alpine   "docker-entrypoint.s..."   db         10 minutes ago       Up 15 minutes (healthy) 0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp
```

### Exercise 6: EASM Alerts API (Bonus)
Deployed the Alerts API system to warn about detected vulnerabilities and issues, providing endpoint statistics by asset.
*(Results are rendered directly on the live Dashboard)*

### Exercise 7 & 8: Cloud Deployment & HTTPS
The application is deployed on a live DigitalOcean Droplet at IP `159.223.60.128` and the domain has been successfully pointed. The Let's Encrypt SSL/TLS secure lock is active.
**Live check link:** [https://dungsocool-asm.duckdns.org](https://dungsocool-asm.duckdns.org)

### Exercise 9: Auto Deploy on Merge
GitHub Actions is configured with SSH Secrets. Every time new code is merged into the `main` branch, the server automatically updates to the latest version.
