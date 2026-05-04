package main

import (
	//  "net/http"

	"github.com/Kongsavanh12/resume-backend/config"
	"github.com/Kongsavanh12/resume-backend/controller/login"
	"github.com/Kongsavanh12/resume-backend/controller/project"
	"github.com/Kongsavanh12/resume-backend/controller/resume"
	"github.com/Kongsavanh12/resume-backend/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.InitDB()
	config.InitCloudinary()

	r := gin.Default()

	r.Use(CORSMiddleware())
	r.GET("/resume/:id", resume.GetResumeByID)
	r.POST("/login", login.AddLogin)

	authorized := r.Group("")
	authorized.Use(middlewares.Authorizes())
	{
		// Project routes
		projectRoutes := authorized.Group("/projects")
		{
			projectRoutes.POST("", project.CreateProject)
			projectRoutes.GET("", project.GetAllProjects)
			projectRoutes.GET("/:id", project.GetProjectByID)
			projectRoutes.PUT("/:id", project.UpdateProject)
			projectRoutes.DELETE("/:id", project.DeleteProject)
		}

		// Resume routes
		resumeRoutes := authorized.Group("/resumes")
		{
			resumeRoutes.POST("", resume.CreateResume)
			resumeRoutes.PUT("/:id", resume.UpdateResume)
		}
		
	}

	r.Run(":8080")

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
