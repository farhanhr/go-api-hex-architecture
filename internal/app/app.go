package app

import (
	"context"
	"gonews/config"
	"gonews/internal/adapter/handler"
	"gonews/internal/adapter/imagekit"
	"gonews/internal/adapter/repository"
	"gonews/internal/core/service"
	"gonews/lib/auth"
	"gonews/lib/middleware"
	"gonews/lib/pagination"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}

	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatalf("Error connecting to temp dir: %v", err)
		return
	}

	// Imagekit
	ikAdapter := imagekit.NewImageKitAdapter(cfg)
	

	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(cfg)
	_ = pagination.NewPagination()

	//repository
	authRepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)


	//service
	authService := service.NewAuthService(authRepo, cfg, jwt)
	categoryService := service.NewCategoryService(categoryRepo)
	contentService := service.NewContentService(contentRepo, cfg, ikAdapter)
	userService := service.NewUserService(userRepo)

	//handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] %{ip} %{status} - %{latency} %{method} %{path}\n",
	}))

	if os.Getenv("APP_ENV") != "production" {
		cfg := swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.json",
			Path:     "docs",
			Title:    "Swagger API Docs",
		}
		
		app.Use(swagger.New(cfg))
	}


	api := app.Group("/api")
	api.Post("/login", authHandler.Login)

	adminApp := api.Group("/admin")
	adminApp.Use(middlewareAuth.CheckToken())

	//category
	categoryApp := adminApp.Group("/categories")
	categoryApp.Get("/", categoryHandler.GetCategories)
	categoryApp.Post("/", categoryHandler.CreateCategory)
	categoryApp.Get("/:categoryID", categoryHandler.GetCategoryByID)
	categoryApp.Put("/:categoryID", categoryHandler.EditCategory)
	categoryApp.Delete("/:categoryID", categoryHandler.DeleteCategory)
	
	//content
	contentApp := adminApp.Group("/contents")
	contentApp.Get("/", contentHandler.GetContents) 
	contentApp.Post("/", contentHandler.CreateContent) 
	contentApp.Get("/:contentID", contentHandler.GetContentById) 
	contentApp.Put("/:contentID", contentHandler.UpdateContent) 
	contentApp.Delete("/:contentID", contentHandler.DeleteContent) 
	contentApp.Post("/upload-image", contentHandler.UploadImageR2)

	//user 
	userApp := adminApp.Group("/users")
	userApp.Get("/profile", userHandler.GetUserByID)
	userApp.Put("/update-password", userHandler.UpdatePassword)

	//fe
	feApp := api.Group("/fe")
	feApp.Get("/categories", categoryHandler.GetCategoryFE)
	feApp.Get("/contents", contentHandler.GetContentWithQuery)
	feApp.Get("/contents/:contentID", contentHandler.GetContentDetail)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatalf("error when starting server: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("server shutdown on 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}