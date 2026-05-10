package handler

import (
	"face-recognition-fyp/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type FaceRPCAHandler struct {
	uc domain.FaceRPCAUsecase
}

func NewFaceRPCAHandler(uc domain.FaceRPCAUsecase) *FaceRPCAHandler {
	return &FaceRPCAHandler{
		uc: uc,
	}
}

// LogAttendance godoc
// @Summary      Log single attendance from Raspi
// @Tags         Attendance
// @Param        body body object true "user_id, presence_time, confidence, latency_ms, status_id"
// @Success      200 {object} map[string]string
// @Router       /api/attendance [post]
func (h *FaceRPCAHandler) LogAttendance(c *fiber.Ctx) error {
	var body domain.BulkAttendanceRecord
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid request"})
	}
	if err := h.uc.LogAttendanceRaw(c.Context(), body); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success"})
}

// LogAttendanceBulk godoc
// @Summary      Bulk sync attendance from Raspi
// @Tags         Attendance
// @Param        body body object true "records array"
// @Success      200 {object} map[string]interface{}
// @Router       /api/attendance/bulk [post]
func (h *FaceRPCAHandler) LogAttendanceBulk(c *fiber.Ctx) error {
	var body struct {
		Records []domain.BulkAttendanceRecord `json:"records"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid request"})
	}
	synced, err := h.uc.LogAttendanceBulk(c.Context(), body.Records)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "synced": synced})
}

// GetUserStatusToday godoc
// @Summary      Get user's latest attendance status today
// @Tags         Attendance
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]int
// @Router       /api/user_status_today/{id} [get]
func (h *FaceRPCAHandler) GetUserStatusToday(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("id"))
	status, err := h.uc.GetUserStatusToday(c.Context(), userID)
	if err != nil || status == nil {
		return c.JSON(fiber.Map{"status_id": 0})
	}
	return c.JSON(fiber.Map{"status_id": *status})
}

// GetRecentPresence godoc
// @Summary      Get 10 most recent attendance records
// @Tags         Dashboard
// @Success      200 {object} map[string]interface{}
// @Router       /api/recent_presence [get]
func (h *FaceRPCAHandler) GetRecentPresence(c *fiber.Ctx) error {
	records, err := h.uc.GetRecentPresence(c.Context())
	if err != nil || records == nil {
		return c.JSON(fiber.Map{"records": []any{}})
	}
	return c.JSON(fiber.Map{"records": records})
}

// GetAttendanceHistory godoc
// @Summary      Paginated attendance history
// @Tags         History
// @Param        page     query int false "Page number" default(1)
// @Param        per_page query int false "Items per page" default(10)
// @Success      200 {object} map[string]interface{}
// @Router       /api/attendance_history [get]
func (h *FaceRPCAHandler) GetAttendanceHistory(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	result, err := h.uc.GetAttendanceHistory(c.Context(), page, perPage)
	if err != nil {
		return c.JSON(fiber.Map{"records": []any{}, "total": 0, "page": page, "per_page": perPage, "total_pages": 0})
	}
	return c.JSON(result)
}

// GetAllUsers godoc
// @Summary      List all users
// @Tags         User Management
// @Success      200 {object} map[string]interface{}
// @Router       /api/users [get]
func (h *FaceRPCAHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.uc.GetAllUsers(c.Context())
	if err != nil {
		return c.JSON(fiber.Map{"users": []any{}})
	}
	if users == nil {
		users = []domain.UserListRow{}
	}
	return c.JSON(fiber.Map{"users": users})
}

// UpdateUser godoc
// @Summary      Update user info
// @Tags         User Management
// @Param        id   path int    true "User ID"
// @Param        body body object true "name, nip, active"
// @Success      200 {object} map[string]string
// @Router       /api/users/{id} [put]
func (h *FaceRPCAHandler) UpdateUser(c *fiber.Ctx) error {
	var body struct {
		Name   string `json:"name"`
		NIP    string `json:"nip"`
		Active *bool  `json:"active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid body"})
	}
	userID, _ := strconv.Atoi(c.Params("id"))
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	if err := h.uc.UpdateUser(c.Context(), userID, body.Name, body.NIP, active); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success"})
}

// DeleteUser godoc
// @Summary      Delete user
// @Tags         User Management
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]string
// @Router       /api/users/{id} [delete]
func (h *FaceRPCAHandler) DeleteUser(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("id"))
	if err := h.uc.DeleteUser(c.Context(), userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User dan data wajah berhasil dihapus"})
}

// CreateUser godoc
// @Summary      Create a new user
// @Tags         User Management
// @Param        body body object true "name, nip, email"
// @Success      200 {object} map[string]interface{}
// @Router       /api/create_user [post]
func (h *FaceRPCAHandler) CreateUser(c *fiber.Ctx) error {
	var body struct {
		Name  string `json:"name"`
		NIP   string `json:"nip"`
		Email string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid request"})
	}
	userID, err := h.uc.CreateUser(c.Context(), body.Name, body.NIP, body.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "user_id": userID})
}

// @Summary      Upload dataset for a user
// @Tags         User Management
// @Param        id   path int true "User ID"
// @Param        files formData file true "Images"
// @Success      200 {object} map[string]interface{}
// @Router       /api/dataset/{id} [post]
func (h *FaceRPCAHandler) UploadDataset(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid user ID"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Failed to parse form"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "No files uploaded"})
	}

	if err := h.uc.UploadDataset(c.Context(), userID, files); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Dataset uploaded successfully"})
}

