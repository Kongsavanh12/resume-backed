package entity

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	UserID uint
	User   *User `gorm:"foreignKey:UserID"`

	SenderName  string `gorm:"type:varchar(250)"`
	SenderEmail string `gorm:"type:varchar(250)"`
	Subject     string `gorm:"type:varchar(250)"`
	Message     string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(20)"`
}
