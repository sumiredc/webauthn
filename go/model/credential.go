package model

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

type Credential struct {
	CredentialID string              `gorm:"column:credential_id"`
	UserID       uint64              `gorm:"column:user_id"`
	Data         webauthn.Credential `gorm:"column:data;serializer:json"`
}

func (Credential) TableName() string {
	return "credential"
}
