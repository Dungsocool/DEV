$ErrorActionPreference = "Stop"

Write-Host "Tạo Asset..."
$domain_res = try { Invoke-RestMethod -Method Post -Uri "http://localhost:8080/assets" -ContentType "application/json" -Body '{"name":"google.com","type":"domain"}' } catch { $null }
$ip_res = try { Invoke-RestMethod -Method Post -Uri "http://localhost:8080/assets" -ContentType "application/json" -Body '{"name":"127.0.0.1","type":"ip"}' } catch { $null }

$domain_id = if ($domain_res) { $domain_res.data.id } else { $null }
$ip_id = if ($ip_res) { $ip_res.data.id } else { $null }

Write-Host "Tạo Scan Job..."
$scans = @()
if ($ip_id) {
    try { $scans += Invoke-RestMethod -Method Post -Uri "http://localhost:8080/assets/$ip_id/scan" -ContentType "application/json" -Body '{"scan_type":"ip"}' } catch { Write-Host "Error creating IP scan" }
    try { $scans += Invoke-RestMethod -Method Post -Uri "http://localhost:8080/assets/$ip_id/scan" -ContentType "application/json" -Body '{"scan_type":"port"}' } catch { Write-Host "Error creating Port scan" }
}

Write-Host "Chờ 10 giây để scan hoàn thành..."
Start-Sleep -Seconds 10
$results = @{}

foreach ($scan in $scans) {
    if ($scan -and $scan.id) {
        $job_id = $scan.id
        $scan_type = $scan.scan_type
        
        Write-Host "Lấy kết quả cho job $job_id (loại $scan_type)..."
        try {
            $job_res = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/scan-jobs/$job_id/results"
            $results[$scan_type] = $job_res
        } catch {
            Write-Host "Error: $_"
            $results[$scan_type] = "Error retrieving results"
        }
    }
}

$results | ConvertTo-Json -Depth 10 | Out-File "c:\Users\xxx\DEV\scan_results_log.txt"
Write-Host "Xong! Kiểm tra file scan_results_log.txt"
