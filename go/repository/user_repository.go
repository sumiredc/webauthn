package repository

import (
	"github.com/sumiredc/webauthn/model"

	"github.com/go-webauthn/webauthn/webauthn"
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

func (r *UserRepository) GetWithCredential(credentialID, userHandle string) (*model.User, error) {
	type userWithCredential struct {
		model.User
		Credential webauthn.Credential `gorm:"column:data;serializer:json"`
	}

	uWithC := &userWithCredential{}

	err := r.db.
		Table(model.User{}.TableName()+" AS u").
		Joins("INNER JOIN credential AS c ON u.id = c.user_id").
		Where("c.credential_id = ?", credentialID).
		Where("u.user_handle = ?", userHandle).
		Select([]string{
			"u.*",
			"c.data",
		}).
		First(uWithC).
		Error

	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:          uWithC.ID,
		Username:    uWithC.Username,
		UserHandle:  uWithC.UserHandle,
		Credentials: []webauthn.Credential{uWithC.Credential},
	}, nil
}
