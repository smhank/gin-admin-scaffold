# gin-admin-scaffold

基于 Gin + Vue3 的后台管理系统脚手架。

## 技术栈

### 后端
- **框架**: Gin (Go HTTP 框架)
- **ORM**: GORM (MySQL)
- **架构**: Clean Architecture / DDD 风格分层
- **缓存**: Redis
- **消息队列**: Redis Streams / MySQL 双驱动
- **定时任务**: 内置轻量级调度器
- **数据迁移**: 类 Yii2 迁移机制
- **日志**: Zap

### 前端
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP 客户端**: Axios
- **桌面端**: Electron (可选)

## 项目结构

```
├── main.go                          # 应用入口
├── cmd/
│   ├── gen/                         # 代码生成工具
│   └── migrate/                     # 数据迁移 CLI
├── internal/
│   ├── application/                 # 应用服务层 (用例编排)
│   ├── domain/
│   │   ├── model/                   # 领域模型
│   │   └── repository/              # 仓储接口
│   ├── infras/
│   │   ├── cache/                   # Redis 缓存
│   │   ├── config/                  # 配置管理 (Viper)
│   │   ├── cron/                    # 定时任务调度器
│   │   ├── global/                  # 全局单例
│   │   ├── logger/                  # 日志 (Zap)
│   │   ├── migration/               # 数据迁移
│   │   ├── persistence/             # 持久化 (GORM)
│   │   ├── query/                   # 查询 (GORM gen)
│   │   └── queue/                   # 消息队列
│   └── interfaces/
│       ├── app/                     # 应用初始化
│       ├── handler/                 # HTTP 处理器
│       ├── middleware/              # 中间件
│       ├── response/                # 响应封装
│       └── router/                  # 路由注册
├── web/                             # 前端项目
├── docs/                            # 文档
├── runtime/                         # 运行时文件 (日志等)
├── docker-compose.yml               # Docker 编排
├── Dockerfile                       # 后端 Dockerfile
└── Makefile                         # 构建命令
```

## 快速开始

### 1. 配置

复制配置文件并根据环境修改：

```bash
cp internal/infras/config/config.yaml config.yaml
```

### 2. 启动后端

```bash
go run main.go
```

### 3. 启动前端

```bash
cd web
npm install
npm run dev
```

### 4. 执行数据迁移

```bash
go run cmd/migrate/main.go up
```

## 主要功能

- ✅ 用户认证 (JWT)
- ✅ 角色管理 (RBAC)
- ✅ 权限管理
- ✅ 菜单管理
- ✅ 路径管理
- ✅ 操作日志
- ✅ 数据迁移
- ✅ 定时任务
- ✅ 消息队列 (Redis Streams / MySQL)
- ✅ 缓存 (Redis)
- ✅ 请求限流
- ✅ CORS 跨域

## 文档

- [架构图](architecture-diagram.md)
- [定时任务使用指南](docs/cron-guide.md)
- [数据迁移使用指南](docs/migration-guide.md)
- [消息队列使用指南](docs/queue-guide.md)
