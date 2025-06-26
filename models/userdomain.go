package models

import "time"

type UserDomain struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	UserID        string         `gorm:"not null" validate:"-"`
	DomainID      uint           `gorm:"not null" validate:"required"`
	Domain        Domain         `gorm:"foreignKey:DomainID" validate:"-"`
	User          User           `gorm:"foreignKey:UserID" validate:"-"`
	UserQuestions []UserQuestion `gorm:"foreignKey:UserDomainID" validate:"-"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
