package models

import (
	"time"
)

type Feedback struct {
	ID                      string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ResponseID              string    `gorm:"not null;index"`
	Clarity                 string    `validate:"required"`
	Tone                    string    `validate:"required"`
	Relevance               string    `validate:"required"`
	OverallScore            string    `validate:"required"`
	NextQuestion            string    `validate:"required"`
	NextQuestionDifficuilty string    `validate:"required,oneof=EASY MEDIUM HARD"`
	Suggestion              string    `validate:"required"`
	Explanation             string    `gorm:"type:text" validate:"required"`
	CreatedAt               time.Time `gorm:"autoCreateTime"`
}
