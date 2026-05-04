package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cld *cloudinary.Cloudinary

func InitCloudinary() {

    log.Println("CLOUD_NAME:", os.Getenv("CLOUDINARY_CLOUD_NAME"))
    log.Println("API_KEY:", os.Getenv("CLOUDINARY_API_KEY"))

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatal("❌ Cloudinary init failed:", err)
	}

	Cld = cld
	log.Println("✅ Cloudinary connected successfully")
}