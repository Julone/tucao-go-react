package models

// User struct to describe User object.
type User struct {
	ID           int    `gorm:"column:id;type:bigint;not null;primaryKey;auto_increment" json:"id" `
	Username     string `gorm:"column:username" json:"username" validate:"required,lte=255"`
	PasswordHash string `gorm:"column:password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	UserStatus   int    `gorm:"column:user_status" json:"user_status" validate:"required,len=1"`
	UserRole     string `gorm:"column:user_role" json:"user_role" validate:"required,lte=25"`
	BaseDbTime
}
type UserRepo struct {
	Curd[User]
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}
