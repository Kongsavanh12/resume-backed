package entity

import (
	"gorm.io/gorm"
)

type Role struct{
	gorm.Model
	RoleName 	string

	Users 		[]User			`gorm:"foreignKey:RoleID"`
} 