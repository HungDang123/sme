import { useCallback, useEffect, useMemo, useState } from 'react';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';
const SOURCES = ['facebook', 'youtube', 'manual'];
const SENTIMENTS = ['positive', 'neutral', 'negative'];

export default function App() {
  const [brands, setBrands] = useState([]);
  const [mentions, setMentions] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [brandForm, setBrandForm] = useState({
    name: 'Demo Spa',
    keywords: 'tri mun, cham soc da',
    telegramChatId: ''
  });
  const [filters, setFilters] = useState({
    brandId: '',
    keyword: '',
    source: '',
    sentiment: '',
    from: '',
    to: ''
  });
  const [showManual, setShowManual] = useState(false);
  const [mentionForm, setMentionForm] = useState({
    brandId: '',
    brandName: 'Demo Spa',
    keyword: 'tri mun',
    source: 'manual',
    content: 'Dich vu tot, minh rat hai long'
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const stats = useMemo(() => {
    return mentions.reduce(
      (acc, mention) => {
        acc.total += 1;
        acc[mention.sentiment] = (acc[mention.sentiment] || 0) + 1;
        return acc;
      },
      { total: 0, positive: 0, neutral: 0, negative: 0 }
    );
  }, [mentions]);

  const dailyChart = useMemo(() => buildDailyChart(mentions, 7), [mentions]);

  const mentionQuery = useMemo(() => {
    const params = new URLSearchParams();
    if (filters.brandId) params.set('brandId', filters.brandId);
    if (filters.keyword) params.set('keyword', filters.keyword);
    if (filters.source) params.set('source', filters.source);
    if (filters.sentiment) params.set('sentiment', filters.sentiment);
    if (filters.from) params.set('from', new Date(filters.from).toISOString());
    if (filters.to) params.set('to', new Date(filters.to + 'T23:59:59').toISOString());
    params.set('limit', '200');
    return `/mentions?${params.toString()}`;
  }, [filters]);

  const loadData = useCallback(async () => {
    setError('');
    setLoading(true);
    try {
      const [brandRes, mentionRes, alertRes] = await Promise.all([
        request('/brands'),
        request(mentionQuery),
        request('/alerts').catch(() => ({ data: [] }))
      ]);
      setBrands(brandRes.data || []);
      setMentions(mentionRes.data || []);
      setAlerts(alertRes.data || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [mentionQuery]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  useEffect(() => {
    const timer = setInterval(loadData, 60000);
    return () => clearInterval(timer);
  }, [loadData]);

  async function createBrand(event) {
    event.preventDefault();
    setError('');
    try {
      const res = await request('/brands', {
        method: 'POST',
        body: JSON.stringify({
          name: brandForm.name,
          keywords: brandForm.keywords.split(',').map((item) => item.trim()).filter(Boolean),
          telegramChatId: brandForm.telegramChatId.trim()
        })
      });
      setMentionForm((current) => ({
        ...current,
        brandId: res.data.id,
        brandName: res.data.name,
        keyword: res.data.keywords?.[0] || current.keyword
      }));
      await loadData();
    } catch (err) {
      setError(err.message);
    }
  }

  async function createMention(event) {
    event.preventDefault();
    setError('');
    try {
      await request('/mentions', { method: 'POST', body: JSON.stringify(mentionForm) });
      await loadData();
    } catch (err) {
      setError(err.message);
    }
  }

  async function triggerCrawl() {
    setError('');
    try {
      await request('/ingest/trigger', { method: 'POST' });
      setTimeout(loadData, 3000);
    } catch (err) {
      setError(err.message);
    }
  }

  const maxDaily = Math.max(...dailyChart.map((d) => d.count), 1);

  return (
    <main className="app-shell">
      <header className="topbar">
        <div>
          <h1>SME Social Listening</h1>
          <p>Theo doi mention, sentiment tieng Viet va alert Telegram cho SME.</p>
        </div>
        <div className="topbar-actions">
          <button type="button" onClick={loadData} disabled={loading}>
            {loading ? 'Dang tai...' : 'Refresh'}
          </button>
          <button type="button" className="secondary" onClick={triggerCrawl}>
            Crawl ngay
          </button>
        </div>
      </header>

      {error && <div className="error">{error}</div>}

      <section className="stats-grid">
        <Stat label="Mentions" value={stats.total} />
        <Stat label="Positive" value={stats.positive} className="positive" />
        <Stat label="Neutral" value={stats.neutral} className="neutral" />
        <Stat label="Negative" value={stats.negative} className="negative" />
      </section>

      <section className="panel chart-panel">
        <h2>Mentions 7 ngay gan nhat</h2>
        <div className="bar-chart">
          {dailyChart.map((day) => (
            <div key={day.label} className="bar-col">
              <div
                className="bar"
                style={{ height: `${Math.max(8, (day.count / maxDaily) * 100)}%` }}
                title={`${day.label}: ${day.count}`}
              />
              <span>{day.count}</span>
              <small>{day.label}</small>
            </div>
          ))}
        </div>
      </section>

      <section className="panel filters-panel">
        <h2>Bo loc</h2>
        <div className="filters-grid">
          <label>
            Brand
            <select value={filters.brandId} onChange={(e) => setFilters({ ...filters, brandId: e.target.value })}>
              <option value="">Tat ca</option>
              {brands.map((b) => (
                <option key={b.id} value={b.id}>{b.name}</option>
              ))}
            </select>
          </label>
          <label>
            Keyword
            <input value={filters.keyword} onChange={(e) => setFilters({ ...filters, keyword: e.target.value })} />
          </label>
          <label>
            Source
            <select value={filters.source} onChange={(e) => setFilters({ ...filters, source: e.target.value })}>
              <option value="">Tat ca</option>
              {SOURCES.map((s) => (
                <option key={s} value={s}>{s}</option>
              ))}
            </select>
          </label>
          <label>
            Sentiment
            <select value={filters.sentiment} onChange={(e) => setFilters({ ...filters, sentiment: e.target.value })}>
              <option value="">Tat ca</option>
              {SENTIMENTS.map((s) => (
                <option key={s} value={s}>{s}</option>
              ))}
            </select>
          </label>
          <label>
            Tu ngay
            <input type="date" value={filters.from} onChange={(e) => setFilters({ ...filters, from: e.target.value })} />
          </label>
          <label>
            Den ngay
            <input type="date" value={filters.to} onChange={(e) => setFilters({ ...filters, to: e.target.value })} />
          </label>
        </div>
      </section>

      <section className="work-grid">
        <form className="panel" onSubmit={createBrand}>
          <h2>Them brand</h2>
          <label>
            Ten brand
            <input value={brandForm.name} onChange={(e) => setBrandForm({ ...brandForm, name: e.target.value })} />
          </label>
          <label>
            Keywords (phay)
            <input value={brandForm.keywords} onChange={(e) => setBrandForm({ ...brandForm, keywords: e.target.value })} />
          </label>
          <label>
            Telegram Chat ID
            <input
              value={brandForm.telegramChatId}
              onChange={(e) => setBrandForm({ ...brandForm, telegramChatId: e.target.value })}
            />
          </label>
          <button type="submit">Luu brand</button>
        </form>

        <section className="panel">
          <h2>Nhap tay (demo)</h2>
          <p className="hint">Du lieu chinh den tu crawl Facebook/YouTube moi 30 phut.</p>
          <button type="button" className="secondary" onClick={() => setShowManual((v) => !v)}>
            {showManual ? 'An form' : 'Hien form nhap tay'}
          </button>
          {showManual && (
            <form className="manual-form" onSubmit={createMention}>
              <label>
                Brand ID
                <input value={mentionForm.brandId} onChange={(e) => setMentionForm({ ...mentionForm, brandId: e.target.value })} />
              </label>
              <label>
                Keyword
                <input value={mentionForm.keyword} onChange={(e) => setMentionForm({ ...mentionForm, keyword: e.target.value })} />
              </label>
              <label>
                Content
                <textarea value={mentionForm.content} onChange={(e) => setMentionForm({ ...mentionForm, content: e.target.value })} rows={3} />
              </label>
              <button type="submit">Tao mention</button>
            </form>
          )}
        </section>
      </section>

      <section className="panel">
        <h2>Brands ({brands.length})</h2>
        <div className="list">
          {brands.map((brand) => (
            <article key={brand.id} className="list-row">
              <strong>{brand.name}</strong>
              <span>{brand.keywords?.join(', ') || 'No keywords'}</span>
              {brand.telegramChatId && <span>Telegram: {brand.telegramChatId}</span>}
              <code>{brand.id}</code>
            </article>
          ))}
          {brands.length === 0 && <p className="empty">Chua co brand.</p>}
        </div>
      </section>

      <section className="panel">
        <h2>Mentions ({mentions.length})</h2>
        <div className="list">
          {mentions.map((mention) => (
            <article key={mention.id} className="list-row mention-row">
              <div>
                <strong>{mention.source} · {mention.keyword}</strong>
                <p>{mention.content}</p>
                <small>{formatDate(mention.publishedAt)}</small>
              </div>
              <div className="mention-meta">
                <span className={`sentiment ${mention.sentiment}`}>{mention.sentiment}</span>
                {mention.url && <a href={mention.url} target="_blank" rel="noreferrer">Mo nguon</a>}
              </div>
            </article>
          ))}
          {mentions.length === 0 && <p className="empty">Chua co mention. Bam Crawl ngay de thu.</p>}
        </div>
      </section>

      <section className="panel">
        <h2>Alerts gan day</h2>
        <div className="list">
          {alerts.slice(0, 20).map((alert) => (
            <article key={alert.id} className="list-row">
              <strong>{alert.alertType || 'new'} · {alert.keyword}</strong>
              <span>{alert.sentiment} · {alert.source}</span>
              <p>{alert.content}</p>
            </article>
          ))}
          {alerts.length === 0 && <p className="empty">Chua co alert.</p>}
        </div>
      </section>
    </main>
  );
}

function Stat({ label, value, className = '' }) {
  return (
    <div className={`stat-card ${className}`}>
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function buildDailyChart(mentions, days) {
  const result = [];
  const now = new Date();
  for (let i = days - 1; i >= 0; i--) {
    const d = new Date(now);
    d.setDate(d.getDate() - i);
    const key = d.toISOString().slice(0, 10);
    result.push({
      label: key.slice(5),
      count: mentions.filter((m) => (m.publishedAt || m.createdAt || '').slice(0, 10) === key).length
    });
  }
  return result;
}

function formatDate(value) {
  if (!value) return '';
  return new Date(value).toLocaleString('vi-VN');
}

async function request(path, options = {}) {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
    ...options
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Request failed with ${res.status}`);
  }
  if (res.status === 204) return {};
  return res.json();
}
