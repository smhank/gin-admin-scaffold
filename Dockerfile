# ============ 构建阶段 ============
FROM golang:1.26-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server main.go

# ============ 运行阶段 ============
FROM alpine:3.19

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制编译产物
COPY --from=builder /app/server .

# 复制配置文件
COPY --from=builder /app/internal/infras/config/config.yaml ./internal/infras/config/config.yaml

# 创建日志目录
RUN mkdir -p runtime/logs

# 暴露端口
EXPOSE 9501

# 启动服务
CMD ["./server"]
