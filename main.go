package main

import (
	"fmt"
	"go-fundraising/entity"
	"go-fundraising/handler"
	"go-fundraising/middlewares"
	"go-fundraising/repository"
	"go-fundraising/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatalln("Error while loading env file")
	}
	DB_HOST := os.Getenv("DB_HOST")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_USER := os.Getenv("DB_USER")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", DB_HOST, DB_USER, DB_PASSWORD, DB_PORT, DB_NAME)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("Database connected")
	db.AutoMigrate(&entity.User{})

	userRepository := repository.NewUserRepository(db)

	jwtService := service.NewJWTService()
	userService := service.NewUserService(userRepository)

	userHandler := handler.NewUserHandler(userService, jwtService)

	router := gin.Default()

	userRouter := router.Group("/api/v1/users")
	{
		userRouter.POST("/register", userHandler.RegisterUser)
		userRouter.POST("/login", userHandler.LoginUser)
		userRouter.POST("/check-email", userHandler.CheckEmailAvaibility)
		userRouter.POST("/avatars", middlewares.AuthorizeToken(jwtService, userService), userHandler.UploadAvatar)

	}
	router.Run(":5000")
}
