#!/bin/bash

# 部署腳本
echo "🚀 開始部署到 Vercel..."

# 檢查是否已安裝 Vercel CLI
if ! command -v vercel &> /dev/null; then
    echo "❌ Vercel CLI 未安裝，正在安裝..."
    npm install -g vercel
fi

# 檢查是否已登入
if ! vercel whoami &> /dev/null; then
    echo "❌ 請先登入 Vercel"
    vercel login
fi

# 構建前端
echo "📦 構建前端..."
cd frontend
npm run build
cd ..

# 部署到 Vercel
echo "🌐 部署到 Vercel..."
vercel --prod

echo "✅ 部署完成！"
echo "📝 請記得在 Vercel Dashboard 設置環境變數："
echo "   - DATABASE_URL: 你的 Supabase 連接字符串"
echo "   - BASE_URL: 你的 Vercel 域名"
