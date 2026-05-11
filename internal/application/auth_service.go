package application

import (
	"errors"
	"gin-admin-base/internal/domain/model"
	"gin-admin-base/internal/domain/repository"
)

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CheckPermission(userID uint, requiredPermission string) (bool, error) {
	perms, err := s.repo.GetPermissionsByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, p := range perms {
		// 如果用户拥有 "*"（所有权限），则任何权限检查都通过
		if p.Code == "*" {
			return true, nil
		}
		if p.Code == requiredPermission {
			return true, nil
		}
	}
	return false, errors.New("permission denied")
}

// GetUserByUsername 根据用户名获取用户
func (s *AuthService) GetUserByUsername(username string) (*model.User, error) {
	return s.repo.GetUserByUsername(username)
}

// GetUserPermissionCodes 获取用户的所有权限代码列表
func (s *AuthService) GetUserPermissionCodes(userID uint) ([]string, error) {
	perms, err := s.repo.GetPermissionsByUserID(userID)
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(perms))
	for _, p := range perms {
		codes = append(codes, p.Code)
	}
	return codes, nil
}

// ClearUserPermissionCache 清除用户权限缓存
func (s *AuthService) ClearUserPermissionCache(userID uint) error {
	return s.repo.ClearUserPermissionCache(userID)
}
