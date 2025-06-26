package models

import "time"

// models/admin_user.go
type AdminUser struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string `gorm:"unique;not null" validate:"required,email"`
	Username  string `gorm:"unique;not null" validate:"required,min=3"`
	Password  string `gorm:"column:password_hash;not null" validate:"required,min=6"`
	Role      string `gorm:"type:varchar(20);default:EDITOR" validate:"required,role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
