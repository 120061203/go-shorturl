# 部署指南

## Vercel 部署

### 1. 準備工作

1. **安裝 Vercel CLI**
   ```bash
   npm i -g vercel
   ```

2. **登入 Vercel**
   ```bash
   vercel login
   ```

### 2. 環境變數設置

在 Vercel 項目設置中添加以下環境變數：

- `DATABASE_URL`: Supabase PostgreSQL 連接字符串
- `BASE_URL`: 你的 Vercel 域名 (例如: https://your-app.vercel.app)

### 3. 部署步驟

1. **初始化部署**
   ```bash
   vercel
   ```

2. **設置環境變數**
   ```bash
   vercel env add DATABASE_URL
   vercel env add BASE_URL
   ```

3. **部署到生產環境**
   ```bash
   vercel --prod
   ```

### 4. 項目結構

```
/
├── api/                    # Vercel Serverless Functions
│   ├── shorten.go         # POST /api/shorten
│   ├── redirect.go        # GET /:short_code
│   └── stats.go           # GET /api/stats/:short_code
├── frontend/              # Vue.js 前端
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── internal/              # Go 後端邏輯
│   ├── db/
│   ├── handlers/
│   └── models/
├── cmd/                   # 本地開發服務器
├── vercel.json           # Vercel 配置
└── go.mod
```

### 5. API 端點

- `POST /api/shorten` - 創建短網址
- `GET /:short_code` - 重定向到原始網址
- `GET /api/stats/:short_code` - 獲取統計數據

### 6. 路由配置

- API 請求路由到對應的 Serverless Functions
- 靜態文件路由到前端構建結果
- 短網址重定向路由到 redirect.go

### 7. 本地開發

```bash
# 啟動後端
go run cmd/server/main.go

# 啟動前端
cd frontend && npm run dev
```

### 8. 故障排除

1. **構建失敗**
   - 檢查 Go 版本兼容性
   - 確認所有依賴都已安裝

2. **環境變數問題**
   - 確認 DATABASE_URL 格式正確
   - 檢查 BASE_URL 設置

3. **路由問題**
   - 檢查 vercel.json 配置
   - 確認 API 端點正確

### 9. 監控和日誌

- 在 Vercel Dashboard 查看函數執行日誌
- 監控 API 響應時間和錯誤率
- 設置告警通知

### 10. 性能優化

- 使用 Supabase 連接池
- 實施緩存策略
- 優化數據庫查詢
