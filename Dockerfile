# 使用包含 Go 和 Node.js 的基础镜像
FROM golang:1.21-alpine AS builder

# 安装 Node.js 和 npm
RUN apk add --no-cache nodejs npm

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制前端代码并构建
COPY frontend/ ./frontend/
WORKDIR /app/frontend
RUN npm install && npm run build

# 回到根目录，复制 Go 代码并构建
WORKDIR /app
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY api/ ./api/
RUN go build -o server ./cmd/server

# 使用轻量级镜像运行
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# 从构建阶段复制可执行文件和前端构建结果
COPY --from=builder /app/server .
COPY --from=builder /app/frontend/dist ./frontend/dist

# 暴露端口
EXPOSE 8080

# 运行服务器
CMD ["./server"]

