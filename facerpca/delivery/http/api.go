package http

import (
	"face-recognition-fyp/domain"
	"face-recognition-fyp/facerpca/delivery/http/handler"

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
func SetupRouter(uc domain.FaceRPCAUsecase) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:   "PNM Attendance API v2.0",
		BodyLimit: 10 * 1024 * 1024, // 10MB limit for image uploads/syncs
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,X-API-Key",
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

	// User Management
	api.Get("/users", h.GetAllUsers)
	api.Put("/users/:id", h.UpdateUser)
	api.Delete("/users/:id", h.DeleteUser)
	api.Post("/create_user", h.CreateUser)
	api.Post("/dataset/:id", h.UploadDataset)
	api.Get("/user_name/:id", h.GetUserName)

	// Attendance
	api.Post("/attendance", h.LogAttendance)
	api.Post("/attendance/bulk", h.LogAttendanceBulk)
	api.Get("/user_status_today/:id", h.GetUserStatusToday)
	api.Get("/recent_presence", h.GetRecentPresence)
	api.Get("/attendance_history", h.GetAttendanceHistory)
	api.Post("/set_cuti", h.SetCuti)

	// Dashboard
	api.Get("/dashboard_stats", h.GetDashboardStats)
	api.Get("/attendance_frequency", h.GetAttendanceFrequency)
	api.Get("/admin_dashboard/:id", h.GetAdminDashboard)

	// Auth
	api.Post("/verify_pin", h.VerifyPin)

	// KPI / Jira
	api.Get("/recent_jira", h.GetRecentJira)
	api.Get("/jira_history", h.GetJiraHistory)
	api.Get("/kpi_accumulation", h.GetKpiAccumulation)

	// Training
	api.Post("/train_model", h.StartTraining)
	api.Get("/train_status", h.GetTrainStatus)

	return app
}
