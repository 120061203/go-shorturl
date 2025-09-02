-- 確保 devuser 用戶存在並設定正確的密碼
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'devuser') THEN
        CREATE USER devuser WITH PASSWORD 'devpass' SUPERUSER;
    ELSE
        ALTER USER devuser WITH PASSWORD 'devpass' SUPERUSER;
    END IF;
END
$$;

-- 確保 shortener 資料庫存在
SELECT 'CREATE DATABASE shortener' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'shortener')\gexec

-- 授予所有權限
GRANT ALL PRIVILEGES ON DATABASE shortener TO devuser;
