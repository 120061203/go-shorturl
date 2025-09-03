#!/bin/bash

echo "🚀 啟動短網址服務系統"
echo "======================"

# 檢查 Docker 是否運行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未運行，請先啟動 Docker Desktop"
    exit 1
fi

# 啟動 PostgreSQL
echo "1. 啟動 PostgreSQL 資料庫..."
docker-compose up -d postgres

# 等待資料庫啟動
echo "2. 等待資料庫啟動..."
sleep 5

# 檢查資料庫連接
if docker exec shorturl-postgres psql -U devuser -d shortener -c "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ 資料庫連接正常"
else
    echo "❌ 資料庫連接失敗"
    exit 1
fi

# 啟動後端服務
echo "3. 啟動後端服務..."
docker run --rm --network go-shorturl_default -p 8080:8080 -v $(pwd):/app -w /app -e DATABASE_URL=postgres://devuser:devpass@shorturl-postgres:5432/shortener golang:1.21 go run cmd/server/main.go &
BACKEND_PID=$!

# 等待後端啟動
echo "4. 等待後端服務啟動..."
sleep 10

# 檢查後端服務
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ 後端服務正常"
else
    echo "❌ 後端服務啟動失敗"
    kill $BACKEND_PID 2>/dev/null
    exit 1
fi

# 啟動前端服務
echo "5. 啟動前端服務..."
cd frontend && npm run dev &
FRONTEND_PID=$!

# 等待前端啟動
echo "6. 等待前端服務啟動..."
sleep 10

# 檢查前端服務
if curl -s http://localhost:5175 > /dev/null; then
    echo "✅ 前端服務正常"
else
    echo "❌ 前端服務啟動失敗"
    kill $FRONTEND_PID 2>/dev/null
    kill $BACKEND_PID 2>/dev/null
    exit 1
fi

echo ""
echo "🎉 系統啟動完成！"
echo ""
echo "📱 前端地址: http://localhost:5175"
echo "🔧 後端地址: http://localhost:8080"
echo "📊 健康檢查: http://localhost:8080/health"
echo ""
echo "💡 管理命令："
echo "   - 停止服務: ./stop.sh"
echo "   - 查看狀態: ./test-system.sh"
echo "   - 查看日誌: docker logs shorturl-postgres"
echo ""
echo "按 Ctrl+C 停止所有服務"

# 等待用戶中斷
trap 'echo ""; echo "🛑 正在停止服務..."; kill $BACKEND_PID 2>/dev/null; kill $FRONTEND_PID 2>/dev/null; docker-compose down; echo "✅ 服務已停止"; exit 0' INT

wait
