-- 添加設備類型和地理位置字段到clicks表
ALTER TABLE clicks 
ADD COLUMN IF NOT EXISTS device_type VARCHAR(30),
ADD COLUMN IF NOT EXISTS location TEXT;

-- 如果字段已存在但長度不夠，更新長度
ALTER TABLE clicks 
ALTER COLUMN device_type TYPE VARCHAR(30);

-- 為新字段建立索引以提升查詢效能
CREATE INDEX IF NOT EXISTS idx_clicks_device_type ON clicks(device_type);
CREATE INDEX IF NOT EXISTS idx_clicks_location ON clicks(location);

