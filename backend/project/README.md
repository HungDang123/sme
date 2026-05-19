# Backend Project

Folder nay chua tai lieu chung cho backend monorepo.

## Architecture

Backend dung mo hinh Monorepo Multi-Service:

- Moi service nam trong mot folder rieng duoi `backend/`.
- Moi service la mot Go project doc lap, co `go.mod`, `Dockerfile`, entrypoint `cmd/server/main.go`.
- Service khong import package truc tiep tu service khac.
- Service giao tiep voi nhau qua HTTP noi bo.
- Gateway la service duy nhat frontend goi truc tiep.

## Service Layout

Moi service nen theo layout:

```text
service-name/
  cmd/server/main.go
  internal/config/
  internal/domain/
  internal/service/
  internal/store/
  internal/transport/http/
  go.mod
  Dockerfile
```

Neu service can goi service khac, dat HTTP client trong `internal/upstream/`.

## Services

- `gateway-service`: expose `/api/*`, route request den service noi bo.
- `brand-service`: expose `/brands`, `/keywords/active`, PostgreSQL.
- `mention-service`: expose `/mentions`, dedup PostgreSQL, goi sentiment va alert.
- `sentiment-service`: expose `/sentiment/analyze` (rules + Gemini).
- `alert-service`: expose `/alerts`, Telegram batching.
- `ingest-service`: cron crawl Facebook/YouTube, `POST /ingest/trigger`.
- `migrations/`: PostgreSQL schema (`brands`, `keywords`, `mentions`, `alerts`).

## Internal HTTP Protocol

Tat ca response thanh cong dung envelope:

```json
{
  "data": {}
}
```

Tat ca response loi dung:

```json
{
  "error": "message"
}
```

## Local Service Discovery

Trong Docker Compose, cac service goi nhau bang hostname noi bo:

- `http://brand-service:8081`
- `http://mention-service:8082`
- `http://sentiment-service:8083`
- `http://alert-service:8084`
- `http://ingest-service:8085`
