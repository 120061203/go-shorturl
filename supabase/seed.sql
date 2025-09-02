-- 插入一些測試資料
INSERT INTO urls (original_url, short_code) VALUES
('https://www.google.com', 'google'),
('https://www.github.com', 'github'),
('https://www.stackoverflow.com', 'stack');

-- 插入一些測試點擊記錄
INSERT INTO clicks (url_id, ip_address, user_agent, referrer) 
SELECT 
    u.id,
    '127.0.0.1',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    'https://example.com'
FROM urls u 
WHERE u.short_code = 'google';

INSERT INTO clicks (url_id, ip_address, user_agent, referrer) 
SELECT 
    u.id,
    '192.168.1.1',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    'https://github.com'
FROM urls u 
WHERE u.short_code = 'github';
