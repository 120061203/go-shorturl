-- 完整的数据库初始化脚本
-- 用于 Zeabur 或其他 PostgreSQL 数据库

-- 1. 创建 urls 表
CREATE TABLE IF NOT EXISTS urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    original_url TEXT NOT NULL,
    short_code VARCHAR(16) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- 2. 创建 clicks 表（基础版本）
CREATE TABLE IF NOT EXISTS clicks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url_id UUID REFERENCES urls(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP DEFAULT now(),
    ip_address TEXT,
    user_agent TEXT,
    referrer TEXT
);

-- 3. 添加设备类型和地理位置字段
ALTER TABLE clicks 
ADD COLUMN IF NOT EXISTS device_type VARCHAR(30),
ADD COLUMN IF NOT EXISTS location TEXT;

-- 4. 添加详细地理位置字段
ALTER TABLE clicks 
ADD COLUMN IF NOT EXISTS location_isp VARCHAR(255),
ADD COLUMN IF NOT EXISTS location_hostname VARCHAR(255),
ADD COLUMN IF NOT EXISTS location_country VARCHAR(100),
ADD COLUMN IF NOT EXISTS location_region VARCHAR(100),
ADD COLUMN IF NOT EXISTS location_city VARCHAR(100),
ADD COLUMN IF NOT EXISTS location_zip VARCHAR(20);

-- 5. 创建索引以提升查询效能
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
CREATE INDEX IF NOT EXISTS idx_clicks_url_id ON clicks(url_id);
CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);
CREATE INDEX IF NOT EXISTS idx_clicks_device_type ON clicks(device_type);
CREATE INDEX IF NOT EXISTS idx_clicks_location ON clicks(location);
CREATE INDEX IF NOT EXISTS idx_clicks_location_country ON clicks(location_country);
CREATE INDEX IF NOT EXISTS idx_clicks_location_region ON clicks(location_region);
CREATE INDEX IF NOT EXISTS idx_clicks_location_city ON clicks(location_city);
CREATE INDEX IF NOT EXISTS idx_clicks_location_isp ON clicks(location_isp);

-- 6. 验证表是否创建成功
SELECT 'Tables created successfully!' as status;
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' AND table_name IN ('urls', 'clicks');

