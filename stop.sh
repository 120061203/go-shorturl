#!/bin/bash

echo "🛑 停止短網址服務系統"
echo "====================="

# 停止前端服務
echo "1. 停止前端服務..."
pkill -f "npm run dev" 2>/dev/null
pkill -f "vite" 2>/dev/null

# 停止後端服務
echo "2. 停止後端服務..."
pkill -f "go run" 2>/dev/null
docker stop $(docker ps -q --filter "ancestor=golang:1.21") 2>/dev/null

# 停止資料庫
echo "3. 停止資料庫..."
docker-compose down

echo ""
echo "✅ 所有服務已停止"
echo ""
echo "💡 重新啟動："
echo "   - 完整啟動: ./start.sh"
echo "   - 僅資料庫: docker-compose up -d"
echo "   - 僅前端: cd frontend && npm run dev"
