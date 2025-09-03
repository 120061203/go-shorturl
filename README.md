# 🚀 xsong 短網址服務

一個簡潔、快速、個人化的短網址工具，專為個人使用而設計。

## ✨ 功能特色

- **簡潔設計**：極簡主義的用戶界面
- **快速縮短**：一鍵生成短網址
- **QR Code**：自動生成並可下載 QR Code
- **自訂短碼**：支援自訂短網址代碼
- **統計分析**：查看點擊統計和裝置分析
- **個人風格**：專屬於 xsong 的設計風格

## 🛠️ 技術架構

- **後端**：Go + Fiber + PostgreSQL
- **前端**：Vue 3 + TypeScript + Tailwind CSS
- **資料庫**：PostgreSQL (Docker)
- **部署**：Docker + Vercel (可選)

## 🚀 快速啟動

### 方法一：一鍵啟動（推薦）

```bash
# 啟動所有服務
./start.sh

# 停止所有服務
./stop.sh

# 測試系統狀態
./test-system.sh
```

### 方法二：手動啟動

```bash
# 1. 啟動資料庫
docker-compose up -d

# 2. 啟動後端服務
docker run --rm --network go-shorturl_default -p 8080:8080 \
  -v $(pwd):/app -w /app \
  -e DATABASE_URL=postgres://devuser:devpass@shorturl-postgres:5432/shortener \
  golang:1.21 go run cmd/server/main.go

# 3. 啟動前端服務（新終端）
cd frontend && npm run dev
```

## 📱 訪問地址

- **前端界面**：http://localhost:5175
- **後端 API**：http://localhost:8080
- **健康檢查**：http://localhost:8080/health

## 🔧 API 端點

### 創建短網址
```bash
POST /api/shorten
Content-Type: application/json

{
  "url": "https://www.example.com",
  "custom_code": "my-custom-code"  // 可選
}
```

### 查詢統計
```bash
GET /api/stats/{short_code}
```

### 重定向
```bash
GET /{short_code}
```

## 🎨 設計特色

### 個人風格
- **深色主題**：現代化的深色設計
- **漸變色彩**：紫色到粉色的漸變效果
- **毛玻璃效果**：backdrop-blur 的現代設計
- **簡潔佈局**：去除不必要的商業元素

### 功能簡化
- 專注於核心功能：縮短網址和 QR Code
- 移除商業元素：公司介紹、服務條款等
- 個人化體驗：專屬於 xsong 的設計

## 🐳 Docker 管理

```bash
# 查看容器狀態
docker ps

# 查看資料庫日誌
docker logs shorturl-postgres

# 重啟資料庫
docker-compose restart

# 清理所有容器
docker-compose down -v
```

## 🔍 故障排除

### 常見問題

1. **資料庫連接失敗**
   ```bash
   # 檢查容器狀態
   docker ps
   
   # 重新啟動資料庫
   docker-compose down && docker-compose up -d
   ```

2. **前端無法訪問**
   ```bash
   # 檢查端口
   lsof -i :5175
   
   # 重新啟動前端
   cd frontend && npm run dev
   ```

3. **後端 API 錯誤**
   ```bash
   # 檢查後端日誌
   docker logs $(docker ps -q --filter "ancestor=golang:1.21")
   
   # 重新啟動後端
   ./stop.sh && ./start.sh
   ```

## 📝 開發指南

### 專案結構
```
go-shorturl/
├── cmd/server/          # 後端入口
├── internal/            # 內部包
│   ├── db/             # 資料庫連接
│   ├── handlers/       # API 處理器
│   └── models/         # 資料模型
├── frontend/           # Vue 前端
├── db/                 # 資料庫腳本
├── scripts/            # 工具腳本
└── docker-compose.yml  # Docker 配置
```

### 開發命令
```bash
# 後端開發
go run cmd/server/main.go

# 前端開發
cd frontend && npm run dev

# 資料庫遷移
supabase db push
```

## 📄 授權

MIT License

## 🤝 關於

這是一個專屬於 xsong 的個人短網址服務，簡潔、快速、實用。

---

**xsong** - 個人短網址服務 🚀
# Updated at Wed Sep  3 16:33:18 CST 2025
