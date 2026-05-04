package upload

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/Kongsavanh12/resume-backend/config"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func uploadToCloudinary(c *gin.Context, folder string) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))

	uploadResult, err := config.Cld.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			Folder:   folder,
			PublicID: fileName,
		},
	)
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
