package project

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Kongsavanh12/resume-backend/config"
	"github.com/Kongsavanh12/resume-backend/entity"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// POST /projects
func CreateProject(c *gin.Context) {
	userIDStr := c.PostForm("UserID")
	userID, _ := strconv.Atoi(userIDStr)

	// Upload รูป
	coverURL := ""
	file, err := c.FormFile("image")
	if err == nil {
		src, _ := file.Open()
		defer src.Close()

		fileName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
		uploadResult, err := config.Cld.Upload.Upload(
			context.Background(),
			src,
			uploader.UploadParams{
				Folder:   "projects",
				PublicID: fileName,
			},
		)
		if err == nil {
			coverURL = uploadResult.SecureURL
		}
	}

	displayOrder, _ := strconv.Atoi(c.PostForm("DisplayOrder"))

	// รับ TechStacks ที่ user พิมพ์เอง
	techStackNames := c.PostFormArray("TechStackNames[]")

	var techStacks []entity.ProjectTechStack
	for _, name := range techStackNames {
		if name == "" {
			continue
		}
		techStacks = append(techStacks, entity.ProjectTechStack{
			Name: name,
		})
	}

	// รับ Features เป็น array เช่น Features[]=User auth&Features[]=Dark mode
	featuresStr := c.PostFormArray("Features[]")
	orderIndexsStr := c.PostFormArray("OrderIndexs[]")

	// สร้าง ProjectFeature
	var features []entity.ProjectFeature
	for i, desc := range featuresStr {
		orderIndex := i + 1 // default order
		if i < len(orderIndexsStr) {
			orderIndex, _ = strconv.Atoi(orderIndexsStr[i])
		}
		features = append(features, entity.ProjectFeature{
			Description: desc,
			OrderIndex:  orderIndex,
		})
	}

	project := entity.Project{
		UserID:           uint(userID),
		ProjectTitle:     c.PostForm("ProjectTitle"),
		ProjectSlug:      c.PostForm("ProjectSlug"),
		ShortDescription: c.PostForm("ShortDescription"),
		FullDescription:  c.PostForm("FullDescription"),
		CoverImageURL:    coverURL,
		GitHubURL:        c.PostForm("GitHubURL"),
		DemoURL:          c.PostForm("DemoURL"),
		Status:           c.PostForm("Status"),
		DisplayOrder:     displayOrder,
		ProjectTechStack: techStacks,
		ProjectFeatures:  features,
	}

	log.Println("coverURL:", coverURL)

	if err := config.DB().Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": project})
}

// GET /projects
func GetAllProjects(c *gin.Context) {
	var projects []entity.Project

	if err := config.DB().Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// Get /projects/:id
func GetProjectByID(c *gin.Context) {
	var project entity.Project
	id := c.Param("id")

	if err := config.DB().First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// PUT /projects/:id
func UpdateProject(c *gin.Context) {
	id := c.Param("id")

	// เช็คว่ามีอยู่ใน db
	var project entity.Project
	if err := config.DB().First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Upload รูปใหม่ถ้ามีการส่งมา
	coverURL := project.CoverImageURL // ← ใช้รูปเดิมก่อน
	file, err := c.FormFile("image")
	if err == nil { // มีไฟล์ใหม่ส่งมา
		src, _ := file.Open()
		defer src.Close()

		fileName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
		uploadResult, err := config.Cld.Upload.Upload(
			context.Background(),
			src,
			uploader.UploadParams{
				Folder:   "projects",
				PublicID: fileName,
			},
		)
		if err == nil {
			coverURL = uploadResult.SecureURL // ← อัพเดทรูปใหม่
		}
	}

	// อัพเดท TechStacks ถ้ามีส่งมา
	techStackNames := c.PostFormArray("TechStackNames[]")
	if len(techStackNames) > 0 {
		// ลบ TechStacks เดิมออกก่อน
		config.DB().Where("project_id = ?", project.ID).Delete(&entity.ProjectTechStack{})

		// สร้างใหม่
		var techStacks []entity.ProjectTechStack
		for _, name := range techStackNames {
			if name == "" {
				continue
			}
			techStacks = append(techStacks, entity.ProjectTechStack{
				ProjectID: project.ID,
				Name:      name,
			})
		}
		config.DB().Create(&techStacks)
	}

	// อัพเดท Features ถ้ามีส่งมา
	featuresStr := c.PostFormArray("Features[]")
	orderIndexsStr := c.PostFormArray("OrderIndexs[]")
	if len(featuresStr) > 0 {
		// ลบ Features เดิมออกก่อน
		config.DB().Where("project_id = ?", project.ID).Delete(&entity.ProjectFeature{})

		// สร้างใหม่
		var features []entity.ProjectFeature
		for i, desc := range featuresStr {
			orderIndex := i + 1
			if i < len(orderIndexsStr) {
				orderIndex, _ = strconv.Atoi(orderIndexsStr[i])
			}
			features = append(features, entity.ProjectFeature{
				ProjectID:   project.ID,
				Description: desc,
				OrderIndex:  orderIndex,
			})
		}
		config.DB().Create(&features)
	}

	// อัพเดท field อื่นๆ ถ้ามีส่งมา
	displayOrder, _ := strconv.Atoi(c.PostForm("DisplayOrder"))

	updates := map[string]interface{}{}

	if v := c.PostForm("ProjectTitle"); v != "" {
		updates["project_title"] = v
	}
	if v := c.PostForm("ProjectSlug"); v != "" {
		updates["project_slug"] = v
	}
	if v := c.PostForm("ShortDescription"); v != "" {
		updates["short_description"] = v
	}
	if v := c.PostForm("FullDescription"); v != "" {
		updates["full_description"] = v
	}
	if v := c.PostForm("GitHubURL"); v != "" {
		updates["git_hub_url"] = v
	}
	if v := c.PostForm("DemoURL"); v != "" {
		updates["demo_url"] = v
	}
	if v := c.PostForm("Status"); v != "" {
		updates["status"] = v
	}
	if displayOrder != 0 {
		updates["display_order"] = displayOrder
	}
	updates["cover_image_url"] = coverURL

	if err := config.DB().Model(&project).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ดึงข้อมูลใหม่พร้อม relations
	config.DB().
		Preload("ProjectTechStack").
		Preload("ProjectFeatures").
		First(&project, id)

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// DELETE /projects/:id
func DeleteProject(c *gin.Context) {
	var project entity.Project
	id := c.Param("id")

	// เช็คว่ามีอยู่ใน db
	if err := config.DB().First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	//Soft delete
	if err := config.DB().Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Project deleted successfully"})
}
