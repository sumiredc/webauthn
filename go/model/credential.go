package model

import "github.com/go-webauthn/webauthn/webauthn"

type Credential struct {
	CredentialID string              `gorm:"column:credential_id"`
	UserID       uint64              `gorm:"column:user_id"`
	Object       webauthn.Credential `gorm:"serializer:json;column:json"`
}

func (Credential) TableName() string {
	return "credential"
}
