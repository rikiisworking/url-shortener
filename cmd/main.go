package main

import (
	"log"

	"gitgub.com/rikiisworking/url-shortener/internal/config"
	"gitgub.com/rikiisworking/url-shortener/internal/handler"
	"gitgub.com/rikiisworking/url-shortener/internal/service"
	"gitgub.com/rikiisworking/url-shortener/internal/storage"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg := config.Load()

	// Storages
	postgresRepo, err := storage.NewPostgresRepo(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}
	defer postgresRepo.Close()

	redisRepo := storage.NewRedisRepo(cfg)

	// Service & Handler
	urlService := service.NewURLService(postgresRepo, redisRepo, cfg.ShortCodeLength)
	urlHandler := handler.NewURLHandler(urlService)

	app := fiber.New(fiber.Config{
		AppName: "URL Shortener - Go Fiber + pgx",
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(func(c fiber.Ctx) error {
		c.Locals("port", cfg.ServerPort)
		return c.Next()
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})
	app.Post("/api/shorten", urlHandler.Shorten)
	app.Get("/:shortCode", urlHandler.Redirect)

	log.Printf("Server starting on :%s", cfg.ServerPort)
	log.Fatal(app.Listen(":" + cfg.ServerPort))
}
