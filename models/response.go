package models

import "time"

type Response struct {
	ResponseID string           `gorm:"column:responseid;type:uuid;default:gen_random_uuid();primaryKey"`
	SessionID  string           `gorm:"not null" validate:"required"`
	Session    InterviewSession `gorm:"foreignKey:SessionID" validate:"-"`

	UserQuestionID string        `gorm:"not null;unique" validate:"required"`
	UserQuestion   *UserQuestion `gorm:"foreignKey:UserQuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" validate:"-"`

	Answer      string    `validate:"required"`
	SubmittedAt time.Time `gorm:"autoCreateTime"`
	Feedback    Feedback  `gorm:"foreignKey:ResponseID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" validate:"-"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
