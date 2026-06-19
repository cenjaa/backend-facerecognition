package http

import (
	"face-recognition-fyp/domain"
	"face-recognition-fyp/facerpca/delivery/http/handler"
	"face-recognition-fyp/facerpca/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/spf13/viper"
)

func APIKeyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip auth for health check
		if c.Path() == "/health" {
			return c.Next()
		}

		apiKey := c.Get("X-API-Key")
		expectedKey := viper.GetString("server.api_key")

		if expectedKey != "" && apiKey != expectedKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}
		return c.Next()
	}
}

// SetupRouter creates the Fiber app and registers all routes.
func SetupRouter(uc domain.FaceRPCAUsecase, redisRepo domain.FaceRPCARedisRepository) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:   "PNM Attendance API v2.0",
		BodyLimit: 10 * 1024 * 1024, // 10MB limit for image uploads/syncs
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,X-API-Key,Authorization",
	}))

	// Swagger UI at /swagger/*
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check (no auth)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "mode": "vps"})
	})

	// API key middleware for all /api routes
	app.Use(APIKeyMiddleware())

	h := handler.NewFaceRPCAHandler(uc)

	api := app.Group("/api")

	// Public / Pi Routes (Still protected by API Key)
	api.Post("/attendance", h.LogAttendance)
	api.Post("/attendance/bulk", h.LogAttendanceBulk)
	api.Get("/user_status_today/:id", h.GetUserStatusToday)
	api.Get("/user_name/:id", h.GetUserName)
	api.Post("/verify_pin", h.VerifyPin)
	api.Post("/login_admin", h.LoginAdmin)
	api.Post("/infer", h.InferFace)

	// Admin Routes (Protected by JWT + Redis)
	admin := api.Group("/", middleware.AuthMiddleware(redisRepo))

	// User Management
	admin.Get("/users", h.GetAllUsers)
	admin.Put("/users/:id", h.UpdateUser)
	admin.Delete("/users/:id", h.DeleteUser)
	admin.Post("/create_user", h.CreateUser)
	admin.Post("/dataset/:id", h.UploadDataset)

	// History / Jira
	admin.Get("/attendance_history", h.GetAttendanceHistory)
	admin.Get("/recent_jira", h.GetRecentJira)
	admin.Get("/jira_history", h.GetJiraHistory)
	admin.Get("/kpi_accumulation", h.GetKpiAccumulation)

	// Dashboard
	admin.Get("/dashboard_stats", h.GetDashboardStats)
	admin.Get("/attendance_frequency", h.GetAttendanceFrequency)
	admin.Get("/admin_dashboard/:id", h.GetAdminDashboard)
	admin.Get("/recent_presence", h.GetRecentPresence)

	// Training
	admin.Post("/train_model", h.StartTraining)
	admin.Get("/train_status", h.GetTrainStatus)

	// Exports
	admin.Get("/export_attendance", h.ExportAttendance)
	admin.Get("/export_jira", h.ExportJira)

	// Logout
	admin.Post("/logout_admin", h.LogoutAdmin)

	return app
}
