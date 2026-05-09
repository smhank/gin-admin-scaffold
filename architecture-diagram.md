# 架构图（Mermaid）

> 说明：把本文件在支持 Mermaid 的 Markdown 预览器里打开即可看到图形（例如 GitLab/GitHub/多数 IDE 插件）。

## 1）系统整体组件图（前后端 + 基础设施）

```mermaid
flowchart LR
  %% Clients
  subgraph Clients[客户端]
    Browser[浏览器<br/>Vue3 SPA]
    Electron[Electron 桌面壳<br/>加载 Vite Dev / dist]
  end

  %% Frontend
  subgraph FE[前端（web/）]
    ViteDev[Vite DevServer :5173<br/>Proxy /api -> :9501]
    Dist[dist 静态资源]
    Axios[Axios<br/>baseURL=/api<br/>Authorization=token]
  end

  %% Backend
  subgraph BE[后端（Go / Gin）]
    Gin[Gin HTTP Server<br/>main.go]
    Routes[/api 路由组]
    MW[Middleware<br/>OperationLog]
    Handlers[Handlers<br/>User/Role/Permission/Menu/Path/Auth/...]
    App[Application Services<br/>AuthService]
    Domain[Domain<br/>Models + Repo Interface]
    Infras[Infras<br/>Config/Logger/Persistence/Cache/Queue/Cron/Migration]
  end

  %% Infra services
  subgraph Infra[外部依赖]
    MySQL[(MySQL)]
    Redis[(Redis)]
    Q[(Queue<br/>Redis Streams / MySQL)]
  end

  %% Flows
  Browser -->|HTTP(S)| ViteDev
  Electron --> ViteDev
  Electron -->|生产模式| Dist
  ViteDev --> Axios
  Dist --> Axios
  Axios -->|/api/*| Gin

  Gin --> Routes --> MW --> Handlers
  Handlers --> App --> Domain
  Handlers -->|部分 CRUD 直连| Infras
  Infras --> MySQL
  Infras --> Redis
  Infras --> Q
  MW -->|写操作落库| MySQL
```

## 2）后端分层（Clean/DDD 风格）

```mermaid
flowchart TB
  subgraph Interfaces[interfaces（接口/交付层）]
    H[handler/*<br/>HTTP API]
    M[middleware/*]
  end

  subgraph Application[application（用例层）]
    S[auth_service.go 等<br/>用例编排/业务服务]
  end

  subgraph Domain[domain（领域层）]
    Model[model/*<br/>实体/聚合]
    Repo[repository/*<br/>Repo 接口]
  end

  subgraph Infras[infras（基础设施层）]
    Cfg[config/*<br/>Viper + config.yaml]
    Log[logger/*<br/>zap]
    Persist[persistence/*<br/>GORM + MySQL]
    Cache[cache/*<br/>Redis Client]
    Queue[queue/*<br/>Redis Streams / MySQL]
    Cron[cron/*<br/>Scheduler]
    Mig[migration/*<br/>migrate/up/down/history]
    Global[global/*<br/>全局单例引用]
  end

  H --> S --> Repo
  Repo --> Persist
  Persist --> Mig
  M --> Persist
  S --> Model
  H -->|部分直接使用| Persist
  H --> Cfg
  H --> Log
  Persist --> Global
  Cache --> Global
  Queue --> Global
  Cron --> Global
```

## 3）后端启动时序（main.go）

```mermaid
sequenceDiagram
  autonumber
  participant Main as main.go
  participant Cfg as config.InitConfig
  participant Log as logger.InitLogger
  participant DB as persistence.InitDB
  participant Rds as cache.InitRedis
  participant Q as queue.InitQueue
  participant Cron as cron.Scheduler
  participant Gin as gin.Engine

  Main->>Cfg: 读取 config.yaml / env（Viper）
  Main->>Log: 初始化 zap Logger
  Main->>DB: 连接 MySQL + 自动迁移（Up）
  alt DB 成功
    Main->>Rds: 初始化 Redis Client
    Main->>Q: 初始化队列（redis/mysql）
    Main->>Cron: 注册并启动定时任务
  else DB 失败
    Note over Main: handler 使用 mock 数据模式
  end
  Main->>Gin: 注册 /api 路由 + OperationLogMiddleware
  Main->>Gin: Run 0.0.0.0:9501
```

