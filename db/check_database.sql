-- 数据库检查脚本
-- 在 Zeabur PostgreSQL 的 "執行 SQL 指令" 中运行

-- 1. 检查所有表是否存在
SELECT 
    table_name,
    table_type
FROM information_schema.tables 
WHERE table_schema = 'public'
ORDER BY table_name;

-- 2. 检查 urls 表是否存在及其结构
SELECT 
    column_name,
    data_type,
    character_maximum_length,
    is_nullable
FROM information_schema.columns
WHERE table_schema = 'public' 
AND table_name = 'urls'
ORDER BY ordinal_position;

-- 3. 检查 clicks 表是否存在及其结构
SELECT 
    column_name,
    data_type,
    character_maximum_length,
    is_nullable
FROM information_schema.columns
WHERE table_schema = 'public' 
AND table_name = 'clicks'
ORDER BY ordinal_position;

-- 4. 查看 urls 表的数据（如果有）
SELECT 
    id,
    short_code,
    original_url,
    created_at
FROM urls
ORDER BY created_at DESC
LIMIT 10;

-- 5. 查看 clicks 表的数据统计
SELECT 
    COUNT(*) as total_clicks,
    COUNT(DISTINCT url_id) as unique_urls,
    MIN(clicked_at) as first_click,
    MAX(clicked_at) as last_click
FROM clicks;

-- 6. 查看 clicks 表的详细数据（最近 10 条）
SELECT 
    c.id,
    c.clicked_at,
    c.ip_address,
    c.device_type,
    c.location_country,
    c.location_city,
    u.short_code
FROM clicks c
LEFT JOIN urls u ON c.url_id = u.id
ORDER BY c.clicked_at DESC
LIMIT 10;

-- 7. 检查索引
SELECT 
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
AND tablename IN ('urls', 'clicks')
ORDER BY tablename, indexname;



