package repository

import (
	"context"

	"github.com/Amierza/pos-broissant/entity"
	"gorm.io/gorm"
)

type (
	UserRepository interface {
		CheckEmailOrPhoneNumber(ctx context.Context, tx *gorm.DB, email, phoneNumber string) (entity.User, bool, error)
		RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	}
	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRespository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CheckEmailOrPhoneNumber(ctx context.Context, tx *gorm.DB, email, phoneNumber string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("email = ? OR phone_number = ?", email, phoneNumber).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

func (r *userRepository) RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}
