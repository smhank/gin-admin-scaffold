package repository

import "gin-admin-base/internal/domain/model"

type AuthRepository interface {
	GetUserByUsername(username string) (*model.User, error)
	GetPermissionsByUserID(userID uint) ([]model.Permission, error)
	ClearUserPermissionCache(userID uint) error
}
