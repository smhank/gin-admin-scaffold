package persistence

import (
	"gin-admin-base/internal/domain/model"
	"gin-admin-base/internal/domain/repository"

	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) repository.AuthRepository {
	return &authRepo{db: db}
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
	var user model.User
	err := r.db.Preload("Roles.Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	var perms []model.Permission
	for _, role := range user.Roles {
		perms = append(perms, role.Permissions...)
	}
	return perms, nil
}
