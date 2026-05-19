# Seed demo brand + keywords for pilot (run after stack is up)
$Api = if ($env:API_BASE_URL) { $env:API_BASE_URL } else { "http://localhost:8080/api" }

Write-Host "Seeding pilot brand..."
$body = @{
  name = "Demo Spa Pilot"
  keywords = @("tri mun", "cham soc da", "spa quan 1")
  telegramChatId = ""
} | ConvertTo-Json

Invoke-RestMethod -Method Post -Uri "$Api/brands" -ContentType "application/json" -Body $body

Write-Host "Triggering crawl..."
Invoke-RestMethod -Method Post -Uri "$Api/ingest/trigger"

Write-Host "Done. Check dashboard at http://localhost:5173"
