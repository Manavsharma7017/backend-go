package models

import "time"

type UserQuestion struct {
	ID             string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Questionnumber int        `gorm:"not null" validate:"required"`
	Text           string     `gorm:"not null" validate:"required"`
	Explanation    string     `gorm:"type:text" validate:"required"`
	Difficulty     Difficulty `gorm:"type:varchar(10)" validate:"required,oneof=EASY MEDIUM HARD"`
	UserDomainID   uint       `gorm:"not null" validate:"required"`
	UserDomain     UserDomain `gorm:"foreignKey:UserDomainID" validate:"-"`

	SessionID string           `gorm:"not null;column:sessionid" validate:"required"`
	Session   InterviewSession `gorm:"foreignKey:SessionID" validate:"-"`
	UserID    string           `gorm:"not null"`
	User      User             `gorm:"foreignKey:UserID" validate:"-"`

	Responses *Response `gorm:"foreignKey:UserQuestionID" validate:"-"` // use pointer
	CreatedAt time.Time
	UpdatedAt time.Time
}
