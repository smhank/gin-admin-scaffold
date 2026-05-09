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
