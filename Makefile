.PHONY: build run clean test migrate-up migrate-down migrate-create gen help

# 默认目标
all: build

# 变量
APP_NAME := gin-admin-base
BUILD_DIR := build
MAIN_FILE := main.go
MIGRATE_CMD := cmd/migrate/main.go
GEN_CMD := cmd/gen/main.go

# Go 相关
GO := go
GOFLAGS := -ldflags="-s -w"
GOTEST := $(GO) test

# 颜色输出
GREEN := \033[32m
YELLOW := \033[33m
CYAN := \033[36m
RESET := \033[0m

## build: 编译项目
build:
	@echo "$(CYAN)▶ 编译项目...$(RESET)"
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "$(GREEN)✓ 编译完成: $(BUILD_DIR)/$(APP_NAME)$(RESET)"

## build-all: 编译所有包（检查编译错误）
build-all:
	@echo "$(CYAN)▶ 编译所有包...$(RESET)"
	$(GO) build ./...
	@echo "$(GREEN)✓ 所有包编译通过$(RESET)"

## run: 启动服务
run:
	@echo "$(CYAN)▶ 启动服务...$(RESET)"
	$(GO) run $(MAIN_FILE)

## clean: 清理构建产物
clean:
	@echo "$(YELLOW)▶ 清理构建产物...$(RESET)"
	rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✓ 清理完成$(RESET)"

## test: 运行测试
test:
	@echo "$(CYAN)▶ 运行测试...$(RESET)"
	$(GOTEST) -v ./...
	@echo "$(GREEN)✓ 测试完成$(RESET)"

## test-coverage: 运行测试并生成覆盖率报告
test-coverage:
	@echo "$(CYAN)▶ 运行测试并生成覆盖率报告...$(RESET)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ 覆盖率报告已生成: coverage.html$(RESET)"

## migrate-up: 执行数据库迁移
migrate-up:
	@echo "$(CYAN)▶ 执行数据库迁移...$(RESET)"
	$(GO) run $(MIGRATE_CMD) up
	@echo "$(GREEN)✓ 迁移完成$(RESET)"

## migrate-down: 回滚数据库迁移（默认回滚 1 步，可指定步数）
migrate-down:
	@echo "$(YELLOW)▶ 回滚数据库迁移...$(RESET)"
	$(GO) run $(MIGRATE_CMD) down $(steps)
	@echo "$(GREEN)✓ 回滚完成$(RESET)"

## migrate-history: 查看迁移历史（默认显示 10 条）
migrate-history:
	@echo "$(CYAN)▶ 查看迁移历史...$(RESET)"
	$(GO) run $(MIGRATE_CMD) history $(limit)

## migrate-create: 创建新的迁移文件（用法: make migrate-create name=<迁移名称>）
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "$(YELLOW)用法: make migrate-create name=<迁移名称>$(RESET)"; \
		exit 1; \
	fi
	@echo "$(CYAN)▶ 创建迁移: $(name)$(RESET)"
	$(GO) run $(MIGRATE_CMD) create $(name)
	@echo "$(GREEN)✓ 迁移文件已创建$(RESET)"

## gen: 代码生成
gen:
	@echo "$(CYAN)▶ 执行代码生成...$(RESET)"
	$(GO) run $(GEN_CMD)
	@echo "$(GREEN)✓ 代码生成完成$(RESET)"

## tidy: 整理 Go 模块依赖
tidy:
	@echo "$(CYAN)▶ 整理依赖...$(RESET)"
	$(GO) mod tidy
	@echo "$(GREEN)✓ 依赖整理完成$(RESET)"

## fmt: 格式化代码
fmt:
	@echo "$(CYAN)▶ 格式化代码...$(RESET)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ 代码格式化完成$(RESET)"

## vet: 代码静态检查
vet:
	@echo "$(CYAN)▶ 代码静态检查...$(RESET)"
	$(GO) vet ./...
	@echo "$(GREEN)✓ 静态检查通过$(RESET)"

## lint: 运行 golangci-lint（需安装）
lint:
	@echo "$(CYAN)▶ 运行 golangci-lint...$(RESET)"
	golangci-lint run ./...
	@echo "$(GREEN)✓ lint 检查通过$(RESET)"

## help: 显示帮助信息
help:
	@echo "$(CYAN)用法: make <target>$(RESET)"
	@echo ""
	@echo "$(CYAN)可用目标:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?## / { printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(CYAN)迁移命令示例:$(RESET)"
	@echo "  make migrate-up                          # 执行所有未应用的迁移"
	@echo "  make migrate-down steps=1                # 回滚 1 步"
	@echo "  make migrate-history limit=10            # 查看最近 10 条迁移记录"
	@echo "  make migrate-create name=add_user_table  # 创建新的迁移文件"
