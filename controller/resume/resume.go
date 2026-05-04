package resume

import (
	"encoding/json"
	"net/http"

	"github.com/Kongsavanh12/resume-backend/config"
	"github.com/Kongsavanh12/resume-backend/entity"
	"github.com/gin-gonic/gin"
)

type ContactInput struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Linkedin    string `json:"linkedin"`
	LinkedinUrl string `json:"linkedinUrl"`
	Github      string `json:"github"`
	GithubUrl   string `json:"githubUrl"`
}

type ResumeInput struct {
	UserID    uint         `json:"UserID"`
	JobTitle  string       `json:"JobTitle"`
	Contact   ContactInput `json:"contact"`
	Objective string       `json:"objective"`

	Skills          interface{} `json:"skills"`
	Experience      interface{} `json:"experience"`
	Education       interface{} `json:"education"`
	ExtraCurricular interface{} `json:"extracurricular"`
	Languages       interface{} `json:"languages"`
}

func toJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func formatResume(resume entity.Resume) gin.H {
	var skills, experience, education, extracurricular, languages interface{}
	json.Unmarshal([]byte(resume.Skills), &skills)
	json.Unmarshal([]byte(resume.Experience), &experience)
	json.Unmarshal([]byte(resume.Education), &education)
	json.Unmarshal([]byte(resume.ExtraCurricular), &extracurricular)
	json.Unmarshal([]byte(resume.Languages), &languages)

	return gin.H{
		"ID":       resume.ID,
		"UserID":   resume.UserID,
		"JobTitle": resume.JobTitle,
		"contact": gin.H{
			"email":       resume.Email,
			"phone":       resume.Phone,
			"linkedin":    resume.Linkedin,
			"linkedinUrl": resume.LinkedinUrl,
			"github":      resume.Github,
			"githubUrl":   resume.GithubUrl,
		},
		"objective":       resume.Objective,
		"skills":          skills,
		"experience":      experience,
		"education":       education,
		"extracurricular": extracurricular,
		"languages":       languages,
	}
}

// GET /resume/:id — Guest เข้าถึงได้
func GetResumeByID(c *gin.Context) {
	var resume entity.Resume
	id := c.Param("id")

	if err := config.DB().First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": formatResume(resume)})
}

// POST /resumes — User เท่านั้น
func CreateResume(c *gin.Context) {
	var input ResumeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เช็คว่ามี resume อยู่แล้วมั้ย
	var existing entity.Resume
	if err := config.DB().Where("user_id = ?", input.UserID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume already exists"})
		return
	}

	resume := entity.Resume{
		UserID:          input.UserID,
		JobTitle:        input.JobTitle,
		Email:           input.Contact.Email,
		Phone:           input.Contact.Phone,
		Linkedin:        input.Contact.Linkedin,
		LinkedinUrl:     input.Contact.LinkedinUrl,
		Github:          input.Contact.Github,
		GithubUrl:       input.Contact.GithubUrl,
		Objective:       input.Objective,
		Skills:          toJSON(input.Skills),
		Experience:      toJSON(input.Experience),
		Education:       toJSON(input.Education),
		ExtraCurricular: toJSON(input.ExtraCurricular),
		Languages:       toJSON(input.Languages),
	}

	if err := config.DB().Create(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": formatResume(resume)})
}

// PUT /resumes/:id — User เท่านั้น
func UpdateResume(c *gin.Context) {
	var resume entity.Resume
	id := c.Param("id")

	if err := config.DB().First(&resume, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	var input ResumeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.JobTitle != "" {
		updates["job_title"] = input.JobTitle
	}
	if input.Contact.Email != "" {
		updates["email"] = input.Contact.Email
	}
	if input.Contact.Phone != "" {
		updates["phone"] = input.Contact.Phone
	}
	if input.Contact.Linkedin != "" {
		updates["linkedin"] = input.Contact.Linkedin
	}
	if input.Contact.LinkedinUrl != "" {
		updates["linkedin_url"] = input.Contact.LinkedinUrl
	}
	if input.Contact.Github != "" {
		updates["github"] = input.Contact.Github
	}
	if input.Contact.GithubUrl != "" {
		updates["github_url"] = input.Contact.GithubUrl
	}
	if input.Objective != "" {
		updates["objective"] = input.Objective
	}
	if input.Skills != nil {
		updates["skills"] = toJSON(input.Skills)
	}
	if input.Experience != nil {
		updates["experience"] = toJSON(input.Experience)
	}
	if input.Education != nil {
		updates["education"] = toJSON(input.Education)
	}
	if input.ExtraCurricular != nil {
		updates["extra_curricular"] = toJSON(input.ExtraCurricular)
	}
	if input.Languages != nil {
		updates["languages"] = toJSON(input.Languages)
	}

	if err := config.DB().Model(&resume).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ดึงข้อมูลใหม่หลัง update
	config.DB().First(&resume, id)

	c.JSON(http.StatusOK, gin.H{"data": formatResume(resume)})
}