package models

import "time"

// models/user.go
type User struct {
	ID            string             `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string             `validate:"required,min=2"`
	Email         string             `gorm:"unique;not null" validate:"required,email"`
	Password      string             `gorm:"column:password_hash;not null" validate:"required"`
	Role          string             `gorm:"type:varchar(20);default:USER"`
	Sessions      []InterviewSession `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" validate:"-"`
	UserDomains   []UserDomain       `gorm:"foreignKey:UserID" validate:"-"`
	UserQuestions []UserQuestion     `gorm:"foreignKey:UserID" validate:"-"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