// GetUserName godoc
// @Summary      Get user name by ID
// @Tags         User Management
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]interface{}
// @Router       /api/user_name/{id} [get]
func (h *FaceRPCAHandler) GetUserName(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("id"))
	name, err := h.uc.GetUserName(c.Context(), userID)
	if err != nil {
		return c.JSON(fiber.Map{"name": nil})
	}
	return c.JSON(fiber.Map{"name": name})
}

// VerifyPin godoc
// @Summary      Verify admin PIN (face-first auth step 2)
// @Tags         Auth
// @Param        body body object true "user_id, pin"
// @Success      200 {object} map[string]interface{}
// @Router       /api/verify_pin [post]
func (h *FaceRPCAHandler) VerifyPin(c *fiber.Ctx) error {
	var body struct {
		UserID int    `json:"user_id"`
		Pin    string `json:"pin"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid request"})
	}
	valid, name, err := h.uc.VerifyAdminPin(c.Context(), body.UserID, body.Pin)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	if !valid {
		return c.JSON(fiber.Map{"status": "fail", "message": "Invalid PIN or not an admin"})
	}
	return c.JSON(fiber.Map{"status": "success", "name": name})
}

// GetDashboardStats godoc
// @Summary      Today's attendance counts by status
// @Tags         Dashboard
// @Success      200 {object} domain.DashboardStats
// @Router       /api/dashboard_stats [get]
func (h *FaceRPCAHandler) GetDashboardStats(c *fiber.Ctx) error {
	stats, err := h.uc.GetDashboardStats(c.Context())
	if err != nil {
		return c.JSON(fiber.Map{"menunggu": 0, "hadir": 0, "tidak_hadir": 0, "telat": 0, "keluar": 0, "cuti": 0, "updated_at": ""})
	}
	return c.JSON(stats)
}

// SetCuti godoc
// @Summary      Mark a user as on leave (Cuti) for today
// @Tags         Attendance
// @Param        body body object true "user_id"
// @Success      200 {object} map[string]string
// @Router       /api/set_cuti [post]
func (h *FaceRPCAHandler) SetCuti(c *fiber.Ctx) error {
	var body struct {
		UserID int `json:"user_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "Invalid request"})
	}
	if err := h.uc.SetUserCuti(c.Context(), body.UserID); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User marked as cuti"})
}

// GetAttendanceFrequency godoc
// @Summary      Hourly login frequency chart data (05:00-22:00)
// @Tags         Dashboard
// @Success      200 {object} domain.HourlyFrequency
// @Router       /api/attendance_frequency [get]
func (h *FaceRPCAHandler) GetAttendanceFrequency(c *fiber.Ctx) error {
	freq, err := h.uc.GetAttendanceFrequency(c.Context())
	if err != nil {
		return c.JSON(fiber.Map{"labels": []string{}, "data": []int{}})
	}
	return c.JSON(freq)
}

// GetRecentJira godoc
// @Summary      10 most recent Jira KPI updates
// @Tags         KPI
// @Success      200 {object} map[string]interface{}
// @Router       /api/recent_jira [get]
func (h *FaceRPCAHandler) GetRecentJira(c *fiber.Ctx) error {
	records, err := h.uc.GetRecentJira(c.Context())
	if err != nil || records == nil {
		return c.JSON(fiber.Map{"records": []any{}})
	}
	return c.JSON(fiber.Map{"records": records})
}

// GetJiraHistory godoc
// @Summary      Paginated Jira ticket history
// @Tags         History
// @Param        page     query int false "Page" default(1)
// @Param        per_page query int false "Per page" default(10)
// @Success      200 {object} map[string]interface{}
// @Router       /api/jira_history [get]
func (h *FaceRPCAHandler) GetJiraHistory(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	result, err := h.uc.GetJiraHistory(c.Context(), page, perPage)
	if err != nil {
		return c.JSON(fiber.Map{"records": []any{}, "total": 0, "page": page, "per_page": perPage, "total_pages": 0})
	}
	return c.JSON(result)
}

// GetKpiAccumulation godoc
// @Summary      Paginated KPI accumulation
// @Tags         History
// @Param        page     query int false "Page" default(1)
// @Param        per_page query int false "Per page" default(10)
// @Success      200 {object} map[string]interface{}
// @Router       /api/kpi_accumulation [get]
func (h *FaceRPCAHandler) GetKpiAccumulation(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	result, err := h.uc.GetKpiAccumulation(c.Context(), page, perPage)
	if err != nil {
		return c.JSON(fiber.Map{"records": []any{}, "total": 0, "page": page, "per_page": perPage, "total_pages": 0})
	}
	return c.JSON(result)
}

// StartTraining godoc
// @Summary      Trigger model training via ML microservice
// @Tags         Training
// @Success      200 {object} map[string]string
// @Router       /api/train_model [post]
func (h *FaceRPCAHandler) StartTraining(c *fiber.Ctx) error {
	if err := h.uc.TrainModel(c.Context()); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "started"})
}

// GetTrainStatus godoc
// @Summary      Get training progress from ML microservice
// @Tags         Training
// @Success      200 {object} domain.TrainStatus
// @Router       /api/train_status [get]
func (h *FaceRPCAHandler) GetTrainStatus(c *fiber.Ctx) error {
	status, err := h.uc.GetTrainStatus(c.Context())
	if err != nil {
		return c.JSON(fiber.Map{"status": "idle", "progress": 0, "message": ""})
	}
	return c.JSON(status)
}
