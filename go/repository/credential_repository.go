package repository

import (
	"encoding/base64"

	"github.com/sumiredc/webauthn/model"

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
	credential := &model.Credential{
		CredentialID: base64.StdEncoding.EncodeToString(c.ID),
		UserID:       userID,
		Data:         *c,
	}

	return credential, r.db.
		Model(credential).
		Create(credential).
		Error
}

func (r *CredentialRepository) Update(c *webauthn.Credential) (*model.Credential, error) {
	credentialID := base64.StdEncoding.EncodeToString(c.ID)
	credential := &model.Credential{
		Data: *c,
	}

	return credential, r.db.
		Model(credential).
		Where("credential_id = ?", credentialID).
		Updates(credential).
		Error
}
