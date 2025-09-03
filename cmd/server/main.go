package main

import (
	"log"
	"os"

	"go-shorturl/internal/db"
	"go-shorturl/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 初始化資料庫
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// 建立 Fiber 應用程式
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// 中間件
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// API 路由
	api := app.Group("/api")
	api.Post("/shorten", handlers.ShortenURL)
	api.Get("/stats/:short_code", handlers.GetStats)

	// 重定向路由 (必須放在最後，因為它會匹配所有路徑)
	app.Get("/shorturl/:short_code", handlers.RedirectURL)
	app.Get("/:short_code", handlers.RedirectURL) // 保持向後兼容

	// 健康檢查端點
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Short URL service is running",
		})
	})

	// 根路徑
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Short URL Service",
			"version": "1.0.0",
			"endpoints": fiber.Map{
				"POST /api/shorten": "Create a short URL",
				"GET /:short_code": "Redirect to original URL",
				"GET /api/stats/:short_code": "Get URL statistics",
				"GET /health": "Health check",
			},
		})
	})

	// 啟動伺服器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
