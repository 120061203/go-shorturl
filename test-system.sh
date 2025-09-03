#!/bin/bash

echo "🚀 測試短網址服務系統"
echo "========================"

# 測試後端健康檢查
echo "1. 測試後端健康檢查..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ 後端服務正常"
else
    echo "❌ 後端服務異常"
    exit 1
fi

# 測試創建短網址
echo "2. 測試創建短網址..."
SHORT_URL_RESPONSE=$(curl -s -X POST http://localhost:8080/api/shorten \
    -H "Content-Type: application/json" \
    -d '{"url": "https://www.example.com"}')

if echo "$SHORT_URL_RESPONSE" | grep -q "short_url"; then
    echo "✅ 短網址創建成功"
    SHORT_CODE=$(echo "$SHORT_URL_RESPONSE" | grep -o '"short_code":"[^"]*"' | cut -d'"' -f4)
    echo "   短碼: $SHORT_CODE"
else
    echo "❌ 短網址創建失敗"
    echo "$SHORT_URL_RESPONSE"
    exit 1
fi

# 測試統計查詢
echo "3. 測試統計查詢..."
STATS_RESPONSE=$(curl -s http://localhost:8080/api/stats/$SHORT_CODE)

if echo "$STATS_RESPONSE" | grep -q "total_clicks"; then
    echo "✅ 統計查詢成功"
else
    echo "❌ 統計查詢失敗"
    echo "$STATS_RESPONSE"
    exit 1
fi

# 測試前端服務
echo "4. 測試前端服務..."
if curl -s http://localhost:5176 > /dev/null; then
    echo "✅ 前端服務正常"
else
    echo "❌ 前端服務異常"
    exit 1
fi

echo ""
echo "🎉 所有測試通過！系統運行正常"
echo ""
echo "📱 前端地址: http://localhost:5176"
echo "🔧 後端地址: http://localhost:8080"
echo "📊 健康檢查: http://localhost:8080/health"
echo ""
echo "💡 提示："
echo "   - 在瀏覽器中打開 http://localhost:5176 查看前端"
echo "   - 使用 API 測試工具測試後端功能"
echo "   - 查看 Docker 容器狀態: docker ps"
