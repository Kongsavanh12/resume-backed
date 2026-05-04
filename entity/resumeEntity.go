package entity

import (
	"gorm.io/gorm"
)

type Resume struct {
	gorm.Model
	UserID uint
	User   *User `gorm:"foreignKey:UserID"`

	// Contact
	Email       string `gorm:"type:varchar(250)"`
	Phone       string `gorm:"type:varchar(250)"`
	Linkedin    string `gorm:"type:varchar(250)"`
	LinkedinUrl string `gorm:"type:varchar(250)"`
	Github      string `gorm:"type:varchar(250)"`
	GithubUrl   string `gorm:"type:varchar(250)"`

	JobTitle  string `gorm:"type:varchar(250)"`
	Objective string `gorm:"type:text"`

	Skills          string `gorm:"type:jsonb"` 
	Experience      string `gorm:"type:jsonb"` 
	Education       string `gorm:"type:jsonb"` 
	ExtraCurricular string `gorm:"type:jsonb"` 
	Languages       string `gorm:"type:jsonb"` 
}