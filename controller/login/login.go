package login

import (
	"net/http"

	"github.com/Kongsavanh12/resume-backend/config"
	"github.com/Kongsavanh12/resume-backend/entity"
	"github.com/Kongsavanh12/resume-backend/services"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UserEmail    string `json:"user_email" binding:"required,email"`
	UserPassword string `json:"user_password" binding:"required"`
}

func AddLogin(c *gin.Context) {
	var req LoginRequest

	// ---------- validate input ----------
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	db := config.DB()

	var user entity.User

	if err := db.Preload("Role").
		Where("user_email = ?", req.UserEmail).
		First(&user).Error; err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !config.CheckPasswordHash(req.UserPassword, user.UserPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ---------- generate JWT ----------
	jwtWrapper := services.JwtWrapper{
		SecretKey:       config.JwtSecret,
		Issuer:          config.JwtIssuer,
		ExpirationHours: 24 * 180,
	}

	signedToken, err := jwtWrapper.GenerateToken(user.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error signing token"})
		return
	}

	var roleName string
	if user.Role != nil {
		roleName = user.Role.RoleName
	}

	c.JSON(http.StatusOK, gin.H{
		"token_type": "Bearer",
		"token":      signedToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.UserName,
			"email": user.UserEmail,
			"role":  roleName,
		},
	})
}
