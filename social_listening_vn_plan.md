# Social Listening Vietnam - Market & Implementation Plan

## 1. Tong quan

File SVG `social_listening_vn_flow.svg` mo ta mot flow san pham Social Listening danh cho SME Viet Nam.

Muc tieu la xay dung cong cu theo doi nhac den thuong hieu, keyword va sentiment voi chi phi thap, de dung, uu tien tieng Viet va phu hop voi doanh nghiep nho.

## 2. Phan tich thi truong

### 2.1. Doi thu enterprise

- Buzzmetrics / Kompa dang phuc vu chu yeu nhom tap doan lon.
- Chi phi tham chieu trong SVG: khoang `$500+ / bao cao`.
- Mo hinh nay kho tiep can voi shop nho, SME, local brand hoac agency nho.

### 2.2. Khoang trong SME Viet Nam

- SME can biet ai dang noi ve brand, san pham, doi thu va keyword nganh hang.
- Nhu cau chinh:
  - Theo doi mention tren social va news.
  - Phat hien khung hoang hoac comment tieu cuc som.
  - Xem sentiment tong quan.
  - Co bao cao don gian, de hieu, khong can analyst chuyen sau.
- Khoang trong: chua co tool gia re, Viet hoa tot, setup nhanh cho SME.

### 2.3. Co hoi san pham

- Dinh vi: Social Listening nhe, gia re, danh cho SME Viet Nam.
- Goi muc tieu: `300-500k/thang`.
- Gia tri chinh:
  - Setup keyword nhanh.
  - Dashboard don gian.
  - Alert Telegram khi co mention moi hoac sentiment xau.
  - Phan tich tieng Viet bang AI.

## 3. Kien truc tong quan

### 3.1. Data sources

Nguon du lieu ban dau:

- Facebook Scraper qua RapidAPI.
- Instagram API qua RapidAPI.
- YouTube API qua RapidAPI.
- Google News qua RSS/API.

Thu tu uu tien MVP:

1. Facebook.
2. YouTube.
3. Google News.
4. Instagram.

Ly do: Facebook va YouTube co gia tri cao voi SME Viet Nam, de demo ro use case, va du de ship MVP som.

### 3.2. Backend aggregator

Backend dung Go + Gin.

Thanh phan chinh:

- API server quan ly brands, keywords, mentions va dashboard data.
- Worker pool de crawl nhieu keyword/source song song.
- Rate limit theo tung provider de tranh vuot quota.
- Redis cache de giam goi API lap lai.
- Dedup pipeline de loai bo mention trung.
- Cron job crawl moi 30 phut.

### 3.3. Sentiment Analysis

Phan tich sentiment dung Gemini API, uu tien tieng Viet.

Output toi thieu:

- `positive`
- `neutral`
- `negative`

Co the luu them:

- `confidence`
- `reason`
- `detected_language`

Keyword rules dung de bo tro AI trong cac truong hop ro rang, vi du:

- Tu khoa khieu nai: "lua dao", "te", "khong giao hang", "bao hanh kem".
- Tu khoa tich cuc: "tot", "hai long", "se ung ho", "dang tien".

### 3.4. Database

Dung PostgreSQL.

Bang toi thieu:

- `brands`: thong tin brand/khach hang.
- `keywords`: keyword can theo doi.
- `mentions`: mention thu thap duoc tu cac source.

Truong quan trong cua `mentions`:

- `brand_id`
- `keyword_id`
- `source`
- `external_id`
- `url`
- `author_name`
- `content`
- `published_at`
- `engagement_count`
- `sentiment`
- `sentiment_confidence`
- `sentiment_reason`
- `created_at`

Dedup theo uu tien:

1. `source + external_id` neu provider co ID on dinh.
2. Hash cua `source + url`.
3. Hash cua `source + content + published_at` neu thieu ID va URL.

### 3.5. Dashboard

Dashboard co the lam bang React hoac HTML don gian trong MVP.

Man hinh can co:

- Tong so mentions theo ngay.
- Ty le sentiment positive/neutral/negative.
- Danh sach mention moi nhat.
- Filter theo brand, keyword, source va sentiment.
- Link mo nguon mention goc.

### 3.6. Alert Telegram

Alert Telegram dung cho case can phan ung nhanh.

Trigger ban dau:

- Co mention moi.
- Co mention sentiment negative.
- Co keyword canh bao xuat hien trong content.

Noi dung alert:

- Brand.
- Keyword matched.
- Source.
- Sentiment.
- Doan content ngan.
- URL nguon.

## 4. Plan trien khai theo tung y

### 4.1. Market validation

- Chon mot nganh doc de test truoc, vi du spa, nha khoa, F&B local hoac shop online.
- Tim 3-5 chu SME/agency nho de hoi ve nhu cau theo doi brand mention.
- Xac nhan muc gia chap nhan duoc: `199k`, `499k`, `1.5tr`.
- Lay 1 khach hang pilot truoc khi mo rong source va tinh nang.

Ket qua can co:

- Mot nhom nganh uu tien.
- Mot danh sach keyword mau.
- Mot khach hang hoac use case pilot.

### 4.2. Data ingestion

- Tao abstraction cho moi data source theo interface chung.
- Ket noi Facebook va YouTube truoc.
- Chuan hoa payload ve schema mention chung.
- Luu raw payload neu can debug provider.
- Them retry co gioi han cho loi tam thoi.
- Log loi provider nhung khong lam dung toan bo crawl job.

