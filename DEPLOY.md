# Deploy SME Social Listening (VPS / Pilot)

## Prerequisites

- Docker + Docker Compose on VPS
- Optional API keys: `RAPIDAPI_KEY`, `GEMINI_API_KEY`, `TELEGRAM_BOT_TOKEN`

## Quick start

```bash
cp .env.example .env
# Edit .env with real secrets
docker compose up --build -d
```

URLs:

| Service | URL |
|---------|-----|
| Dashboard | http://YOUR_HOST:5173 |
| API Gateway | http://YOUR_HOST:8080/api |
| Ingest (manual crawl) | POST http://YOUR_HOST:8080/api/ingest/trigger |

## Seed pilot data

Windows:

```powershell
.\scripts\seed-pilot.ps1
```

Linux/macOS:

```bash
sh scripts/seed-pilot.sh
```

## Pilot test checklist

- [ ] Create brand with keywords via dashboard
- [ ] Click **Crawl ngay** — demo mentions appear (without RapidAPI key)
- [ ] Set `RAPIDAPI_KEY` — real Facebook/YouTube data on next crawl
- [ ] Verify dedup: run crawl twice, mention count stable
- [ ] Set `GEMINI_API_KEY` — neutral content gets AI sentiment
- [ ] Set brand `telegramChatId` + `TELEGRAM_BOT_TOKEN` — alerts delivered
- [ ] Filter mentions by brand, source, sentiment, date range
- [ ] Negative mention triggers immediate Telegram alert

## Environment reference

See [.env.example](.env.example) for all variables.

## Notes

- PostgreSQL data persists in Docker volume `postgres_data`
- Crawl runs every 30 minutes via ingest-service cron
- Without `RAPIDAPI_KEY`, ingest uses demo providers for Facebook/YouTube
