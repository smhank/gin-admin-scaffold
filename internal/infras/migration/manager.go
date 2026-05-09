package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"gin-admin-base/internal/domain/model"

	"gorm.io/gorm"
)

// Manager 迁移管理器（类似 Yii2 的 migrate 控制器）
type Manager struct {
	db         *gorm.DB
	migrations []Migration
}

// NewManager 创建迁移管理器
func NewManager(db *gorm.DB) *Manager {
	// 自动创建迁移记录表
	db.AutoMigrate(&model.Migration{})

	return &Manager{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// Register 注册迁移（类似 Yii2 的 migrations 数组）
func (m *Manager) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// RegisterAll 批量注册迁移
func (m *Manager) RegisterAll(migrations []Migration) {
	m.migrations = append(m.migrations, migrations...)
}

// Up 执行所有未执行的迁移（类似 Yii2 的 migrate/up）
func (m *Manager) Up() error {
	// 获取当前最大批次号
	var maxBatch int
	m.db.Model(&model.Migration{}).Select("COALESCE(MAX(batch), 0)").Scan(&maxBatch)
	batch := maxBatch + 1

	// 按名称排序执行
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Name() < m.migrations[j].Name()
	})

	applied := 0
	for _, mig := range m.migrations {
		// 检查是否已执行
		var count int64
		m.db.Model(&model.Migration{}).Where("name = ?", mig.Name()).Count(&count)
		if count > 0 {
			fmt.Printf("  - 跳过: %s (已执行)\n", mig.Name())
			continue
		}

		// 执行迁移
		fmt.Printf("*** 正在应用迁移: %s\n", mig.Name())
		if err := mig.Up(m.db); err != nil {
			fmt.Printf("!!! 迁移失败: %s\n", mig.Name())
			fmt.Printf("!!! 错误详情: %v\n", err)
			return fmt.Errorf("迁移 %s 失败: %w", mig.Name(), err)
		}

		// 记录迁移
		record := model.Migration{
			Name:      mig.Name(),
			Batch:     batch,
			AppliedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		if err := m.db.Create(&record).Error; err != nil {
			fmt.Printf("!!! 记录迁移失败: %s\n", mig.Name())
			fmt.Printf("!!! 错误详情: %v\n", err)
			return fmt.Errorf("记录迁移 %s 失败: %w", mig.Name(), err)
		}

		fmt.Printf("     > 已应用: %s (批次: %d)\n", mig.Name(), batch)
		applied++
	}

	if applied == 0 {
		fmt.Println("没有需要应用的迁移。")
	} else {
		fmt.Printf("成功应用 %d 个迁移。\n", applied)
	}

	return nil
}

// Down 回滚最近一批迁移（类似 Yii2 的 migrate/down）
func (m *Manager) Down(steps int) error {
	// 获取最大批次号
	var maxBatch int
	m.db.Model(&model.Migration{}).Select("COALESCE(MAX(batch), 0)").Scan(&maxBatch)
	if maxBatch == 0 {
		fmt.Println("没有可回滚的迁移。")
		return nil
	}

	// 获取指定批次的迁移
	var records []model.Migration
	query := m.db.Model(&model.Migration{}).Where("batch = ?", maxBatch).Order("name desc")
	if steps > 0 {
		query = query.Limit(steps)
	}
	if err := query.Find(&records).Error; err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Println("没有可回滚的迁移。")
		return nil
	}

	// 查找对应的迁移并回滚
	rolledBack := 0
	for _, record := range records {
		for _, mig := range m.migrations {
			if mig.Name() == record.Name {
				fmt.Printf("*** 正在回滚迁移: %s\n", mig.Name())
				if err := mig.Down(m.db); err != nil {
					return fmt.Errorf("回滚迁移 %s 失败: %w", mig.Name(), err)
				}
				// 删除迁移记录
				m.db.Delete(&record)
				fmt.Printf("     > 已回滚: %s\n", mig.Name())
				rolledBack++
				break
			}
		}
	}

	if rolledBack == 0 {
		fmt.Println("没有找到可回滚的迁移实现。")
	} else {
		fmt.Printf("成功回滚 %d 个迁移。\n", rolledBack)
	}

	return nil
}

// History 查看迁移历史（类似 Yii2 的 migrate/history）
func (m *Manager) History(limit int) ([]model.Migration, error) {
	var records []model.Migration
	query := m.db.Model(&model.Migration{}).Order("batch desc, name desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// NewMigration 创建新的迁移文件（类似 Yii2 的 migrate/create）
func NewMigration(name string, dir string) (string, error) {
	// 生成文件名：mYYYYMMDD_HHIISS_name
	// 注意：文件名不能以 _test 结尾，否则 Go 会将其视为测试文件
	now := time.Now()
	safeName := strings.TrimSuffix(name, "_test")
	safeName = strings.TrimSuffix(safeName, "_test")
	fileName := fmt.Sprintf("m%s_%s", now.Format("20060102_150405"), safeName)
	filePath := filepath.Join(dir, fileName+".go")

	// 检查文件是否已存在
	if _, err := os.Stat(filePath); err == nil {
		return "", fmt.Errorf("迁移文件已存在: %s", filePath)
	}

	// 生成迁移模板
	tpl := template.Must(template.New("migration").Parse(`package migration

import (
	"gorm.io/gorm"
)

// {{.ClassName}} {{.Description}}
type {{.ClassName}} struct {
	BaseMigration
}

// New{{.ClassName}} 创建迁移
func New{{.ClassName}}() *{{.ClassName}} {
	return &{{.ClassName}}{
		BaseMigration: NewBaseMigration("{{.FileName}}"),
	}
}

// Up 执行迁移
func (m *{{.ClassName}}) Up(db *gorm.DB) error {
	// TODO: 在这里编写迁移逻辑
	// 示例: db.Create(&model.Permission{Name: "示例权限", Code: "example:perm"})
	return nil
}

// Down 回滚迁移
func (m *{{.ClassName}}) Down(db *gorm.DB) error {
	// TODO: 在这里编写回滚逻辑
	// 示例: db.Where("code = ?", "example:perm").Delete(&model.Permission{})
	return nil
}
`))

	className := strings.ReplaceAll(fileName, "_", " ")
	className = strings.ReplaceAll(className, ".", " ")
	className = strings.Title(className)
	className = strings.ReplaceAll(className, " ", "")

	data := struct {
		ClassName   string
		FileName    string
		Description string
	}{
		ClassName:   className,
		FileName:    fileName,
		Description: name,
	}

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := tpl.Execute(f, data); err != nil {
		return "", err
	}

	return filePath, nil
}
