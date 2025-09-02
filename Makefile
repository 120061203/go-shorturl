.PHONY: help build run test clean docker-up docker-down db-migrate

# 預設目標
help:
	@echo "可用的命令:"
	@echo "  build        - 建構應用程式"
	@echo "  run          - 執行應用程式"
	@echo "  test         - 執行測試"
	@echo "  clean        - 清理建構檔案"
	@echo "  docker-up    - 啟動 Docker 容器"
	@echo "  docker-down  - 停止 Docker 容器"
	@echo "  db-migrate   - 執行資料庫遷移"

# 建構應用程式
build:
	go build -o bin/server cmd/server/main.go

# 執行應用程式
run:
	go run cmd/server/main.go

# 執行測試
test:
	go test ./...

# 清理建構檔案
clean:
	rm -rf bin/

# 啟動 Docker 容器
docker-up:
	docker-compose up -d

# 停止 Docker 容器
docker-down:
	docker-compose down

# 執行資料庫遷移
db-migrate:
	@echo "請手動執行資料庫遷移:"
	@echo "psql \$$DATABASE_URL -f db/schema.sql"

# 安裝依賴
deps:
	go mod tidy
	go mod download

# 格式化程式碼
fmt:
	go fmt ./...

# 檢查程式碼
vet:
	go vet ./...

# 執行所有檢查
check: fmt vet test
