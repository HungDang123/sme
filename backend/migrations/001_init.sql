CREATE TABLE IF NOT EXISTS brands (
  id               TEXT PRIMARY KEY,
  name             TEXT NOT NULL,
  telegram_chat_id TEXT,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS keywords (
  id         TEXT PRIMARY KEY,
  brand_id   TEXT NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
  value      TEXT NOT NULL,
  is_active  BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (brand_id, value)
);

CREATE TABLE IF NOT EXISTS mentions (
  id                   TEXT PRIMARY KEY,
  brand_id             TEXT NOT NULL REFERENCES brands(id),
  keyword_id           TEXT REFERENCES keywords(id),
  keyword              TEXT NOT NULL,
  source               TEXT NOT NULL,
  external_id          TEXT,
  url                  TEXT,
  author_name          TEXT,
  content              TEXT NOT NULL,
  published_at         TIMESTAMPTZ NOT NULL,
  engagement_count     INT DEFAULT 0,
  sentiment            TEXT NOT NULL DEFAULT 'unknown',
  sentiment_confidence DOUBLE PRECISION,
  sentiment_reason     TEXT,
  raw_payload          JSONB,
  dedup_key            TEXT NOT NULL,
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (dedup_key)
);

CREATE INDEX IF NOT EXISTS idx_mentions_brand_published ON mentions(brand_id, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_mentions_sentiment ON mentions(sentiment);
CREATE INDEX IF NOT EXISTS idx_mentions_source ON mentions(source);

CREATE TABLE IF NOT EXISTS alerts (
  id         TEXT PRIMARY KEY,
  brand_id   TEXT NOT NULL,
  keyword    TEXT NOT NULL,
  source     TEXT NOT NULL,
  sentiment  TEXT NOT NULL,
  content    TEXT NOT NULL,
  url        TEXT,
  alert_type TEXT NOT NULL DEFAULT 'new',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_alerts_brand_created ON alerts(brand_id, created_at DESC);
