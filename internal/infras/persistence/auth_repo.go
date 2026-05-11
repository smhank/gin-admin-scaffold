package persistence

import (
	"encoding/json"
	"fmt"
	"gin-admin-base/internal/domain/model"
	"gin-admin-base/internal/domain/repository"
	"gin-admin-base/internal/infras/cache"
	"time"

	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) repository.AuthRepository {
	return &authRepo{db: db}
}

// cacheKey 生成用户权限缓存键
func cacheKey(userID uint) string {
	return fmt.Sprintf("user_permissions:%d", userID)
}

func (r *authRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).Preload("Roles.Permissions").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) GetPermissionsByUserID(userID uint) ([]model.Permission, error) {
	// 尝试从 Redis 缓存获取
	key := cacheKey(userID)
	if cache.RedisClient != nil {
		// 先检查 Redis 是否可用
		if err := cache.RedisClient.Ping(cache.Ctx).Err(); err == nil {
			data, err := cache.RedisClient.Get(cache.Ctx, key).Bytes()
			if err == nil {
				var perms []model.Permission
				if json.Unmarshal(data, &perms) == nil {
					return perms, nil
				}
			}
		}
	}

	// 缓存未命中，从数据库查询
	var user model.User
	err := r.db.Preload("Roles.Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	var perms []model.Permission
	permMap := make(map[uint]bool) // 去重
	for _, role := range user.Roles {
		for _, p := range role.Permissions {
			if !permMap[p.ID] {
				permMap[p.ID] = true
				perms = append(perms, p)
			}
		}
	}

	// 写入 Redis 缓存
	if cache.RedisClient != nil {
		if err := cache.RedisClient.Ping(cache.Ctx).Err(); err == nil {
			if data, err := json.Marshal(perms); err == nil {
				cache.RedisClient.Set(cache.Ctx, key, data, 30*time.Minute)
			}
		}
	}

	return perms, nil
}

// ClearUserPermissionCache 清除用户权限缓存
func (r *authRepo) ClearUserPermissionCache(userID uint) error {
	if cache.RedisClient != nil {
		key := cacheKey(userID)
		return cache.RedisClient.Del(cache.Ctx, key).Err()
	}
	return nil
}
