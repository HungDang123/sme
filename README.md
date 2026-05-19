# SME Social Listening

Source base cho he thong Social Listening danh cho SME Viet Nam.

## Kien truc

- `frontend`: React + Vite.
- `backend/project`: tai lieu kien truc, protocol va convention chung cua backend.
- `backend/gateway-service`: public API cho frontend, proxy HTTP noi bo den cac service.
- `backend/brand-service`: quan ly brand va keyword.
- `backend/mention-service`: tao/list mention, goi `sentiment-service` va `alert-service` qua HTTP noi bo.
- `backend/sentiment-service`: phan tich sentiment ban dau bang keyword rules, co san cho Gemini sau nay.
- `backend/alert-service`: alert + Telegram bot (batch anti-spam).
- `backend/ingest-service`: cron crawl Facebook/YouTube, worker pool, rate limit.
- `backend/migrations`: PostgreSQL schema.

Backend la Monorepo Multi-Service: moi folder service la mot Go project doc lap, co `go.mod` va `Dockerfile` rieng. Cac service khong import code truc tiep cua nhau; giao tiep qua HTTP protocol noi bo.

## Chay local bang Docker Compose

```bash
cp .env.example .env
docker compose up --build
```

Xem [DEPLOY.md](DEPLOY.md) cho pilot checklist va seed script.

URL:

- Frontend: `http://localhost:5173`
- Gateway API: `http://localhost:8080`
- Brand service: `http://localhost:8081`
- Mention service: `http://localhost:8082`
- Sentiment service: `http://localhost:8083`
- Alert service: `http://localhost:8084`

## Chay rieng backend

```bash
cd backend
cd sentiment-service && go mod tidy && go run ./cmd/server
cd ../brand-service && go mod tidy && go run ./cmd/server
cd ../alert-service && go mod tidy && go run ./cmd/server
cd ../mention-service && go mod tidy && go run ./cmd/server
cd ../gateway-service && go mod tidy && go run ./cmd/server
```

Neu chay ngoai Docker, set env neu can:

```bash
BRAND_SERVICE_URL=http://localhost:8081
MENTION_SERVICE_URL=http://localhost:8082
SENTIMENT_SERVICE_URL=http://localhost:8083
ALERT_SERVICE_URL=http://localhost:8084
```

## Chay rieng frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend mac dinh goi gateway qua `http://localhost:8080/api`.
# sme
