-- 添加詳細地理位置字段到clicks表
ALTER TABLE clicks 
ADD COLUMN IF NOT EXISTS location_isp VARCHAR(255),
ADD COLUMN IF NOT EXISTS location_hostname VARCHAR(255),
ADD COLUMN IF NOT EXISTS location_country VARCHAR(100),
ADD COLUMN IF NOT EXISTS location_region VARCHAR(100),
ADD COLUMN IF NOT EXISTS location_city VARCHAR(100),
ADD COLUMN IF NOT EXISTS location_zip VARCHAR(20);

-- 為新字段建立索引以提升查詢效能
CREATE INDEX IF NOT EXISTS idx_clicks_location_country ON clicks(location_country);
CREATE INDEX IF NOT EXISTS idx_clicks_location_region ON clicks(location_region);
CREATE INDEX IF NOT EXISTS idx_clicks_location_city ON clicks(location_city);

