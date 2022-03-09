package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-fundraising/entity"
	"go-fundraising/handler"
	"go-fundraising/repository"
	"go-fundraising/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
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
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	userRouter := router.Group("/api/v1/users")
	{
		userRouter.POST("/register", userHandler.RegisterUser)
		userRouter.POST("/login", userHandler.LoginUser)
	}
	router.Run(":5000")
}
