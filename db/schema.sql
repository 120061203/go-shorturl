-- 短網址表
CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    original_url TEXT NOT NULL,
    short_code VARCHAR(16) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- 點擊紀錄表
CREATE TABLE clicks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url_id UUID REFERENCES urls(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP DEFAULT now(),
    ip_address TEXT,
    user_agent TEXT,
    referrer TEXT
);

-- 建立索引以提升查詢效能
CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_clicks_url_id ON clicks(url_id);
CREATE INDEX idx_clicks_clicked_at ON clicks(clicked_at);
