package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gin-admin-base/internal/domain/model"
	"gin-admin-base/internal/infras/config"
	"gin-admin-base/internal/infras/global"
	"gin-admin-base/internal/infras/logger"
	migrationlib "gin-admin-base/internal/infras/migration"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置和日志
	config.InitConfig()
	global.Logger = logger.InitLogger(config.GetLoggerConfig())
	defer global.Logger.Sync()

	// 连接数据库
	mysqlConfig := config.GetMySQLConfig()
	db, err := gorm.Open(mysql.Open(mysqlConfig.DSN()), &gorm.Config{})
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}

	// 自动建表（确保 migrations 表存在）
	if err := db.AutoMigrate(&model.Migration{}); err != nil {
		fmt.Printf("创建 migrations 表失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ migrations 表已就绪")

	// 创建迁移管理器
	manager := migrationlib.NewManager(db)

	// 注册所有迁移
	registerMigrations(manager)

	// 解析命令
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "up":
		if err := manager.Up(); err != nil {
			fmt.Printf("迁移失败: %v\n", err)
			os.Exit(1)
		}
	case "down":
		steps := 1
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &steps)
		}
		if err := manager.Down(steps); err != nil {
			fmt.Printf("回滚失败: %v\n", err)
			os.Exit(1)
		}
	case "history":
		limit := 10
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &limit)
		}
		records, err := manager.History(limit)
		if err != nil {
			fmt.Printf("查询历史失败: %v\n", err)
			os.Exit(1)
		}
		if len(records) == 0 {
			fmt.Println("没有迁移记录。")
			return
		}
		fmt.Println("\n迁移历史:")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-3s %-50s %-8s %-20s\n", "ID", "迁移名称", "批次", "应用时间")
		fmt.Println(strings.Repeat("-", 80))
		for _, r := range records {
			fmt.Printf("%-3d %-50s %-8d %-20s\n", r.ID, r.Name, r.Batch, r.AppliedAt)
		}
		fmt.Println(strings.Repeat("-", 80))
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("用法: go run cmd/migrate/main.go create <迁移名称>")
			os.Exit(1)
		}
		name := os.Args[2]
		dir := filepath.Join("internal", "infras", "migration")
		filePath, err := migrationlib.NewMigration(name, dir)
		if err != nil {
			fmt.Printf("创建迁移失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("已创建迁移文件: %s\n", filePath)

		// 自动注册迁移到 registerMigrations 函数
		if err := autoRegisterMigration(filePath); err != nil {
			fmt.Printf("警告: 自动注册迁移失败: %v\n", err)
			fmt.Println("请手动将新迁移注册到 registerMigrations 函数中")
		} else {
			fmt.Println("✓ 已自动注册迁移到 registerMigrations")
		}
	default:
		printUsage()
	}
}

// registerMigrations 注册所有迁移
func registerMigrations(manager *migrationlib.Manager) {
	manager.RegisterAll([]migrationlib.Migration{
		migrationlib.NewM20250506000000InitTables(),               // 建表
		migrationlib.NewM20250506000001InitSeedData(),             // 初始化种子数据
		migrationlib.NewM20250506000002AddOperationLogMenu(),      // 添加操作历史菜单
		migrationlib.NewM20260507101614CreateTableTestMigration(), // create table test migration
	})
}

// autoRegisterMigration 自动将新创建的迁移注册到 registerMigrations 函数中
func autoRegisterMigration(filePath string) error {
	// 从文件路径中提取文件名（不含扩展名）
	fileName := filepath.Base(filePath)
	fileName = strings.TrimSuffix(fileName, ".go")

	// 生成构造函数名：将文件名中的下划线替换为空格，首字母大写，再移除空格
	// 例如: m20260507_095637_create_table_test_migration -> M20260507095637CreateTableTestMigration
	funcName := strings.ReplaceAll(fileName, "_", " ")
	funcName = strings.Title(funcName)
	funcName = strings.ReplaceAll(funcName, " ", "")

	// 读取 main.go 文件内容
	mainFile := "cmd/migrate/main.go"
	content, err := os.ReadFile(mainFile)
	if err != nil {
		return fmt.Errorf("读取 %s 失败: %w", mainFile, err)
	}

	// 检查是否已经注册
	registerLine := fmt.Sprintf("migrationlib.New%s()", funcName)
	if strings.Contains(string(content), registerLine) {
		fmt.Printf("  迁移 %s 已注册，跳过\n", funcName)
		return nil
	}

	// 在 registerMigrations 函数的最后一个注册项后面追加新的注册代码
	// 查找 "func registerMigrations" 到第一个 "})" 之间的内容
	oldContent := string(content)

	// 先找到 registerMigrations 函数
	funcStart := strings.Index(oldContent, "func registerMigrations")
	if funcStart == -1 {
		return fmt.Errorf("未找到 registerMigrations 函数")
	}

	// 从函数开始位置查找第一个 "	})"（即函数内的结束标记）
	searchTarget := "\t})"
	lastIndex := strings.Index(oldContent[funcStart:], searchTarget)
	if lastIndex == -1 {
		return fmt.Errorf("未找到迁移注册代码结束标记")
	}
	insertPos := funcStart + lastIndex

	// 构造要插入的注册代码
	comment := fileName
	// 从文件名中提取描述部分（去掉时间戳前缀）
	if idx := strings.Index(fileName, "_"); idx != -1 {
		if idx2 := strings.Index(fileName[idx+1:], "_"); idx2 != -1 {
			comment = fileName[idx+idx2+2:]
		}
	}
	comment = strings.ReplaceAll(comment, "_", " ")
	newLine := fmt.Sprintf("\t\tmigrationlib.New%s(), // %s\n", funcName, comment)

	// 在 "	})" 之前插入新代码
	newContent := oldContent[:insertPos] + newLine + oldContent[insertPos:]

	// 写回文件
	if err := os.WriteFile(mainFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入 %s 失败: %w", mainFile, err)
	}

	return nil
}

func printUsage() {
	fmt.Println("用法:")
	fmt.Println("  go run cmd/migrate/main.go up                  # 执行所有未应用的迁移")
	fmt.Println("  go run cmd/migrate/main.go down [步骤数]        # 回滚迁移")
	fmt.Println("  go run cmd/migrate/main.go history [条数]       # 查看迁移历史")
	fmt.Println("  go run cmd/migrate/main.go create <迁移名称>    # 创建新的迁移文件")
}
