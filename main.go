package main

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"go-fundraising/entity"
	"go-fundraising/handler"
	"go-fundraising/middlewares"
	"go-fundraising/repository"
	"go-fundraising/service"
	webHandler "go-fundraising/web/handler"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"

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
	db.AutoMigrate(&entity.User{}, &entity.Campaign{}, &entity.CampaignImage{}, &entity.Transaction{})

	userRepository := repository.NewUserRepository(db)
	campaignRepository := repository.NewCampaignRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	jwtService := service.NewJWTService()
	userService := service.NewUserService(userRepository)
	campaignService := service.NewCampaignService(campaignRepository)
	paymentService := service.NewPaymentService()
	transactionService := service.NewTransactionService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, jwtService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	userWebHandler := webHandler.NewUserHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images/avatar", "./images/avatars")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")
	// Multiple template render HTML
	router.HTMLRender = loadTemplates("./web/templates")

	router.POST("/midtrans/callback", transactionHandler.GetNotification)

	userRouter := router.Group("/api/v1/users")
	{
		userRouter.POST("/register", userHandler.RegisterUser)
		userRouter.POST("/login", userHandler.LoginUser)
		userRouter.POST("/check-email", userHandler.CheckEmailAvaibility)
		userRouter.POST("/avatars", middlewares.AuthorizeToken(jwtService, userService), userHandler.UploadAvatar)
		userRouter.GET("/profile", middlewares.AuthorizeToken(jwtService, userService), userHandler.MyProfile)
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

	transactionRouter := router.Group("/api/v1/campaigns")
	{
		transactionRouter.GET("/:id/transactions",
			middlewares.AuthorizeToken(jwtService, userService),
			transactionHandler.GetTransactionsByCampaignID)
	}

	transactionUserRouter := router.Group("/api/v1/transactions")
	{
		transactionUserRouter.GET("/",
			middlewares.AuthorizeToken(jwtService, userService),
			transactionHandler.GetTransactionsByUserID)
		transactionUserRouter.POST("/",
			middlewares.AuthorizeToken(jwtService, userService),
			transactionHandler.CreateTransaction)
	}

	userWebRouter := router.Group("/users")
	{
		userWebRouter.GET("/", userWebHandler.Index)
		userWebRouter.GET("/add", userWebHandler.Add)
		userWebRouter.POST("/store", userWebHandler.Store)
	}
	router.Run(":5000")
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// Load template from all the files inside layouts folder
	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		log.Fatalln(err.Error())
	}
	// Load template from all the files/folders inside templates folder
	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
