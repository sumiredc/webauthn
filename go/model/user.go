package model

import "github.com/go-webauthn/webauthn/webauthn"

type User struct {
	ID         uint64 `gorm:"column:id;primaryKey" json:"id"`
	Username   string `gorm:"column:username" json:"username"`
	UserHandle string `gorm:"column:user_handle" json:"user_handle"`
}

var _ webauthn.User = (*User)(nil)

func (User) TableName() string {
	return "user"
}

func (u *User) WebAuthnID() []byte {
	return []byte(u.UserHandle)
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.Username
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return []webauthn.Credential{}
}
