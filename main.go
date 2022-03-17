package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-fundraising/entity"
	"go-fundraising/handler"
	"go-fundraising/middlewares"
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
	db.AutoMigrate(&entity.User{}, &entity.Campaign{}, &entity.CampaignImage{})

	userRepository := repository.NewUserRepository(db)
	campaignRepository := repository.NewCampaignRepository(db)

	jwtService := service.NewJWTService()
	userService := service.NewUserService(userRepository)
	campaignService := service.NewCampaignService(campaignRepository)

	userHandler := handler.NewUserHandler(userService, jwtService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images/avatar", "./images/avatars")

	userRouter := router.Group("/api/v1/users")
	{
		userRouter.POST("/register", userHandler.RegisterUser)
		userRouter.POST("/login", userHandler.LoginUser)
		userRouter.POST("/check-email", userHandler.CheckEmailAvaibility)
		userRouter.POST("/avatars", middlewares.AuthorizeToken(jwtService, userService), userHandler.UploadAvatar)

	}

	campaignRouter := router.Group("/api/v1/campaigns")
	{
		campaignRouter.GET("/", campaignHandler.GetCampaigns)
		campaignRouter.POST("/",
			middlewares.AuthorizeToken(jwtService, userService),
			campaignHandler.CreateCampaign)
		campaignRouter.GET("/:id", campaignHandler.GetCampaignByID)
		campaignRouter.PUT("/:id", middlewares.AuthorizeToken(jwtService, userService), campaignHandler.UpdateCampaign)
		//campaignRouter.GET("/:slug", campaignHandler.GetCampaignBySlug)
	}

	campaignImageRouter := router.Group("/api/v1/campaigns-images")
	{
		campaignImageRouter.POST("/",
			middlewares.AuthorizeToken(jwtService, userService),
			campaignHandler.UploadCampaignImage)
	}
	router.Run(":5000")
}
