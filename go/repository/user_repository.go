package repository

import (
	"wa/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var cnt int64

	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&cnt).Error
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
	user := &model.User{}
	err := r.db.Model(user).First(user, id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(username, userHandle string) (*model.User, error) {
	user := &model.User{
		Username:   username,
		UserHandle: userHandle,
	}

	return user, r.db.Model(user).Create(user).Error
}
