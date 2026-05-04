package config

import (
    "fmt"
    "log"

    "github.com/Kongsavanh12/resume-backend/entity"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {
    return db
}

func ConnecttionDB() {
    dsn := "host=localhost user=postgres password=1234 dbname=resume port=5432 sslmode=disable TimeZone=Asia/Bangkok"
    if dsn == "" {
        log.Fatal("DATABASE_URL is not set")
    }

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to PostgreSQL: %v", err)
    }

    db = database
    fmt.Println("✅ Database connected successfully")
}

func SetupDatabase() {
    err := db.AutoMigrate(
        &entity.Role{},                    
        &entity.User{},             
        &entity.Resume{},           
        &entity.Contact{},          
        &entity.Project{},          
        &entity.ProjectTechStack{}, 
        &entity.ProjectFeature{},   
    )

	// Role
	Role := entity.Role{
		RoleName: "Admin",
	}
	db.FirstOrCreate(&Role, &entity.Role{RoleName: "Admin"})

	hashedPassword, err := HashPassword("123")
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}

	//User
	User := entity.User{
		RoleID: 1,
		UserName: "Kongsavanh",
		UserEmail: "Kongsavanh.ch@gmail.com",
		UserPassword: hashedPassword,
	}
	db.FirstOrCreate(&User, &entity.User{UserEmail: "Kongsavanh.ch@gmail.com"})

    if err != nil {
        log.Fatal("❌ Migration failed:", err)
    }

    fmt.Println("✅ Database migrated successfully")
}

func InitDB() {
    ConnecttionDB()
    SetupDatabase()
}