package main

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"go-fundraising/entity"
	"go-fundraising/handler"
	"go-fundraising/middlewares"
	"go-fundraising/repository"
	"go-fundraising/service"
	webHandler "go-fundraising/web/handler"
	"log"
	"net/http"
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
	COOKIE_SECRET := os.Getenv("COOKIE_SECRET")

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
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	authenticationWebHandler := webHandler.NewAuthenticationHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())
	store := cookie.NewStore([]byte(COOKIE_SECRET))
	router.Use(sessions.Sessions("routerSession", store))
	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")
	// Multiple template render HTML
	router.HTMLRender = loadTemplates("./web/templates")

	router.GET("/", middlewares.AuthorizeAdmin(), func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/campaigns")
	})
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

	userWebRouter := router.Group("/users", middlewares.AuthorizeAdmin())
	{
		userWebRouter.GET("/", userWebHandler.Index)
		userWebRouter.GET("/add", userWebHandler.Add)
		userWebRouter.POST("/store", userWebHandler.Store)
		userWebRouter.GET("/edit/:id", userWebHandler.Edit)
		userWebRouter.POST("/update/:id", userWebHandler.Update)
		userWebRouter.POST("/delete/:id", userWebHandler.Delete)
		userWebRouter.GET("/avatar/:id", userWebHandler.UploadAvatar)
		userWebRouter.POST("/avatar/store/:id", userWebHandler.StoreAvatar)
	}

	campaignWebRouter := router.Group("/campaigns", middlewares.AuthorizeAdmin())
	{
		campaignWebRouter.GET("/", campaignWebHandler.Index)
		campaignWebRouter.GET("/add", campaignWebHandler.Add)
		campaignWebRouter.POST("/store", campaignWebHandler.Store)
		campaignWebRouter.GET("/image/:id", campaignWebHandler.UploadImage)
		campaignWebRouter.POST("/image/store/:id", campaignWebHandler.StoreImage)
		campaignWebRouter.GET("/edit/:id", campaignWebHandler.Edit)
		campaignWebRouter.POST("/update/:id", campaignWebHandler.Update)
		campaignWebRouter.GET("/show/:id", campaignWebHandler.Show)
	}
	transactionWebRouter := router.Group("/transactions", middlewares.AuthorizeAdmin())
	{
		transactionWebRouter.GET("/", transactionWebHandler.Index)
	}

	authenticationWebRouter := router.Group("/auth")
	{
		authenticationWebRouter.GET("/login", authenticationWebHandler.LoginIndex)
		authenticationWebRouter.POST("/login", authenticationWebHandler.LoginStore)
		authenticationWebRouter.POST("/logout", authenticationWebHandler.Logout)
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
