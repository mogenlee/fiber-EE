.PHONY: build run dev test clean gen swagger docker-build docker-up docker-down lint

# 应用名称
APP_NAME=fiber-ee
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 构建
build:
	go build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/server

# 运行
run: build
	./bin/$(APP_NAME)

# 开发模式 (需要 air: go install github.com/air-verse/air@latest)
dev:
	air

# 测试
test:
	go test -v -cover ./...

# 测试覆盖率
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 清理
clean:
	rm -rf bin/ coverage.out coverage.html

# 生成 GORM Gen 代码
gen:
	go run ./cmd/gen_orm

# 生成模块代码 (用法: make module name=article)
module:
	go run ./cmd/gen_module -name=$(name)

# 生成 Swagger 文档 (需要 swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	swag init -g cmd/server/main.go -o internal/docs --parseDependency

# 代码检查 (需要 golangci-lint)
lint:
	golangci-lint run ./...

# 格式化
fmt:
	go fmt ./...
	goimports -w .

# 依赖整理
tidy:
	go mod tidy

# Docker 构建
docker-build:
	docker build -t $(APP_NAME):$(VERSION) .

# Docker Compose 启动
docker-up:
	docker-compose up -d

# Docker Compose 停止
docker-down:
	docker-compose down

# 数据库迁移 (示例)
migrate-up:
	@echo "Running migrations..."

migrate-down:
	@echo "Rolling back migrations..."
