package entity

import (
	"gorm.io/gorm"
)

type ProjectFeature struct {
	gorm.Model
	ProjectID 		uint
	Project			*Project		`gorm:"foreignKey:ProjectID"`

	Description 	string
	OrderIndex 		int
}