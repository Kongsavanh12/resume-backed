package entity

import(
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	UserID 				uint
	User				*User				`gorm:"foreignKey:UserID"`

    ProjectTitle     string             `gorm:"type:varchar(250)"`
    ProjectSlug      string             `gorm:"type:varchar(250)"`
    ShortDescription string             `gorm:"type:varchar(250)"`
    FullDescription  string             `gorm:"type:text"`
    CoverImageURL    string             `gorm:"type:varchar(250)"`
    GitHubURL        string             `gorm:"type:varchar(250)"`
    DemoURL          string             `gorm:"type:varchar(250)"`
    Status           string             `gorm:"type:varchar(20)"` 
    DisplayOrder     int

	ProjectTechStack 	[]ProjectTechStack 	`gorm:"foreignKey:ProjectID"`
	ProjectFeatures 	[]ProjectFeature 	`gorm:"foreignKey:ProjectID"`
}