package models

type Domain struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	Name        string     `gorm:"unique;not null" validate:"required"`
	Description *string    `gorm:"type:text" validate:"required"`
	Questions   []Question `gorm:"foreignKey:DomainID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
