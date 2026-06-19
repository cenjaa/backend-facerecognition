package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"face-recognition-fyp/constant"
	"face-recognition-fyp/domain"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"github.com/xuri/excelize/v2"
)

type FaceRPCAUsecase struct {
	MinioRepo    domain.FaceRPCAMinIORepository
	SQLRepo      domain.FaceRPCASQLRepository
	JiraRepo     domain.FaceRPCAJiraRepository
	RedisRepo    domain.FaceRPCARedisRepository
	MLServiceURL string
}

func NewFaceRPCAUsecase(m domain.FaceRPCAMinIORepository, s domain.FaceRPCASQLRepository, j domain.FaceRPCAJiraRepository, r domain.FaceRPCARedisRepository, mlServiceURL string) *FaceRPCAUsecase {
	return &FaceRPCAUsecase{
		MinioRepo:    m,
		SQLRepo:      s,
		JiraRepo:     j,
		RedisRepo:    r,
		MLServiceURL: strings.TrimRight(mlServiceURL, "/"),
	}
}

func (uc *FaceRPCAUsecase) CreateUser(ctx context.Context, name, nip, email, createdBy, datasetPath string) (int, error) {
	if name == "" || nip == "" {
		return 0, fmt.Errorf("name and nip are required")
	}
	return uc.SQLRepo.CreateUser(ctx, name, nip, email, createdBy, datasetPath)
}

func (uc *FaceRPCAUsecase) UpdateUser(ctx context.Context, userID int, name, nip string, active bool) error {
	return uc.SQLRepo.Update(ctx, userID, name, nip, active)
}

func (uc *FaceRPCAUsecase) DeleteUser(ctx context.Context, userID int) error {
	return uc.SQLRepo.Delete(ctx, userID)
}

func (uc *FaceRPCAUsecase) VerifyAdminPin(ctx context.Context, userID int, pin string) (string, string, error) {
	isAdmin, err := uc.SQLRepo.IsAdmin(ctx, userID)
	if err != nil {
		return "", "", err
	}
	if !isAdmin {
		return "", "", fmt.Errorf("user is not an admin")
	}

	valid, err := uc.SQLRepo.VerifyAdminPin(ctx, userID, pin)
	if err != nil {
		return "", "", err
	}
	if !valid {
		return "", "", fmt.Errorf("invalid pin")
	}

	name, _ := uc.SQLRepo.GetUserNameByID(ctx, userID)

	// Generate JWT for the verified admin
	token, err := uc.generateToken(ctx, userID)
	if err != nil {
		return "", name, err
	}

	return token, name, nil
}

// RecordAttendance handles face recognition login. No manual logoff — clock out is automatic at 20:00.
// Before 09:00 → Status 1 (Hadir). After 09:00 → Status 3 (Telat Hadir).
func (uc *FaceRPCAUsecase) RecordAttendance(ctx context.Context, userID int, mode string, confidence, latency float64) error {
	today := time.Now()

	currentStatusPtr, err := uc.SQLRepo.GetUserStatusToday(ctx, userID, today)
	if err != nil {
		return err
	}

	// Check if user is already logged in (Hadir or Telat)
	if currentStatusPtr != nil {
		cs := *currentStatusPtr
		if cs == constant.StatusHadir || cs == constant.StatusTelat {
			return fmt.Errorf("user already logged in today")
		}
		if cs == constant.StatusKeluar {
			return fmt.Errorf("user already clocked out today")
		}
		if cs == constant.StatusCuti {
			return fmt.Errorf("user is on leave (cuti) today")
		}
	}

	// Determine status based on time of day
	status := constant.StatusHadir // Before 09:00
	if today.Hour() >= 9 {
		status = constant.StatusTelat // After 09:00
	}

	record := domain.AttendanceRecord{
		IDUser:       userID,
		PresenceTime: today,
		Confidence:   domain.Confidence(confidence),
		LatencyMS:    latency,
		Status:       status,
	}
	return uc.SQLRepo.SaveAttendance(ctx, record)
}

func (uc *FaceRPCAUsecase) LogAttendanceRaw(ctx context.Context, rec domain.BulkAttendanceRecord) error {
	return uc.SQLRepo.SaveRaw(ctx, rec)
}

func (uc *FaceRPCAUsecase) LogAttendanceBulk(ctx context.Context, records []domain.BulkAttendanceRecord) (int, error) {
	return uc.SQLRepo.SaveBulk(ctx, records)
}

