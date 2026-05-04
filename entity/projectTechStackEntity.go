package entity

import (
	"gorm.io/gorm"
)

type ProjectTechStack struct {
	gorm.Model
	ProjectID uint
	Project   *Project `gorm:"foreignKey:ProjectID"`

	Name string `gorm:"type:varchar(250)"`
}
