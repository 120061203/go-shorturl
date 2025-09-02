#!/bin/bash

# 測試短網址 API 的腳本

BASE_URL="http://localhost:8080"

echo "🚀 測試短網址服務 API"
echo "========================"

# 測試健康檢查
echo "1. 測試健康檢查..."
curl -s "$BASE_URL/health" | jq .
echo ""

# 測試建立短網址
echo "2. 測試建立短網址..."
SHORTEN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/shorten" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com", "custom_code": "test123"}')

echo "$SHORTEN_RESPONSE" | jq .

# 提取短碼
SHORT_CODE=$(echo "$SHORTEN_RESPONSE" | jq -r '.short_code')
echo "短碼: $SHORT_CODE"
echo ""

# 測試重定向 (只檢查狀態碼)
echo "3. 測試重定向..."
REDIRECT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/$SHORT_CODE")
echo "重定向狀態碼: $REDIRECT_STATUS"
echo ""

# 測試統計
echo "4. 測試統計..."
curl -s "$BASE_URL/api/stats/$SHORT_CODE" | jq .
echo ""

echo "✅ 測試完成！"
