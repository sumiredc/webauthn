package repository

import (
	"encoding/base64"
	"wa/model"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type CredentialRepository struct {
	db *gorm.DB
}

func NewCredentialRepository(db *gorm.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

func (r *CredentialRepository) Create(userID uint64, c *webauthn.Credential) (*model.Credential, error) {
	println(string(c.ID))
	credential := &model.Credential{
		CredentialID: base64.StdEncoding.EncodeToString(c.ID),
		UserID:       userID,
		Object:       *c,
	}

	return credential, r.db.Model(credential).Create(credential).Error
}
