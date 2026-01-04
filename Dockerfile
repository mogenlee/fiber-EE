# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git ca-certificates tzdata

# 安装 swag（用于生成 swagger 文档）
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 生成 swagger 文档
RUN swag init -g cmd/server/main.go -o internal/docs --parseDependency

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装证书和时区
RUN apk --no-cache add ca-certificates tzdata

# 从构建阶段复制
COPY --from=builder /app/server .
COPY --from=builder /app/config/config.yaml ./config/
COPY --from=builder /app/internal/locales ./internal/locales/
COPY --from=builder /app/internal/docs ./internal/docs/

# 创建非 root 用户
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

ENTRYPOINT ["./server"]
