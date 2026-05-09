# 数据迁移系统使用指南

## 概述

本项目的数据迁移系统参考 Yii2 框架的迁移机制设计，支持通过 CLI 命令管理数据库结构变更和数据初始化。每个迁移是一个独立的 Go 文件，包含 `Up()`（执行）和 `Down()`（回滚）方法。

## 目录结构

```
internal/infras/migration/          # 迁移核心库
├── migration.go                    # Migration 接口 + BaseMigration 基类
├── manager.go                      # 迁移管理器（up/down/history/create）
├── m20250506_000001_init_seed_data.go        # 迁移示例：初始化种子数据
└── m20250506_000002_add_operation_log_menu.go # 迁移示例：添加操作历史菜单

cmd/migrate/main.go                 # CLI 命令入口
```

## 快速开始

### 查看迁移历史

```bash
go run cmd/migrate/main.go history
```

输出示例：
```
迁移历史:
--------------------------------------------------------------------------------
ID  迁移名称                                              批次     应用时间
1   m20250506_000001_init_seed_data                       1        2026-05-06 15:30:00
2   m20250506_000002_add_operation_log_menu               1        2026-05-06 15:30:01
--------------------------------------------------------------------------------
```

### 执行未应用的迁移

```bash
go run cmd/migrate/main.go up
```

### 回滚最近一批迁移

```bash
# 回滚 1 个迁移
go run cmd/migrate/main.go down

# 回滚 3 个迁移
go run cmd/migrate/main.go down 3
```

### 创建新的迁移文件

```bash
go run cmd/migrate/main.go create add_user_avatar_field
```

会在 `internal/infras/migration/` 目录下生成类似 `m20250506_161500_add_user_avatar_field.go` 的文件。

## 创建迁移

### 1. 使用 CLI 生成模板

```bash
go run cmd/migrate/main.go create add_user_avatar_field
```

### 2. 编辑生成的迁移文件

```go
package migration

import (
    "gin-admin-base/internal/domain/model"
    "gin-admin-base/internal/infras/migration"
    "gorm.io/gorm"
)

// M20250506161500AddUserAvatarField 添加用户头像字段
type M20250506161500AddUserAvatarField struct {
    migration.BaseMigration
}

// NewM20250506161500AddUserAvatarField 创建迁移
func NewM20250506161500AddUserAvatarField() *M20250506161500AddUserAvatarField {
    return &M20250506161500AddUserAvatarField{
        BaseMigration: migration.NewBaseMigration("m20250506_161500_add_user_avatar_field"),
    }
}

// Up 执行迁移
func (m *M20250506161500AddUserAvatarField) Up(db *gorm.DB) error {
    // 添加字段
    return db.Migrator().AddColumn(&model.User{}, "avatar")
}

// Down 回滚迁移
func (m *M20250506161500AddUserAvatarField) Down(db *gorm.DB) error {
    // 删除字段
    return db.Migrator().DropColumn(&model.User{}, "avatar")
}
```

### 3. 注册迁移

在 `cmd/migrate/main.go` 和 `internal/infras/persistence/db.go` 的 `registerMigrations` 函数中注册：

```go
func registerMigrations(manager *migrationlib.Manager) {
    manager.RegisterAll([]migrationlib.Migration{
        migrationlib.NewM20250506000001InitSeedData(),
        migrationlib.NewM20250506000002AddOperationLogMenu(),
        migrationlib.NewM20250506161500AddUserAvatarField(), // 新增
    })
}
```

## 迁移文件命名规范

```
mYYYYMMDD_HHIISS_name.go
```

- `m` - 固定前缀
- `YYYYMMDD` - 创建日期（如 20250506）
- `HHIISS` - 创建时间（如 161500）
- `name` - 简短描述（蛇形命名）

示例：`m20250506_161500_add_user_avatar_field.go`

## 迁移接口说明

```go
// Migration 迁移接口
type Migration interface {
    Name() string                    // 返回迁移名称
    Up(db *gorm.DB) error            // 执行迁移
    Down(db *gorm.DB) error          // 回滚迁移
}

// BaseMigration 基础迁移（嵌入到具体迁移中）
type BaseMigration struct {
    // 已实现 Name() 方法
    // Down() 默认返回 nil（不可回滚）
}
```

## 常用迁移操作

### 建表

```go
func (m *MyMigration) Up(db *gorm.DB) error {
    return db.AutoMigrate(&model.MyTable{})
}

func (m *MyMigration) Down(db *gorm.DB) error {
    return db.Migrator().DropTable(&model.MyTable{})
}
```

### 增删改数据

```go
func (m *MyMigration) Up(db *gorm.DB) error {
    // 插入数据
    return db.Create(&model.Permission{Name: "示例", Code: "example:perm"}).Error
}

func (m *MyMigration) Down(db *gorm.DB) error {
    // 删除数据
    return db.Where("code = ?", "example:perm").Delete(&model.Permission{}).Error
}
```

### 添加字段

```go
func (m *MyMigration) Up(db *gorm.DB) error {
    return db.Migrator().AddColumn(&model.User{}, "avatar")
}

func (m *MyMigration) Down(db *gorm.DB) error {
    return db.Migrator().DropColumn(&model.User{}, "avatar")
}
```

### 创建索引

```go
func (m *MyMigration) Up(db *gorm.DB) error {
    return db.Migrator().CreateIndex(&model.User{}, "idx_email")
}

func (m *MyMigration) Down(db *gorm.DB) error {
    return db.Migrator().DropIndex(&model.User{}, "idx_email")
}
```

## 自动迁移 vs 数据迁移

| 类型 | 机制 | 用途 |
|------|------|------|
| AutoMigrate | GORM 自动同步模型与表结构 | 建表、加字段（增量） |
| 数据迁移 | 自定义迁移文件 | 初始化数据、数据转换、复杂 DDL |

两者配合使用：
1. `AutoMigrate` 保证表结构与模型一致
2. 数据迁移负责种子数据、数据迁移等业务逻辑

## 迁移记录表

迁移执行记录保存在 `migrations` 表中：

```sql
CREATE TABLE `migrations` (
  `id`         bigint(20)   NOT NULL AUTO_INCREMENT,
  `name`       varchar(255) NOT NULL COMMENT '迁移名称',
  `batch`      int(11)      NOT NULL DEFAULT 1 COMMENT '批次号',
  `applied_at` varchar(64)           DEFAULT NULL COMMENT '应用时间',
  `created_at` datetime(3)           DEFAULT NULL,
  `updated_at` datetime(3)           DEFAULT NULL,
  `deleted_at` datetime(3)           DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 最佳实践

1. **迁移不可修改** - 已合并到主分支的迁移文件不要修改，如需变更请创建新的迁移
2. **幂等性** - 迁移应可重复执行（已执行过的迁移会自动跳过）
3. **回滚测试** - 创建迁移时务必实现 `Down()` 方法，确保可回滚
4. **小步提交** - 每个迁移只做一件事，便于回滚和排查问题
5. **先测试后合并** - 在开发环境执行 `up` 和 `down` 验证无误后再提交
