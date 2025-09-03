package handler

import (
	"net/http"

	"go-shorturl/pkg/db"
	"go-shorturl/pkg/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func init() {
	// 初始化資料庫
	if err := db.InitDB(); err != nil {
		panic(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
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
	api.Get("/stats/:short_code", handlers.GetStats)

	// 處理請求
	adaptor.FiberApp(app).ServeHTTP(w, r)
}