func (uc *FaceRPCAUsecase) UploadDataset(ctx context.Context, userID int, fileHeaders []*multipart.FileHeader) error {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("dataset_%d", userID))
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // clean up afterwards

	for i, header := range fileHeaders {
		file, err := header.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %d: %w", i, err)
		}

		dstPath := filepath.Join(tempDir, fmt.Sprintf("%d.jpg", i))
		dstFile, err := os.Create(dstPath)
		if err != nil {
			file.Close()
			return fmt.Errorf("failed to create file %s: %w", dstPath, err)
		}

		if _, err := io.Copy(dstFile, file); err != nil {
			dstFile.Close()
			file.Close()
			return fmt.Errorf("failed to write file %s: %w", dstPath, err)
		}
		dstFile.Close()
		file.Close()
	}

	// Upload folder to MinIO
	_, err = uc.MinioRepo.UploadUserFaces(ctx, userID, tempDir)
	return err
}

func (uc *FaceRPCAUsecase) InferFace(ctx context.Context, fileHeader *multipart.FileHeader) (map[string]interface{}, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", uc.MLServiceURL+"/api/infer", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ML service unreachable: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// High-visibility confidence log
	if detected, ok := result["detected"].(bool); ok && detected {
		fmt.Printf("🎯 [FACE] User %v detected with %.2f%% confidence\n", result["user_id"], result["confidence"])
	} else {
		fmt.Printf("👁️ [FACE] No match or low confidence (%.2f%%)\n", result["confidence"])
	}

	return result, nil
}

func (uc *FaceRPCAUsecase) TrainModel(ctx context.Context) error {
	if uc.MLServiceURL == "" {
		return fmt.Errorf("ML service URL not configured")
	}

	resp, err := http.Post(uc.MLServiceURL+"/train", "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to reach ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ML service returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode ML response: %w", err)
	}

	if result.Status == "error" {
		return fmt.Errorf("ML service error: %s", result.Message)
	}

	return nil
}

func (uc *FaceRPCAUsecase) GetTrainStatus(ctx context.Context) (*domain.TrainStatus, error) {
	if uc.MLServiceURL == "" {
		return &domain.TrainStatus{Status: "idle", Progress: 0, Message: "ML service not configured"}, nil
	}

	resp, err := http.Get(uc.MLServiceURL + "/train_status")
	if err != nil {
		return &domain.TrainStatus{Status: "error", Progress: 0, Message: "ML service unreachable"}, nil
	}
	defer resp.Body.Close()

	var status domain.TrainStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return &domain.TrainStatus{Status: "error", Progress: 0, Message: "Invalid ML response"}, nil
	}

	return &status, nil
}

// HandleDailyInit runs at 05:00 AM. Assigns status 0 (Menunggu Kehadiran) to all active users.
func (uc *FaceRPCAUsecase) HandleDailyInit(ctx context.Context) error {
	today := time.Now().Truncate(24 * time.Hour)
	return uc.SQLRepo.InitializeDailyStatus(ctx, today)
}

// HandleMorningCutoff runs at 09:00 AM. Users still at status 0 become status 2 (Tidak Hadir).
func (uc *FaceRPCAUsecase) HandleMorningCutoff(ctx context.Context) error {
	today := time.Now().Truncate(24 * time.Hour)
	now := time.Now()
	return uc.SQLRepo.MarkAbsentUsers(ctx, today, now)
}

// HandleEveningCutoff runs at 20:00. Forces all Hadir (1) and Telat (3) users to Keluar (4),
// then accumulates their KPI.
func (uc *FaceRPCAUsecase) HandleEveningCutoff(ctx context.Context) error {
	today := time.Now().Truncate(24 * time.Hour)

	affectedUsers, err := uc.SQLRepo.ForceLogoffUsers(ctx, today)
	if err != nil {
		return err
	}

	for _, userID := range affectedUsers {
		_ = uc.AccumulateKPI(ctx, userID, today)
	}

	return nil
}

// PollJiraTasks polls Jira for all currently logged-in users today.
func (uc *FaceRPCAUsecase) PollJiraTasks(ctx context.Context) error {
	today := time.Now().Truncate(24 * time.Hour)

	loggedInUsers, err := uc.SQLRepo.GetLoggedInUsers(ctx, today)
	if err != nil {
		return err
	}

	if len(loggedInUsers) == 0 {
		return nil
	}

	return uc.pollJiraForUsers(ctx, loggedInUsers, today)
}

// PollJiraTasksForDate polls Jira for users who worked on a specific date (used for overnight catch-up at 07:00 AM).
func (uc *FaceRPCAUsecase) PollJiraTasksForDate(ctx context.Context, date time.Time) error {
	users, err := uc.SQLRepo.GetUsersWhoWorkedOnDate(ctx, date)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	return uc.pollJiraForUsers(ctx, users, date)
}

