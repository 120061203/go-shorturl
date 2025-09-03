#!/bin/bash

# 資料庫查詢腳本
# 使用方法: ./scripts/db-query.sh [command]
# 命令選項: tables, urls, clicks, schema, stats

case "$1" in
    "tables")
        echo "=== 查看所有表格 ==="
        docker exec -e PGPASSWORD=devpass shorturl-postgres psql -U devuser -d shortener -c "\dt"
        ;;
    "urls")
        echo "=== 查看 urls 表格內容 ==="
        docker exec -e PGPASSWORD=devpass shorturl-postgres psql -U devuser -d shortener -c "SELECT id, original_url, short_code, created_at FROM urls ORDER BY created_at DESC;"
        ;;
    "clicks")
        echo "=== 查看 clicks 表格內容 ==="
        docker exec -e PGPASSWORD=devpass shorturl-postgres psql -U devuser -d shortener -c "SELECT id, url_id, clicked_at, ip_address FROM clicks ORDER BY clicked_at DESC LIMIT 10;"
        ;;
    "schema")
        echo "=== 查看 urls 表格結構 ==="
        docker exec -e PGPASSWORD=devpass shorturl-postgres psql -U devuser -d shortener -c "\d urls"
        echo ""
        echo "=== 查看 clicks 表格結構 ==="
        docker exec -e PGPASSWORD=devpass shorturl-postgres psql -U devuser -d shortener -c "\d clicks"
        ;;
    "stats")
        echo "=== 統計資訊 ==="
        echo "URLs 總數:"
        docker exec -e PGPASSWORD=devpass shorturl-postgres psql -U devuser -d shortener -c "SELECT COUNT(*) as total_urls FROM urls;"
        echo ""
        echo "Clicks 總數:"
        docker exec -e PGPASSWORD=devpass shorturl-postgres psql -U devuser -d shortener -c "SELECT COUNT(*) as total_clicks FROM clicks;"
        echo ""
        echo "最受歡迎的短網址 (按點擊數排序):"
        docker exec -e PGPASSWORD=devpass shorturl-postgres psql -U devuser -d shortener -c "
        SELECT u.short_code, u.original_url, COUNT(c.id) as click_count
        FROM urls u
        LEFT JOIN clicks c ON u.id = c.url_id
        GROUP BY u.id, u.short_code, u.original_url
        ORDER BY click_count DESC;"
        ;;
    *)
        echo "使用方法: $0 [command]"
        echo "命令選項:"
        echo "  tables  - 查看所有表格"
        echo "  urls    - 查看 urls 表格內容"
        echo "  clicks  - 查看 clicks 表格內容"
        echo "  schema  - 查看表格結構"
        echo "  stats   - 查看統計資訊"
        echo ""
        echo "範例:"
        echo "  $0 urls"
        echo "  $0 stats"
        ;;
esac
