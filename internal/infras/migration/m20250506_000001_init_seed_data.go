package migration

import (
	"fmt"

	"gin-admin-base/internal/domain/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// M20250506000001InitSeedData 初始化种子数据
type M20250506000001InitSeedData struct {
	BaseMigration
}

// NewM20250506000001InitSeedData 创建迁移
func NewM20250506000001InitSeedData() *M20250506000001InitSeedData {
	return &M20250506000001InitSeedData{
		BaseMigration: NewBaseMigration("m20250506_000001_init_seed_data"),
	}
}

// Up 执行迁移
func (m *M20250506000001InitSeedData) Up(db *gorm.DB) error {
	// 检查是否已有数据
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("  数据库已有用户数据，跳过种子数据初始化")
		return nil
	}

	// 1. 创建权限
	permissions := []model.Permission{
		{Name: "用户列表", Code: "user:list"},
		{Name: "创建用户", Code: "user:create"},
		{Name: "编辑用户", Code: "user:edit"},
		{Name: "删除用户", Code: "user:delete"},
		{Name: "角色列表", Code: "role:list"},
		{Name: "创建角色", Code: "role:create"},
		{Name: "编辑角色", Code: "role:edit"},
		{Name: "删除角色", Code: "role:delete"},
		{Name: "菜单列表", Code: "menu:list"},
		{Name: "创建菜单", Code: "menu:create"},
		{Name: "编辑菜单", Code: "menu:edit"},
		{Name: "删除菜单", Code: "menu:delete"},
		{Name: "API列表", Code: "api:list"},
		{Name: "创建API", Code: "api:create"},
		{Name: "编辑API", Code: "api:edit"},
		{Name: "删除API", Code: "api:delete"},
	}
	for i := range permissions {
		if err := createIfNotExists(db, &permissions[i], map[string]interface{}{"code": permissions[i].Code}); err != nil {
			return fmt.Errorf("创建权限 %s 失败: %w", permissions[i].Name, err)
		}
	}
	fmt.Println("  ✓ 权限数据已就绪")

	// 2. 创建角色
	roles := []model.Role{
		{Name: "超级管理员", Status: 1},
		{Name: "管理员", Status: 1},
		{Name: "普通用户", Status: 1},
		{Name: "测试角色", Status: 1},
	}
	for i := range roles {
		if err := createIfNotExists(db, &roles[i], map[string]interface{}{"name": roles[i].Name}); err != nil {
			return fmt.Errorf("创建角色 %s 失败: %w", roles[i].Name, err)
		}
	}
	fmt.Println("  ✓ 角色数据已就绪")

	// 3. 创建用户（密码: admin123）
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	users := []model.User{
		{
			Username: "admin",
			Password: string(hashedPassword),
			NickName: "Mr.奇森",
			Email:    "admin@example.com",
			Phone:    "13800138000",
			Status:   1,
		},
		{
			Username: "a303176530",
			Password: string(hashedPassword),
			NickName: "用户1",
			Email:    "user1@example.com",
			Phone:    "13800138001",
			Status:   1,
		},
	}
	for i := range users {
		if err := createIfNotExists(db, &users[i], map[string]interface{}{"username": users[i].Username}); err != nil {
			return fmt.Errorf("创建用户 %s 失败: %w", users[i].Username, err)
		}
	}
	fmt.Println("  ✓ 用户数据已就绪")

	// 4. 关联用户和角色（admin 关联超级管理员）
	var adminUser model.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return fmt.Errorf("查询admin用户失败: %w", err)
	}
	var superAdminRole model.Role
	if err := db.Where("name = ?", "超级管理员").First(&superAdminRole).Error; err != nil {
		return fmt.Errorf("查询超级管理员角色失败: %w", err)
	}

	var roleCount int64
	db.Table("user_roles").Where("user_id = ? AND role_id = ?", adminUser.ID, superAdminRole.ID).Count(&roleCount)
	if roleCount == 0 {
		if err := db.Model(&adminUser).Association("Roles").Append(&superAdminRole); err != nil {
			return fmt.Errorf("关联用户角色失败: %w", err)
		}
	}
	fmt.Println("  ✓ 用户角色关联已就绪")

	// 5. 关联角色和权限
	var allPerms []model.Permission
	db.Find(&allPerms)
	perms := make([]interface{}, len(allPerms))
	for i, p := range allPerms {
		perms[i] = p
	}
	if err := db.Model(&superAdminRole).Association("Permissions").Replace(perms...); err != nil {
		return fmt.Errorf("关联角色权限失败: %w", err)
	}
	fmt.Println("  ✓ 角色权限关联已就绪")

	// 6. 创建菜单
	menus := []model.Menu{
		{Name: "仪表盘", Icon: "Odometer", Path: "/dashboard", ParentID: nil, Sort: 1, Status: 1},
		{Name: "系统管理", Icon: "Setting", Path: "/system", ParentID: nil, Sort: 2, Status: 1},
		{Name: "角色管理", Icon: "User", Path: "/roles", ParentID: uintPtr(2), Sort: 1, Status: 1},
		{Name: "菜单管理", Icon: "Menu", Path: "/menus", ParentID: uintPtr(2), Sort: 2, Status: 1},
		{Name: "API管理", Icon: "Connection", Path: "/apis", ParentID: uintPtr(2), Sort: 3, Status: 1},
		{Name: "用户管理", Icon: "UserFilled", Path: "/users", ParentID: uintPtr(2), Sort: 4, Status: 1},
	}
	for i := range menus {
		if err := createIfNotExists(db, &menus[i], map[string]interface{}{"name": menus[i].Name, "path": menus[i].Path}); err != nil {
			return fmt.Errorf("创建菜单 %s 失败: %w", menus[i].Name, err)
		}
	}
	fmt.Println("  ✓ 菜单数据已就绪")

	// 7. 创建API路径
	apiPaths := []model.ApiPath{
		{Name: "用户列表", Path: "/api/users", Method: "GET", Desc: "获取用户列表"},
		{Name: "创建用户", Path: "/api/users", Method: "POST", Desc: "创建用户"},
		{Name: "编辑用户", Path: "/api/users/:id", Method: "PUT", Desc: "编辑用户"},
		{Name: "删除用户", Path: "/api/users/:id", Method: "DELETE", Desc: "删除用户"},
		{Name: "角色列表", Path: "/api/roles", Method: "GET", Desc: "获取角色列表"},
		{Name: "创建角色", Path: "/api/roles", Method: "POST", Desc: "创建角色"},
		{Name: "编辑角色", Path: "/api/roles/:id", Method: "PUT", Desc: "编辑角色"},
		{Name: "删除角色", Path: "/api/roles/:id", Method: "DELETE", Desc: "删除角色"},
		{Name: "设置权限", Path: "/api/roles/:id/permissions", Method: "PUT", Desc: "设置角色权限"},
		{Name: "分配用户", Path: "/api/roles/assign", Method: "POST", Desc: "分配角色给用户"},
		{Name: "菜单列表", Path: "/api/menus", Method: "GET", Desc: "获取菜单列表"},
		{Name: "创建菜单", Path: "/api/menus", Method: "POST", Desc: "创建菜单"},
		{Name: "编辑菜单", Path: "/api/menus/:id", Method: "PUT", Desc: "编辑菜单"},
		{Name: "删除菜单", Path: "/api/menus/:id", Method: "DELETE", Desc: "删除菜单"},
		{Name: "API列表", Path: "/api/paths", Method: "GET", Desc: "获取API列表"},
		{Name: "创建API", Path: "/api/paths", Method: "POST", Desc: "创建API"},
		{Name: "编辑API", Path: "/api/paths/:id", Method: "PUT", Desc: "编辑API"},
		{Name: "删除API", Path: "/api/paths/:id", Method: "DELETE", Desc: "删除API"},
		{Name: "权限列表", Path: "/api/permissions", Method: "GET", Desc: "获取权限列表"},
		{Name: "创建权限", Path: "/api/permissions", Method: "POST", Desc: "创建权限"},
		{Name: "编辑权限", Path: "/api/permissions/:id", Method: "PUT", Desc: "编辑权限"},
		{Name: "删除权限", Path: "/api/permissions/:id", Method: "DELETE", Desc: "删除权限"},
		{Name: "登录", Path: "/api/login", Method: "POST", Desc: "用户登录"},
	}
	for i := range apiPaths {
		if err := createIfNotExists(db, &apiPaths[i], map[string]interface{}{"path": apiPaths[i].Path, "method": apiPaths[i].Method}); err != nil {
			return fmt.Errorf("创建API路径 %s 失败: %w", apiPaths[i].Name, err)
		}
	}
	fmt.Println("  ✓ API路径数据已就绪")

	return nil
}

// Down 回滚迁移
func (m *M20250506000001InitSeedData) Down(db *gorm.DB) error {
	// 删除所有种子数据
	db.Exec("DELETE FROM api_paths")
	db.Exec("DELETE FROM menus")
	db.Exec("DELETE FROM role_permissions")
	db.Exec("DELETE FROM user_roles")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM permissions")
	fmt.Println("  ✓ 种子数据已清除")
	return nil
}

// createIfNotExists 如果不存在则创建记录
func createIfNotExists(db *gorm.DB, model interface{}, condition map[string]interface{}) error {
	var count int64
	db.Model(model).Where(condition).Count(&count)
	if count > 0 {
		return nil
	}
	return db.Create(model).Error
}

// uintPtr 辅助函数
func uintPtr(id uint) *uint {
	return &id
}
