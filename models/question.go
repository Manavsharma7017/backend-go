package models

import "time"

type Question struct {
	ID          string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Text        string     `validate:"required"`
	Difficulty  Difficulty `gorm:"type:varchar(10)" validate:"required,oneof=EASY MEDIUM HARD"`
	DomainID    uint       `gorm:"not null" validate:"required"`
	Explanation string     `gorm:"type:text" validate:"required"`
	Domain      Domain     `gorm:"foreignKey:DomainID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"` // changed to OnDelete:CASCADE
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
