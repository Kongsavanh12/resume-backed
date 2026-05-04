package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	RoleID 			uint
	Role			*Role			`gorm:"foreignKey:RoleID"`

	UserName 		string 
	UserEmail 		string
	UserPassword 	string

	Resumes 		[]Resume		`gorm:"foreignKey:UserID"`
	Contacts 		[]Contact		`gorm:"foreignKey:UserID"`
	Projects 		[]Project		`gorm:"foreignKey:UserID"`
}