package main

import (
	"gin-admin-base/internal/domain/model"

	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/infras/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// 假设数据库连接已在其他地方处理，这里仅生成代码
	g.ApplyBasic(
		model.User{},
		model.Role{},
		model.Permission{},
	)

	g.Execute()
}