Ket qua can co:

- Crawl duoc mention tu 2 source dau tien.
- Moi mention co source, URL, content, timestamp va keyword matched.

### 4.3. Backend aggregator

- Tao Go + Gin service.
- Tao endpoint quan ly brand va keyword.
- Tao cron crawl moi 30 phut.
- Dung worker pool de xu ly nhieu keyword.
- Them rate limit rieng cho tung provider.
- Cache ket qua ngan han bang Redis neu provider bi goi lap.

Ket qua can co:

- Co the them brand/keyword.
- Job crawl tu dong chay theo lich.
- Du lieu vao DB on dinh, khong trung lap nhieu.

### 4.4. Dedup va database

- Tao schema PostgreSQL cho `brands`, `keywords`, `mentions`.
- Them unique constraint cho dedup key.
- Them index cho query theo `brand_id`, `keyword_id`, `source`, `published_at`, `sentiment`.
- Luu mention moi truoc, sau do enqueue sentiment analysis.

Ket qua can co:

- Mention trung khong bi insert lap.
- Dashboard query nhanh theo thoi gian va filter co ban.

### 4.5. Sentiment Analysis

- Tao service goi Gemini API.
- Gui noi dung mention kem ngu canh keyword/brand.
- Parse output ve enum sentiment.
- Neu AI loi hoac timeout, gan sentiment `neutral` hoac `unknown` tuy schema chon.
- Ket hop keyword rules cho cac tu khoa rui ro cao.

Ket qua can co:

- Mention co sentiment.
- Negative mention co the dung de trigger alert.

### 4.6. Dashboard

- MVP co the dung HTML server-rendered hoac React nhe.
- Hien thi chart mentions theo ngay.
- Hien thi sentiment breakdown.
- Hien thi list mentions moi nhat.
- Co filter source, keyword, sentiment.
- Co link den mention goc.

Ket qua can co:

- Nguoi dung xem duoc tinh hinh brand trong 1 man hinh.
- Du de demo va ban pilot.

### 4.7. Telegram alert

- Tao Telegram bot.
- Luu chat ID theo brand/customer.
- Gui alert khi co mention moi hoac negative.
- Them co che chong spam, vi du gom alert trong 5-10 phut neu mention qua nhieu.

Ket qua can co:

- Khach hang nhan duoc alert dung luc.
- Alert co link nguon de xu ly nhanh.

### 4.8. Pricing

Goi goi y tu SVG:

- Starter: `199k/thang`, 2 keywords.
- Pro: `499k/thang`, 10 keywords.
- Agency: `1.5tr/thang`, unlimited.

Goi MVP nen ban truoc:

- Tap trung Pro `499k/thang`.
- Starter dung de lay user nho.
- Agency chi ban khi co nhu cau quan ly nhieu brand.

## 5. Roadmap 2 tuan

### Tuan 1 - Core

- Setup Go + Gin project.
- Setup PostgreSQL.
- Tao schema `brands`, `keywords`, `mentions`.
- Ket noi Facebook + YouTube truoc.
- Tao cron job crawl moi 30 phut.
- Luu data vao DB.
- Them dedup mention.
- Tao API doc noi bo cho endpoint chinh.

### Tuan 2 - Ship

- Tich hop Gemini API cho sentiment tieng Viet.
- Tao dashboard HTML hoac React don gian.
- Hien thi chart mentions, sentiment va list mention.
- Tich hop Telegram alert.
- Deploy len VPS.
- Demo voi 1 khach pilot.

## 6. Acceptance criteria

- Them duoc brand va keyword can theo doi.
- Crawl duoc mention tu it nhat Facebook va YouTube.
- Mention duoc dedup truoc khi luu.
- Mention co sentiment sau khi xu ly AI.
- Dashboard hien thi duoc mentions theo thoi gian va sentiment.
- Telegram gui alert khi co mention moi hoac negative.
- Co pricing va demo flow ro de ban thu cho 1 khach.

## 7. Rui ro va cach giam

### API source khong on dinh

- Giam rui ro bang cache, retry, logging va provider abstraction.
- MVP khong phu thuoc vao qua nhieu source ngay tu dau.

### Chi phi AI tang

- Chi goi Gemini cho mention moi.
- Cache sentiment theo content hash.
- Dung keyword rules de xu ly nhanh cac case ro rang.

### Spam alert

- Gom alert theo batch ngan.
- Chi alert ngay voi negative/high-risk keywords.

### Sai sentiment tieng Viet

- Luu `reason` de debug.
- Thu thap sample that tu khach pilot.
- Them rule rieng theo nganh sau khi co data.

## 8. Test plan

- Test crawl voi keyword mau va data source mau.
- Test dedup bang cach crawl lai cung mot keyword.
- Test sentiment voi cac cau tieng Viet positive, neutral, negative.
- Test dashboard filter theo brand, keyword, source, sentiment.
- Test Telegram alert voi mention moi va mention negative.
- Test job crawl khong dung toan bo pipeline khi mot provider loi.

## 9. Assumptions

- MVP uu tien toc do ship va ban pilot hon la day du tat ca source.
- Facebook va YouTube la 2 source dau tien.
- PostgreSQL la database chinh.
- Redis dung cho cache va co the bo qua trong ban cuc ky som neu can don gian hoa.
- Dashboard React hoac HTML deu chap nhan duoc, mien la demo duoc gia tri cot loi.
- Chua thuc hien migration, deploy, hay tao code trong buoc lap plan nay.
