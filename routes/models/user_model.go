package models

import (
	"github.com/google/uuid"
)

// User struct to describe User object.
type User struct {
	ID           uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Email        string    `db:"email" json:"email" validate:"required,email,lte=255"`
	PasswordHash string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	UserStatus   int       `db:"user_status" json:"user_status" validate:"required,len=1"`
	UserRole     string    `db:"user_role" json:"user_role" validate:"required,lte=25"`
	BaseDbTime
}
type UserRepo struct {
	Curd[User]
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}
