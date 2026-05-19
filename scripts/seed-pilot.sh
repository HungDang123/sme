#!/bin/sh
# Seed demo brand + keywords for pilot (run after stack is up)
API="${API_BASE_URL:-http://localhost:8080/api}"

echo "Seeding pilot brand..."
curl -s -X POST "$API/brands" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Demo Spa Pilot",
    "keywords": ["tri mun", "cham soc da", "spa quan 1"],
    "telegramChatId": ""
  }' | head -c 500
echo ""
echo "Triggering crawl..."
curl -s -X POST "$API/ingest/trigger"
echo ""
echo "Done. Check dashboard at http://localhost:5173"