// pollJiraForUsers is the shared Jira polling logic.
func (uc *FaceRPCAUsecase) pollJiraForUsers(ctx context.Context, userIDs []int, date time.Time) error {
	for _, userID := range userIDs {
		email, err := uc.SQLRepo.GetUserEmailByID(ctx, userID)
		if err != nil || email == "" {
			continue
		}

		name, err := uc.SQLRepo.GetUserNameByID(ctx, userID)
		if err != nil {
			name = "Unknown"
		}

		tasks, err := uc.JiraRepo.GetTasksForUser(ctx, email)
		if err != nil || len(tasks) == 0 {
			continue
		}

		for _, task := range tasks {
			kpi := domain.KPIHarian{
				IDTiketJira:   task.TicketKey,
				StoryJira:     task.StoryJira,
				IDStatusTiket: task.Status,
				NamaTiket:     name,
				IDUser:        userID,
				Date:          date,
			}
			_ = uc.SQLRepo.UpsertKpiHarian(ctx, kpi)
		}
	}

	return nil
}

// SetUserCuti allows an admin to mark a user as on leave (Cuti) for today.
func (uc *FaceRPCAUsecase) SetUserCuti(ctx context.Context, userID int) error {
	today := time.Now().Truncate(24 * time.Hour)

	// Check current status — don't override if already Hadir/Telat/Keluar
	currentStatusPtr, err := uc.SQLRepo.GetUserStatusToday(ctx, userID, time.Now())
	if err != nil {
		return err
	}
	if currentStatusPtr != nil {
		cs := *currentStatusPtr
		if cs == constant.StatusHadir || cs == constant.StatusTelat {
			return fmt.Errorf("user already logged in today, cannot set cuti")
		}
		if cs == constant.StatusKeluar {
			return fmt.Errorf("user already clocked out today, cannot set cuti")
		}
		if cs == constant.StatusCuti {
			return fmt.Errorf("user is already on cuti today")
		}
	}

	return uc.SQLRepo.SetUserCuti(ctx, userID, today)
}

// AccumulateKPI maps to src/usecase/attendance/accumulate_kpi.py
func (uc *FaceRPCAUsecase) AccumulateKPI(ctx context.Context, userID int, today time.Time) error {
	hasDone, err := uc.SQLRepo.HasDoneTasks(ctx, userID, today)
	if err != nil {
		return err
	}

	if !hasDone {
		return nil
	}

	doneCount, err := uc.SQLRepo.GetDoneTasksCount(ctx, userID, today)
	if err != nil {
		return err
	}

	return uc.SQLRepo.UpdateAkumulasiKPI(ctx, userID, doneCount)
}

func (uc *FaceRPCAUsecase) GetAllUsers(ctx context.Context) ([]domain.UserListRow, error) {
	return uc.SQLRepo.GetAll(ctx)
}

func (uc *FaceRPCAUsecase) GetUserName(ctx context.Context, userID int) (string, error) {
	return uc.SQLRepo.GetUserNameByID(ctx, userID)
}

func (uc *FaceRPCAUsecase) GetUserStatusToday(ctx context.Context, userID int) (*int, error) {
	return uc.SQLRepo.GetUserStatusToday(ctx, userID, time.Now())
}

func (uc *FaceRPCAUsecase) GetRecentPresence(ctx context.Context) ([]domain.RecentPresence, error) {
	return uc.SQLRepo.GetRecentPresence(ctx, 10)
}

func (uc *FaceRPCAUsecase) GetAttendanceHistory(ctx context.Context, page, perPage int) (*domain.PaginatedResult[domain.AttendanceHistoryRow], error) {
	return uc.SQLRepo.GetHistory(ctx, page, perPage)
}

func (uc *FaceRPCAUsecase) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	return uc.SQLRepo.GetDashboardStats(ctx)
}

func (uc *FaceRPCAUsecase) GetAttendanceFrequency(ctx context.Context) (*domain.HourlyFrequency, error) {
	return uc.SQLRepo.GetHourlyFrequency(ctx)
}

func (uc *FaceRPCAUsecase) GetRecentJira(ctx context.Context) ([]domain.JiraHistoryRow, error) {
	return uc.SQLRepo.GetRecentJira(ctx, 10)
}

func (uc *FaceRPCAUsecase) GetJiraHistory(ctx context.Context, page, perPage int) (*domain.PaginatedResult[domain.JiraHistoryRow], error) {
	return uc.SQLRepo.GetJiraHistory(ctx, page, perPage)
}

func (uc *FaceRPCAUsecase) GetKpiAccumulation(ctx context.Context, page, perPage int) (*domain.PaginatedResult[domain.KpiAccumulationRow], error) {
	return uc.SQLRepo.GetKpiAccumulation(ctx, page, perPage)
}

func (uc *FaceRPCAUsecase) GetAdminDashboard(ctx context.Context, adminID int) (*domain.AdminDashboard, error) {
	adminName, err := uc.SQLRepo.GetUserNameByID(ctx, adminID)
	if err != nil {
		adminName = "Admin"
	}

	totalUsers, err := uc.SQLRepo.GetTotalUserCount(ctx)
	if err != nil {
		totalUsers = 0
	}

	todayAttendees, err := uc.SQLRepo.GetTodayAttendeeCount(ctx)
	if err != nil {
		todayAttendees = 0
	}

	modelTS, _ := uc.MinioRepo.GetModelTimestamp(ctx)
	datasetTS, _ := uc.MinioRepo.GetDatasetTimestamp(ctx)

	// Needs retrain if dataset is newer than model
	needsRetrain := datasetTS > modelTS

	return &domain.AdminDashboard{
		AdminName:      adminName,
		TotalUsers:     totalUsers,
		TodayAttendees: todayAttendees,
		NeedsRetrain:   needsRetrain,
		ModelTimestamp: modelTS,
	}, nil
}

func (uc *FaceRPCAUsecase) LoginAdmin(ctx context.Context, username, password string) (string, *domain.AdminUser, error) {
	admin, err := uc.SQLRepo.GetAdminByCredentials(ctx, username, password)
	if err != nil {
		return "", nil, err
	}

	tokenString, err := uc.generateToken(ctx, admin.IDUser)
	if err != nil {
		return "", nil, err
	}

	return tokenString, admin, nil
}

func (uc *FaceRPCAUsecase) generateToken(ctx context.Context, userID int) (string, error) {
	// Generate JWT
	secret := viper.GetString("server.jwt_secret")
	if secret == "" {
		secret = "pnm-face-recognition-secret-key-2024" // Fallback
	}

	claims := jwt.MapClaims{
		"admin_id": userID,
		"exp":      time.Now().Add(60 * time.Minute).Unix(), // Extended to 60m for comfort
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Store in Redis
	err = uc.RedisRepo.SetSession(ctx, tokenString, userID, 60*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to store session: %w", err)
	}

	return tokenString, nil
}

func (uc *FaceRPCAUsecase) LogoutAdmin(ctx context.Context, token string) error {
	return uc.RedisRepo.DeleteSession(ctx, token)
}

func (uc *FaceRPCAUsecase) ExportAttendanceCSV(ctx context.Context, date string) ([]byte, error) {
	logs, err := uc.SQLRepo.GetAllLogs(ctx, date)
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	sb.WriteString("No,Nama,NIP,Status,Confidence,Latency,Waktu\n")
	for _, l := range logs {
		sb.WriteString(fmt.Sprintf("%d,%s,%s,%s,%d,%s,%s\n", 
			l.No, l.Name, l.NIP, l.Status, l.Confidence, l.Latency, l.Time))
	}
	return []byte(sb.String()), nil
}

func (uc *FaceRPCAUsecase) ExportJiraExcel(ctx context.Context) ([]byte, error) {
	history, err := uc.SQLRepo.GetAllJiraHistory(ctx)
	if err != nil {
		return nil, err
	}
	accumulation, err := uc.SQLRepo.GetAllKpiAccumulation(ctx)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	defer f.Close()

	// Sheet 1: Jira History
	sheet1 := "History Jira"
	f.SetSheetName("Sheet1", sheet1)
	f.SetCellValue(sheet1, "A1", "No")
	f.SetCellValue(sheet1, "B1", "Ticket ID")
	f.SetCellValue(sheet1, "C1", "Nama Tiket")
	f.SetCellValue(sheet1, "D1", "Status")
	f.SetCellValue(sheet1, "E1", "Tanggal")
	for i, h := range history {
		row := i + 2
		f.SetCellValue(sheet1, fmt.Sprintf("A%d", row), h.No)
		f.SetCellValue(sheet1, fmt.Sprintf("B%d", row), h.Ticket)
		f.SetCellValue(sheet1, fmt.Sprintf("C%d", row), h.Name)
		f.SetCellValue(sheet1, fmt.Sprintf("D%d", row), h.Status)
		f.SetCellValue(sheet1, fmt.Sprintf("E%d", row), h.Date)
	}

	// Sheet 2: KPI Accumulation
	sheet2 := "Akumulasi KPI"
	f.NewSheet(sheet2)
	f.SetCellValue(sheet2, "A1", "No")
	f.SetCellValue(sheet2, "B1", "Nama Karyawan")
	f.SetCellValue(sheet2, "C1", "Total Tiket")
	for i, a := range accumulation {
		row := i + 2
		f.SetCellValue(sheet2, fmt.Sprintf("A%d", row), a.No)
		f.SetCellValue(sheet2, fmt.Sprintf("B%d", row), a.Name)
		f.SetCellValue(sheet2, fmt.Sprintf("C%d", row), a.TotalTiket)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
